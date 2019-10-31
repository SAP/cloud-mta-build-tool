'use strict'

var utils = require('./utils')

var Request = module.exports = function (req, res) {
  var that = this

  this.completed = false
  this._progressTime = process.hrtime()
  this._connection = req.connection
  this._totalBytes = req.headers['content-length']
  this._bytesRead = this._initialBytesRead = req.connection._requestStats ? req.connection._requestStats.bytesRead : 0
  this._bytesWritten = this._initialBytesWritten = req.connection._requestStats ? req.connection._requestStats.bytesWritten : 0

  var done = function () { that.completed = true }
  res.once('finish', done)
  res.once('close', done)
}

Request.prototype.progress = function () {
  var delta = utils.toMilliseconds(process.hrtime(this._progressTime))
  var read = this._connection.bytesRead - this._bytesRead
  var written = this._connection.bytesWritten - this._bytesWritten
  var rSpeed = read / (delta / 1000)
  var wSpeed = written / (delta / 1000)

  this._progressTime = process.hrtime()
  this._bytesRead = this._connection.bytesRead
  this._bytesWritten = this._connection.bytesWritten

  var result = {
    completed: this.completed,
    time: utils.toMilliseconds(process.hrtime(this._progressTime)),
    timeDelta: delta,
    req: {
      bytes: this._connection.bytesRead - this._initialBytesRead,
      bytesDelta: read,
      speed: Number.isNaN(rSpeed) ? 0 : rSpeed
    },
    res: {
      bytes: this._connection.bytesWritten - this._initialBytesWritten,
      bytesDelta: written,
      speed: Number.isNaN(wSpeed) ? 0 : wSpeed
    }
  }

  if (this._totalBytes) {
    var bytesLeft = this._totalBytes - this._connection.bytesRead
    bytesLeft = bytesLeft < 0 ? 0 : bytesLeft
    result.req.bytesLeft = bytesLeft
    result.req.timeLeft = bytesLeft ? bytesLeft / result.req.speed : 0
  }

  return result
}
