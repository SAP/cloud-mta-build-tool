'use strict'

exports.toMilliseconds = function (tuple) {
  return Math.round(tuple[0] * 1000 + tuple[1] / 1000000)
}
