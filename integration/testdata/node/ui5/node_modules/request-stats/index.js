'use strict'

var http = require('http')
var StatsEmitter = require('./lib/stats_emitter')

module.exports = function (req, res, onStats) {
  var statsEmitter = new StatsEmitter()
  if (req instanceof http.Server) statsEmitter._server(req, res)
  else if (req instanceof http.IncomingMessage) statsEmitter._request(req, res, onStats)
  return statsEmitter
}
