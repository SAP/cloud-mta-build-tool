# @sap/xsenv

Utility for easy setup and access of environment variables and services in Cloud Foundry and XSA.

Applications in Cloud Foundry take various configurations from the environment.
For example Cloud Foundry provides properties of bound services in [VCAP_SERVICES](http://docs.cloudfoundry.org/devguide/deploy-apps/environment-variable.html#VCAP-SERVICES) environment variable.

To test locally you need to provide these configurations by yourself. This package allows you to provide default configurations in a separate configuration file.
* This reduces clutter by removing configuration data from the app code.
* You don't have to set env vars manually each time you start your app.
* Different developers can use their own configurations for their local tests without changing files under source control. Just add this configuration file to `.gitignore` and `.cfignore`.

You can provide default configurations on two levels:
* For Cloud Foundry services via `getServices()` and `default-services.json`
* For any environment variable via `loadEnv()` and `default-env.json`

While here we often reference Cloud Foundry, it all applies also to SAP HANA XS Advanced On-premise Runtime which emulates Cloud Foundry as much as possible.

## Install

```sh
npm install --save @sap/xsenv
```

## Service Lookup

Normally in Cloud Foundry you bind a service instance to your application with a command like this one:
```sh
cf bind-service my-app aservice
```

Here is how you can get this service configuration in your node application:
```js
var xsenv = require('@sap/xsenv');

var services = xsenv.readCFServices();
console.log(services.aservice.credentials); // prints { host: '...', port: '...', user: '...', passwrod: '...', ... }
```

Often the service instance name is not known in advance. Then you can pass its name as an environment variable and use it like this:
```js
var services = xsenv.readCFServices();
var svc = services[process.env.SERVICE_NAME];
```

Another alternative is to lookup the service based on its metadata:
```js
var svc = xsenv.cfServiceCredentials({ tag: 'hdb' });
console.log(svc); // prints { host: '...', port: '...', user: '...', passwrod: '...', ... }
```
This example finds a service binding with `hdb` in the tags.
See [Service Query](#service-query) below for description of the supported query values.

You can also look up multiple services in a single call:
```js
var xsenv = require('@sap/xsenv');

var services = xsenv.getServices({
  hana: { tag: 'hdb' },
  scheduler: { label: 'jobs' }
});

var hanaCredentials = services.hana;
var schedulerCredentials = services.scheduler;
```
This examples finds two services - one with tag `hdb` and the other with label `jobs`.
`getServices` function provides additional convenience that default service configuration can be provided in a JSON file.

To test the above example locally, create a file called `default-services.json` in the working directory of your application.
This file should contain something like this:
```json
{
  "hana": {
    "host": "localhost",
    "port": "30015",
    "user": "SYSTEM",
    "password": "secret"
  },
  "scheduler": {
    "host": "localhost",
    "port": "4242",
    "user": "my_user",
    "password": "secret"
  }
}
```
Notice that the result property names (`hana` and `scheduler`) are the same as those in the query object and also those in `default-services.json`.

[Local environment setup](#local-environment-setup) describes an alternative approach to provide service configurations for local testing.

### User-Provided Service Instances

While this package can look up any kind of bound service instances, you should be aware that [User-Provided Service Instances](https://docs.cloudfoundry.org/devguide/services/user-provided.html) have less properties than managed service instances. Here is an example:
```json
  "VCAP_SERVICES": {
    "user-provided": [
      {
        "name": "pubsub",
        "label": "user-provided",
        "tags": [],
        "credentials": {
          "binary": "pubsub.rb",
          "host": "pubsub01.example.com",
          "password": "p@29w0rd",
          "port": "1234",
          "username": "pubsubuser"
        },
        "syslog_drain_url": ""
	  }
    ]
  }
```
As you can see the only usable property is the `name`. In particular, there are no tags for user-provided services.

### Service Query

Both `getServices` and `filterCFServices` use the same service query values.

Query value | Description
------------|------------
{string}    | Matches the service with the same service instance name (`name` property). Same as `{ name: '<string>' }`.
{object}    | All properties of the given object should match corresponding service instance properties as they appear in VCAP_SERVICES. See below.
{function}  | A function that takes a service object as argument and returns `true` only if it is considered a match

If an object is given as a query value, it may have the following properties:

Property | Description
---------|------------
`name`   | Service instance name - the name you use to bind the service
`label`  | Service name - the name shown by `cf marketplace`
`tag`    | Should match any of the service tags
`plan`   | Service instance plan - the plan you use in `cf create-service`

If multiple properties are given, _all_ of them must match.

**Note:** Do not confuse the instance name (`name` property) with service name (`label` property).
Since you can have multiple instances of the same service bound to your app,
instance name is unique while service name is not.

Here ar some examples.

Find service instance by name:
```js
xsenv.cfServiceCredentials('hana');
```

Look up a service by tag:
```js
xsenv.cfServiceCredentials({ tag: 'relational' });
```

Match several properties:
```js
xsenv.cfServiceCredentials({ label: 'hana', plan: 'shared' });
```

Pass a custom filter function:
```js
xsenv.cfServiceCredentials(function(service) {
  return /shared/.test(service.plan) && /hdi/.test(service.label);
});
```
Notice that the filter function is called with the full service object as it appears in VCAP_SERVICES, but `cfServiceCredentials` returns only the `credentials` property of the matching service.

### getServices(query, [servicesFile])

Looks up bound service instances matching the given query.
If a service is not found in VCAP_SERVICES, returns default service configuration loaded from a JSON file.

* `query` - An object describing requested services. Each property value is a filter as described in [Service Query](#service-query)
* `servicesFile` - (optional) path to JSON file to load default service configuration (default is default-services.json).
If `null`, do not load default service configuration.
* _returns_ - An object with the same properties as in query argument where the value of each property is the respective service credentials object
* _throws_ - An error, if for some of the requested services no or multiple instances are found

### cfServiceCredentials(filter)

Looks up a bound service instance matching the given filter.

**Note:** this function does not load default service configuration from default-services.json.

* `filter` - Service lookup criteria as described in [Service Query](#service-query)
* _returns_ - Credentials object of found service
* _throws_ - An error in case no or multiple matching services are found

### filterCFServices(filter)

Returns all bound services that match the given criteria.

* `filter` - Service lookup criteria as described in [Service Query](#service-query)
* _returns_ - An array of credentials objects of matching services. Empty array, if no matches found.

### readCFServices()

Transforms VCAP_SERVICES to an object where each service instance is mapped to its name.

Given this VCAP_SERVICES example:
```
  {
    "hana" : [ {
      "credentials" : {
        ...
      },
      "label" : "hana",
      "name" : "hana1",
      "plan" : "shared",
      "tags" : [ "hana", "relational" ]
    },
    {
      "credentials" : {
        ...
      },
      "label" : "hana",
      "name" : "hana2",
      "plan" : "shared",
      "tags" : [ "hana", "relational", "SP09" ]
    } ]
  }
```
`readCFServices` would return:
```
{
  hana1: {
    "credentials" : {
      ...
    },
    "label" : "hana",
    "name" : "hana1",
    "plan" : "shared",
    "tags" : [ "hana", "relational" ]
  },
  hana2: {
    "credentials" : {
      ...
    },
    "label" : "hana",
    "name" : "hana2",
    "plan" : "shared",
    "tags" : [ "hana", "relational", "SP09" ]
  }
}
```

## Local environment setup

To test your application locally you often need to setup its environment so that resembles the environment in Cloud Foundry.
You can do this easily by defining the necessary environment variables in a JSON file.

For example you can create file _default-env.json_ with the following content in the working directory of the application :
```json
{
  "PORT": 3000,
  "VCAP_SERVICES": {
    "hana": [
      {
        "credentials": {
          "host": "myhana",
          "port": "30015",
          "user": "SYSTEM",
          "password": "secret"
        },
        "label": "hana",
        "name": "hana-R90",
        "tags": [
          "hana",
          "database",
          "relational"
        ]
      }
    ],
    "scheduler": [
      {
        "credentials": {
          "host": "localhost",
          "port": "4242",
          "user": "jobuser",
          "password": "jobpassword"
        },
        "label": "scheduler",
        "name": "jobscheduler",
        "tags": [
          "scheduler"
        ]
      }
    ]
  }
}
```
Then load it in your application:
```js
xsenv.loadEnv();
console.log(process.env.PORT); // prints 3000
console.log(xsenv.cfServiceCredentials('hana-R90')); // prints { host: 'myhana, port: '30015', user: 'SYSTEM', password: 'secret' }
```

This way you don't need in your code conditional logic if it is running in Clod Foundry or locally.

You can also use a different file name:
```js
xsenv.loadEnv('myenv.json');
```

### loadEnv([file])

Loads the environment from a JSON file.
This function converts each top-level property to a string and sets it in the respective environment variable,
unless it is already set. This function does not change existing environment variables. So the file content acts like default values for environment variables.

This function does not complain if the file does not exist.

* `file` - optional name of JSON file to load, `'default-env.json'` by default. Does nothing if the file does not exist.

## Loading SSL Certificates

If SSL is configured in XSA On-Premise Runtime, it will provide one or more
trusted CA certificates that applications can use to make SSL connections.
If present, the file paths of these certificates are listed in `XS_CACERT_PATH`
environment variable separated by `path.delimiter` (`:` on LINUX and `;` on Windows).

### loadCertificates([certPath])

Loads the certificates listed in the given path.
If this argument is not provided, it uses `XS_CACERT_PATH` environment variable instead.
If that is not set either, the function returns `undefined`.
The function returns an array even if a single certificate is provided.
This function is synchronous.

* `certPath` - optional string with certificate files to load. The file names are separated by `path.delimiter`. Default is `process.env.XS_CACERT_PATH`.
* _returns_ - an array of loaded certificates or `undefined` if `certPath` argument is not provided 
* _throws_ - an error, if some of the files could not be loaded

For example this code loads the trusted CA certificates so they are used for all
subsequent outgoing HTTPS connections:
```js
var https = require('https');
var xsenv = require('@sap/xsenv');

https.globalAgent.options.ca = xsenv.loadCertificates();
```

This function can be used also to load SSL certificates for HANA like this:
```js
var hdb = require('hdb');
var xsenv = require('@sap/xsenv');

var client = hdb.createClient({
  host : 'hostname',
  port : 30015,
  ca   : xsenv.loadCertificates(),
  ...
});
```

### loadCaCert()

**Deprecated.** Use `loadCertificates` instead.

This function loads the certificates listed in `XS_CACERT_PATH` environment variable
into [https.globalAgent](https://nodejs.org/api/https.html#https_https_globalagent) of Node.js.
All subsequent outgoing HTTPS connections will use these certificates to verify the certificate
of the remote host. The verification should be successful if the certificate of the
remote host is signed (directly or via some intermediary) by some of these trusted
certificates.

It is suggested to call this function once during the application startup.

```js
xsenv.loadCaCert();
```

If `XS_CACERT_PATH` variable is not set, the function does nothing.
This function is synchronous.
It will throw an error if it cannot load some of the certificates.

## Debugging

Set `DEBUG=xsenv` in the environment to enable debug traces. See [debug](https://www.npmjs.com/package/debug) package for details.
