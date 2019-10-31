'use strict'

var util = require('util')
var EventEmitter = require('events').EventEmitter
var once = require('once')
var httpHeaders = require('http-headers')
var Request = require('./request')
var utils = require('./utils')

var StatsEmitter = module.exports = function () {
  EventEmitter.call(this)
}
util.inherits(StatsEmitter, EventEmitter)

StatsEmitter.prototype._server = function (server, onStats) {
  this._attach(onStats)
  server.on('request', this._request.bind(this))
}

StatsEmitter.prototype._request = function (req, res, onStats) {
  var that = this
  var start = process.hrtime()

  this.emit('request', new Request(req, res))
  this._attach(onStats)

  var emit = once(function (ok) {
    var bytesReadPreviously = req.connection._requestStats ? req.connection._requestStats.bytesRead : 0
    var bytesWrittenPreviously = req.connection._requestStats ? req.connection._requestStats.bytesWritten : 0
    var bytesReadDelta = req.connection.bytesRead - bytesReadPreviously
    var bytesWrittenDelta = req.connection.bytesWritten - bytesWrittenPreviously
    var ip = getIp(req)

    req.connection._requestStats = {
      bytesRead: req.connection.bytesRead,
      bytesWritten: req.connection.bytesWritten
    }

    that.emit('complete', {
      ok: ok,
      time: utils.toMilliseconds(process.hrtime(start)),
      req: {
        bytes: bytesReadDelta,
        headers: req.headers,
        method: req.method,
        path: req.url,
        ip: ip,
        raw: req
      },
      res: {
        bytes: bytesWrittenDelta,
        headers: httpHeaders(res, true),
        status: res.statusCode,
        raw: res
      }
    })
  })

  res.once('finish', emit.bind(null, true))
  res.once('close', emit.bind(null, false))
}

function getIp (req) {
  var ip = req.headers['x-forwarded-for']
  if (!ip) {
    ip = req.connection && req.connection.remoteAddress
    if (!ip && req.socket) {
      ip = req.socket.remoteAddress || (req.socket.connection && req.socket.connection.remoteAddress)
    }
  }
  return ip
}

StatsEmitter.prototype._attach = function (listener) {
  if (typeof listener === 'function') this.on('complete', listener)
}
