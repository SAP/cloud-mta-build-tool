'use strict';

var assert = require('assert');
var VError = require('verror');

exports.readCFServices = readCFServices;
exports.cfServiceCredentials = cfServiceCredentials;
exports.filterCFServices = filterCFServices;

function readCFServices() {
  if (!process.env.VCAP_SERVICES) {
    return;
  }
  try {
    var services = JSON.parse(process.env.VCAP_SERVICES);
  } catch (err) {
    throw new VError(err, 'Environment variable VCAP_SERVICES is not a valid JSON string.');
  }

  var result = {};
  for (var s in services) {
    for (var si in services[s]) {
      var svc = services[s][si];
      result[svc.name] = svc; // name is the service instance id
    }
  }
  return result;
}

/**
 * Reads service configuration from CloudFoundry environment variable <code>VCAP_SERVICES</code>.
 *
 * @param filter Filter used to find a bound Cloud Foundry service, see filterCFServices
 * @return credentials property of found service
 * @throws Error in case no or multiple matching services are found
 */
function cfServiceCredentials(filter) {
  var matches = filterCFServices(filter);
  if (matches.length !== 1) {
    throw new VError('Found %d matching services', matches.length);
  }
  return matches[0].credentials;
}

/**
 * Returns an array of Cloud Foundry services matching the given filter.
 *
 * @param filter {(string|Object|function)}
 *  - if string, returns the service with the same service instance name (name property)
 *  - if Object, should have some of these properties [name, label, tag, plan] and returns all services
 *    where all of the given properties match. Given tag matches if it is present in the tags array.
 *  - if function, should take a service object as argument and return true only if it matches the filter
 * @returns Arrays of matching service objects, empty if no matches
 */
function filterCFServices(filter) {
  assert(typeof filter === 'string' || typeof filter === 'object' || typeof filter === 'function',
    'bad filter type: ' + typeof filter);

  var services = readCFServices();
  if (!services) {
    return [];
  }

  if (typeof filter === 'string') {
    return services[filter] ? [services[filter]] : [];
  }

  var result = [];
  for (var key in services) {
    if (applyFilter(services[key], filter)) {
      result.push(services[key]);
    }
  }
  return result;
}

function applyFilter(service, filter) {
  if (typeof filter === 'function') {
    return filter(service);
  }

  var match = false;
  for (var key in filter) {
    if (service[key] === filter[key] ||
      (/tags?/.test(key) && service.tags && service.tags.indexOf(filter[key]) >= 0)) {
      match = true;
    } else {
      return false;
    }
  }
  return match;
}


