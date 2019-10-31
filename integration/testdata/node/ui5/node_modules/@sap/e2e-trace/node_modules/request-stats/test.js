'use strict'

var http = require('http')
var https = require('https')
var pem = require('https-pem')
var test = require('tape')
var EventEmitter = require('events').EventEmitter
var requestStats = require('./index')
var Request = require('./lib/request')
var StatsEmitter = require('./lib/stats_emitter')
var KeepAliveAgent = require('keep-alive-agent')

var assertStatsCommon = function (t, stats) {
  t.ok(stats.req.bytes > 0, 'stats.req.bytes > 0') // different headers will result in different results
  t.equal(typeof stats.req.headers.connection, 'string')
  t.equal(stats.req.method, 'PUT')
  t.equal(stats.req.path, '/')
  t.equal(stats.res.status, 200)
}

// assertion helper for validating stats from HTTP requests that finished correctly
var assertStatsFinished = function (t, stats) {
  t.ok(stats.ok, 'stats.ok should be truthy')
  t.ok(stats.time >= 9, 'stats.time >= 9') // The reason we don't just do >= 10, is because setTimeout is not that precise
  t.ok(stats.res.bytes > 0, 'stats.res.bytes > 0') // different headers will result in different results
  t.equal(typeof stats.res.headers.connection, 'string')
  assertStatsCommon(t, stats)
}

// assertion helper for validating stats from HTTP requests that are closed before finishing
var assertStatsClosed = function (t, stats) {
  t.ok(!stats.ok)
  t.ok(stats.time >= 0)
  t.equal(stats.res.bytes, 0)
  t.deepEqual(stats.res.headers, {})
  assertStatsCommon(t, stats)
}

var _start = function (server, errorHandler) {
  server.listen(0, function () {
    var transport = server.key ? https : http
    var options = {
      port: server.address().port,
      method: 'PUT',
      rejectUnauthorized: false // if https
    }

    var req = transport.request(options, function (res) {
      res.resume()
      res.once('end', function () {
        server.close()
      })
    })

    if (errorHandler) req.on('error', errorHandler)

    req.end('42')
  })
}

var _respond = function (req, res) {
  req.on('end', function () {
    setTimeout(function () {
      res.end('Answer to the Ultimate Question of Life, The Universe, and Everything')
    }, 10)
  })
  req.resume()
}

test('StatsEmitter', function (t) {
  var statsEmitter = requestStats()
  t.ok(statsEmitter instanceof StatsEmitter, 'should be returned from requestStats()')
  t.ok(statsEmitter instanceof EventEmitter, 'should be an instance of EventEmitter')

  t.test('should emit a "request" event', function (t) {
    var server = http.createServer(_respond)
    requestStats(server).on('request', function () {
      t.end()
    })
    _start(server)
  })

  t.test('should emit a "complete" event', function (t) {
    var server = http.createServer(_respond)
    requestStats(server).on('complete', function () {
      t.end()
    })
    _start(server)
  })
})

test('requestStats(server, onStats)', function (t) {
  t.test('should call the stats-listener on request end', function (t) {
    var server = http.createServer(_respond)
    requestStats(server, function (stats) {
      assertStatsFinished(t, stats)
      t.end()
    })
    _start(server)
  })

  t.test('should call the stats-listener when the request is destroyed', function (t) {
    var server = http.createServer(function (req, res) {
      req.destroy()
    })
    requestStats(server, function (stats) {
      assertStatsClosed(t, stats)
    })
    _start(server, function (err) {
      t.ok(err instanceof Error)
      t.end()
    })
  })

  t.test('should calculate correct bytes read/written on keep-alive connections', function (t) {
    var agent = new http.Agent()
    agent.maxSockets = 1 // force connection reuse for every connection

    var server = http.createServer(function (req, res) {
      req.on('end', function () {
        res.end()
      })
      req.resume()
    })

    server.once('connection', function () {
      server.once('connection', function () {
        t.fail('Expected the TCP connection to be reused')
      })
    })

    requestStats(server, function (stats) {
      // assert req.bytes are around 1 million (we cannot know exactly since
      // request headers will vary depending on the host system)
      t.ok(stats.req.bytes > 1000000, 'req body size should be above 1MB')
      t.ok(stats.req.bytes < 1000300, 'req body size should not be to much above 1MB')
    })

    var performRequest = function (port, callback) {
      var req = http.request({ port: port, method: 'PUT', agent: agent }, function (res) {
        res.on('end', callback)
        res.resume()
      })
      for (var n = 0, s = ''; n < 100000; n++) s += '1234567890' // 1 million
      req.end(s)
    }

    server.listen(0, function () {
      var port = server.address().port
      performRequest(port, function () {
        performRequest(port, function () {
          server.close()
          t.end()
        })
      })
    })
  })
})

