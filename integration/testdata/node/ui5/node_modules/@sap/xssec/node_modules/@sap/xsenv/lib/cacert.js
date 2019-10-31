'use strict';

var debug = require('debug')('xsenv');
var https = require('https');
var path = require('path');
var fs = require('fs');
var util = require('util');
var assert = require('assert');
var VError = require('verror');

exports.loadCaCert = util.deprecate(loadCaCert,
  'xsenv.loadCaCert is deprecated, use xsenv.loadCertificates instead');
exports.loadCertificates = loadCertificates;

function loadCaCert() {
  debug('XS_CACERT_PATH', process.env.XS_CACERT_PATH);
  if (process.env.XS_CACERT_PATH) {
    https.globalAgent.options.ca = loadCertificates(process.env.XS_CACERT_PATH);
  }
}

function loadCertificates(certPath) {
  assert(!certPath || typeof certPath === 'string', 'certPath argument should be a string');

  certPath = certPath || process.env.XS_CACERT_PATH;
  if (certPath) {
    debug('Loading certificate(s) %s', certPath);
    try {
      return certPath
        .split(path.delimiter)
        .map(function (f) {
          return fs.readFileSync(f);
        });
    } catch (err) {
      throw new VError(err, 'Could not load certificate(s) ' + certPath);
    }
  }
}
