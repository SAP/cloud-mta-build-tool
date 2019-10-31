'use strict';

var fs = require('fs');
var debug = require('debug')('xsenv');
var VError = require('verror');

module.exports = loadEnv;

function loadEnv(jsonFile) {
  jsonFile = jsonFile || 'default-env.json';
  if (!fs.existsSync(jsonFile)) {
    return;
  }
  debug('Loading environment from %s', jsonFile);
  try {
    var json = JSON.parse(fs.readFileSync(jsonFile, 'utf8'));
  } catch (err) {
    throw new VError(err, 'Could not parse %s', jsonFile);
  }
  for (var key in json) {
    if (key in process.env) {
      continue; // do not change existing env vars
    }
    var val = json[key];
    // env vars hold only strings
    if (typeof val === 'object') {
      process.env[key] = JSON.stringify(val);
    } else {
      process.env[key] = val + '';
    }
  }
}