test('requestStats(req, res).once(...)', function (t) {
  t.test('should call the stats-listener on request end', function (t) {
    _start(http.createServer(function (req, res) {
      requestStats(req, res).once('complete', function (stats) {
        assertStatsFinished(t, stats)
        t.end()
      })
      _respond(req, res)
    }))
  })
})

test('requestStats(server).once(...)', function (t) {
  t.test('should call the stats-listener on request end', function (t) {
    var server = http.createServer(_respond)
    requestStats(server).once('complete', function (stats) {
      assertStatsFinished(t, stats)
      t.end()
    })
    _start(server)
  })
})

test('requestStats(https-server)', function (t) {
  var server = https.createServer(pem, _respond)
  requestStats(server).once('complete', function (stats) {
    assertStatsFinished(t, stats)
    t.end()
  })
  _start(server)
})

test('requestStats(req, res, onStats)', function (t) {
  t.test('should call the stats-listener on request end', function (t) {
    _start(http.createServer(function (req, res) {
      requestStats(req, res, function (stats) {
        assertStatsFinished(t, stats)
        t.end()
      })
      _respond(req, res)
    }))
  })
})

test('Request instance', function (t) {
  t.test('should expose a .progress() function', function (t) {
    _start(http.createServer(function (req, res) {
      var request = new Request(req, res)
      t.ok(typeof request.progress === 'function')
      t.end()
    }))
  })

  t.test('should be emitted on the "request" event', function (t) {
    var server = http.createServer(_respond)
    var statsEmitter = requestStats(server)
    var request
    statsEmitter.on('request', function (obj) {
      request = obj
    })
    statsEmitter.on('complete', function (stats) {
      t.ok(request instanceof Request)
      t.end()
    })
    _start(server)
  })
})

test('request.progress()', function (t) {
  t.test('should return a progress object', function (t) {
    var server = http.createServer(_respond)
    var statsEmitter = requestStats(server)
    statsEmitter.on('request', function (request) {
      var progress = request.progress()
      t.equal(progress.completed, false, 'should be completed')
      t.ok(progress.time >= 0, 'progress.time >= 0')
      t.ok(progress.timeDelta >= 0, 'progress.timeDelta >= 0')
      t.ok(progress.req.bytes > 0, 'progress.req.bytes > 0')
      t.ok(progress.req.bytesDelta > 0, 'progress.req.bytesDelta > 0')
      t.ok(progress.req.speed > 0, 'progress.req.speed > 0')
      t.equal(progress.res.bytes, 0, 'progress.res.bytes === 0')
      t.equal(progress.res.bytesDelta, 0, 'progress.res.bytesDelta === 0')
      t.equal(progress.res.speed, 0, 'progress.res.speed === 0')
      t.end()
    })
    _start(server)
  })

  t.test('should not mix progress from two request', function (t) {
    var server = http.createServer(_respond)
    var statsEmitter = requestStats(server)
    var requests = []
    var progress = []

    statsEmitter.on('request', function (request) {
      requests.push(request)
      t.equal(typeof request._connection, 'object')
    })

    statsEmitter.on('complete', function (stats) {
      progress.push(requests[requests.length - 1].progress())
      if (requests.length < 2) return
      t.equal(requests[0]._connection, requests[1]._connection, 'should re-use the http connection')
      t.equal(progress[0].req.bytes, progress[1].req.bytes, 'should receive the same amount of data')
      t.equal(progress[0].res.bytes, progress[1].res.bytes, 'should send the same amount of data')
      t.equal(progress[0].req.bytes + progress[1].req.bytes, requests[0]._connection.bytesRead, 'should not accumulate received data')
      t.equal(progress[0].res.bytes + progress[1].res.bytes, requests[0]._connection.bytesWritten, 'should not accumulate sent data')
      t.end()
    })

    server.listen(0, function () {
      var options = {
        port: server.address().port,
        method: 'PUT',
        headers: { Connection: 'keep-alive' },
        agent: new KeepAliveAgent()
      }

      http.request(options, function (res) {
        res.resume()
      }).end('42')

      setTimeout(function () {
        http.request(options, function (res) {
          res.resume()
          res.once('end', function () {
            server.close()
          })
        }).end('42')
      }, 100)
    })
  })
})

test('end', function (t) {
  t.end()
  process.exit()
})
