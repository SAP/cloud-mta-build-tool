# @sap/audit-logging

Provides audit logging functionalities for Node.js applications.

<!-- toc -->

- [Overview](#overview)
  * [General audit logging principles](#general-audit-logging-principles)
  * [Prerequisites](#prerequisites)
  * [Versions](#versions)
- [API - v1](#api---v1)
  * [Importing the library](#importing-the-library)
  * [Data access messages](#data-access-messages)
  * [Data modification messages](#data-modification-messages)
  * [Update data modification](#update-data-modification)
  * [Configuration change messages](#configuration-change-messages)
  * [Update configuration change](#update-configuration-change)
  * [General security messages](#general-security-messages)
  * [Logging a message](#logging-a-message)
- [API - v2](#api---v2)
  * [Importing the library](#importing-the-library-1)
  * [Data access messages](#data-access-messages-1)
  * [Data modification messages](#data-modification-messages-1)
  * [Configuration change messages](#configuration-change-messages-1)
  * [General security messages](#general-security-messages-1)
- [Local development](#local-development)
  * [Without Audit log service](#without-audit-log-service)
  * [With Audit log service](#with-audit-log-service)

<!-- tocstop -->

## Overview

Audit logging is about writing entries in a specific format to a log storage. Subject to audit logging are events of significant importance.
For example, security events which may impact the confidentiality, the integrity or the availability of a system.
Another example of such an event would be access to personal data (both reading and altering) like bank accounts, political opinion,
health status etc.

While the consumer of ordinary logs is a system administrator who would like to keep track of the state of a system,
audit logs are read by an auditor. There are legal requirements (in some countries stricter than in others) regarding audit logging.

In general the events that are supposed to be audit logged can be grouped in 3 main categories:
- changes to system configurations (which may have significant effect on the system itself)
- access to personal data (related to data privacy)
- general security events (like starting/stopping a system, failed authorization checks etc.)


### General audit logging principles

- All attempts to perform an action in a system should be audit logged no matter if they have been successful or not.
- Audit log entries should be consistent with the state of the system. If, for example, the writing of the audit log entry fails,
but the changes to system critical parameters have been applied, then those changes should be reverted. Best practice is to wait for
the callback of the logger before continuing with the execution of other code.
- Especially important is which user (or other agent) has triggered the corresponding event that is being audit logged.
For most of the cases the library will validate that such a field is provided in the message.
- All audit log entries should be in English. Numbers should be converted to strings with English locale.
All time fields should be in UTC time in order to avoid timezone and day light saving time issues.
- Passwords should never be audit logged.

### Prerequisites

An application using the audit log library needs to be bound to an instance of the Audit log service.

### Versions

The Audit log service provides REST APIs that are available to applications for
logging relevant messages. The latest Audit log server supports 2 versions
of the REST APIs. This library provides JavaScript programming interfaces for
both of these versions of the server's APIs.
**Note:** It is recommended to use REST APIs v2 if available on the Audit log server being in use (and respectively the JavaScript v2 APIs).
The initial version of the Audit log server REST APIs is deprecated in favor of the v2 version. The same applies to the JavaScript APIs provided by this library.

## API - v1

The library provides an API for writing audit messages of type configuration changes, data modifications, data accesses and security events.

### Importing the library

```js
var credentials = {
  "user": "user",
  "password": "password",
  "url": "https://host:port"
};
var auditLog = require('@sap/audit-logging')(credentials);
```

`credentials` object is the bound audit log service's credentials.
Take a look at *@sap/xsenv* package for more information on how to retrieve service credentials.

### Data access messages

Let's suppose we need to create an entry for a data access operation over personal data. We can achieve that with the following code:

```js
auditLog.read('user123')
  .attribute('username', true)
  .attribute('first name', true)
  .attribute('last name', true)
  .accessChannel('UI')
  .by('John Doe')
  .tenant('tenantId')
  .log(...);
```

* `read` - takes a string which identifies the object which is being *accessed*.
* `attribute(name, successful)` - sets object attributes. It is **mandatory** to provide at least one attribute.
  * `name` - is the name of the attribute being accessed.
  * `successful` - specifies whether the access was successful or not.
* `by` - takes a string which identifies the *user* performing the action. This is **mandatory**.
* `accessChannel` - takes a string which specifies *channel* of access.
* `attachment(id, name)` - if attachments or files are downloaded or displayed, information identifying the attachment shall be logged.
  * `id` - attachment id
  * `name` - attachment name
* `tenant` - takes a string which specifies the tenant id. The provided value is ignored by older versions of the Audit log service that do not support setting a tenant.
* `log` - See [here](#logging-a-message)

### Data modification messages

Here is how to create an entry for a data modification operation:

```js
auditLog.update('userdata')
  .attribute('first name', 'john', 'John')
  .by('John Doe')
  .tenant('tenantId')
  .log(...);
```

**Note**: Specifying an old and a new value for an attribute is only supported in newer versions of the Audit log service. Providing these values while working with an older version of the service results in an error in the callback. In such cases one may use the `attribute` method with an alternative signature:

```js
auditLog.update('userdata')
  .attribute('password', false)
  .by('John Doe')
  .tenant('tenantId')
  .log(...);
```

* `update` - takes a string which identifies the object which is being *updated*.
* `attribute(name, oldValue, newValue)` - sets object attributes. It is **mandatory** to provide at least one attribute.
  * `name` - is the name of the attribute being modified.
  * `oldValue` - is the current value of the attribute.
  * `newValue` - is the value of the attribute after the change.

  **Note**: One may use this signature of the `attribute` method only if the Audit log service being consumed supports old and new values.

* `attribute(name, successful)` - sets object attributes. It is **mandatory** to provide at least one attribute.
  * `name` - is the name of the attribute being modified.
  * `successful` - specifies whether the modification was successful or not.

  **Note**: this signature of the method is **deprecated**. It should be used only if the consumed Audit log service does not support old and new values.

* `by` - takes a string which identifies the *user* performing the action. This is **mandatory**.
* `tenant` - takes a string which specifies the tenant id. The provided value is ignored by older versions of the Audit log service that do not support setting a tenant.
* `log` - See [here](#logging-a-message)

### Update data modification

```js
auditLog.updateDataModification(id, isSuccessful)
  .log(...);
```

* `updateDataModification(id, isSuccessful)` - takes two arguments.
  * `id` - id of the data modification message saved earlier (see [log](#logging-a-message))
  * `isSuccessful` - denotes whether the data modification was successful or not.
* `log` - See [here](#logging-a-message)

**Note**: This function should only be used with an Audit log service that supports old and new values.

### Configuration change messages

Here is how to create an entry for a configuration change operation:

```js
auditLog.configurationChange('configuration object')
  .attribute('session timeout', '5', '25')
  .by('Application Admin')
  .successful(true)
  .tenant('tenantId')
  .log(...);
```

* `configurationChange` - takes a string which identifies the object which is being *configured*.
* `attribute(name, oldValue, newValue)` - sets object attributes. It is **mandatory** to provide at least one attribute.
  * `name` - is the name of the attribute being accessed.
  * `oldValue` - is the current value of the attribute being changed.
  * `newValue` - is the value of the attribute after the change.
* `successful(isSuccessful)` - used to mark whether the configuration change is finished with success, failure.
  If not called configuration change will be marked as *pending*.
  * `isSuccessful` - should be a valid boolean.
* `by` - takes a string which identifies the *user* performing the action. This is **mandatory**.
* `tenant` - takes a string which specifies the tenant id. The provided value is ignored by older versions of the Audit log service that do not support setting a tenant.
* `log` - See [here](#logging-a-message)

### Update configuration change

```js
auditLog.updateConfigurationChange(id, isSuccessful)
  .log(...);
```

* `updateConfigurationChange(id, isSuccessful)` - takes two arguments.
  * `id` - id of the configuration message saved earlier (see [log](#logging-a-message))
  * `isSuccessful` - denotes whether the configuration change was successful or not.
* `log` - See [here](#logging-a-message)

### General security messages

Here is how to create a general security audit log message:

```js
auditLog.securityMessage('%d unsuccessful login attempts', 3)
  .by('John Doe')
  .externalIP('127.0.0.1')
  .tenant('tenantId')
  .log(...);
```

* `securityMessage` - takes a formatted string as in [util.format()](https://nodejs.org/api/util.html#util_util_format_format_args).
* `externalIP` - states the IP of the machine that contacts the system. It is not mandatory, but it should be either IPv4 or IPv6.
* `by` - takes a string which identifies the *user* performing the action. This is **mandatory**.
* `tenant` - takes a string which specifies the tenant id. The provided value is ignored by older versions of the Audit log service that do not support setting a tenant.
* `log` - See [here](#logging-a-message)

### Logging a message

Use the `log` method to write a message to the Audit log. It takes one argument - a callback function.
Be aware that the state of the audit logs should be consistent with the state of the system.
Make sure you handle errors from the audit log writer properly.
Application code **should wait** for the logging to finish before executing any other code.

```js
var message = /* any of the above example messages */;
message.log(function (err, id) {
    // Do error handling and place all of the remaining logic here
  });
```

* `message` - Any of the following:
  * [`read`](#data-access-messages)
  * [`update`](#data-modification-messages)
  * [`configurationChange`](#configuration-change-messages)
  * [`updateConfigurationChange`](#update-configuration-change)
  * [`securityMessage`](#general-security-messages)
* `err` - error object in case of error.
* `id` - Id of the message that is saved. Use it when you want to do [`updateConfigurationChange`](#update-configuration-change). `id` is undefined in case of [`updateConfigurationChange`](#update-configuration-change).

**Note**: When a message is logged, the library checks for missing properties and will throw an error if some are missing.

## API - v2

### Importing the library

```js
var credentials = {
  "user": "user",
  "password": "password",
  "url": "https://host:port"
};

var auditLogging = require('@sap/audit-logging');
auditLogging.v2(credentials, function(err, auditLog) {
  if (err) {
    // if the Audit log server does not support version 2 of the REST APIs
    // an error in the callback is returned
    return console.log(err);
  }
});
```

`credentials` object with credentials for the Audit log service.
Take a look at *@sap/xsenv* package for more information on how to retrieve service credentials.
The callback will be called with an error if the Audit log server does not support version 2 of the REST APIs.

### Data access messages

```js
auditLog.read({ type: 'accessed-object-type', id: { key: 'value' } })
  .attribute({ name: 'attr-0' })
  .attribute({ name: 'attr-1', successful: true })
  .attachment({ id: '123' })
  .attachment({ id: '456', name: 'file.doc' })
  .dataSubject({ type: 'data-subject-type', id: { key: 'value' }, role: 'role' })
  .accessChannel('UI')
  .tenant('tenantId')
  .by('John Doe')
  .log(function (err) {

  });
```

* `read` - takes a JavaScript object which identifies the object which contains the data being accessed. Should have `type` and `id` properties.
* `attribute(attribute)` - takes an object which describes an attribute. Should have a `name` property and optionally a `successful` property. It is **mandatory** to provide at least one attribute.
* `attachment(attachment)` - takes an object which describes an attachment (used if attachments or files are downloaded or displayed). Should have an `id` property and optionally a `name` property.
* `dataSubject` - takes an object describing the owner of the personal data. Should have `type` and `id` properties. The `role` property is optional. `dataSubject` is **mandatory**.
* `accessChannel` - takes a string which specifies *channel* of access.
* `tenant` - takes a string which specifies the tenant id.
* `by` - takes a string which identifies the *user* performing the action. This is **mandatory**.
* `log` - logs the message.

### Data modification messages

```js
var message = auditLog.update({ type: 'accessed-object-type', id: { key: 'value' } })
  .attribute({ name: 'attr-0' })
  .attribute({ name: 'attr-1' })
  .attribute({ name: 'attr-2', old: 'old value', new: 'new value' })
  .dataSubject({ type: 'data-subject-type', id: { key: 'value' }, role: 'role' })
  .tenant('tenantId')
  .by('John Doe');

message.logPrepare(function (err) {
  message.logSuccess(function (err) { });
  // or
  message.logFailure(function(err) { });
});
```

* `update` - takes a JavaScript object which identifies the object which contains the data being updated. Should have `type` and `id` properties.
* `attribute(attribute)` - takes an object which describes an attribute. Should have a `name` property and optionally - `old` and `new` properties. It is **mandatory** to provide at least one attribute.
* `dataSubject` - takes an object describing the owner of the personal data. Should have `type` and `id` properties. The `role` property is optional. `dataSubject` is **mandatory**.
* `tenant` - takes a string which specifies the tenant id.
* `by` - takes a string which identifies the *user* performing the action. This is **mandatory**.
* `logPrepare` - Used to log that a user has started an operation over the data.
* `logSuccess` - Used to log that the operation over the data has been completed successfully.
* `logFailure` - Used to log that the operation over the data has not been completed successfully.

### Configuration change messages

```js
var message = auditLog.configurationChange({ type: 'accessed-object-type', id: { key: 'value' } })
  .attribute({ name: 'session timeout', old: '5', new: '25' })
  .tenant('tenantId')
  .by('Application Admin');

message.logPrepare(function (err) {
  message.logSuccess(function (err) { });
  // or
  message.logFailure(function(err) { });
});
```

* `configurationChange` - takes a JavaScript object which identifies the object which contains the data being configured. Should have `type` and `id` properties.
* `attribute(attribute)` - takes an object which describes an attribute. Should have a `name`, `old` and `new` properties. It is **mandatory** to provide at least one attribute.
* `tenant` - takes a string which specifies the tenant id.
* `by` - takes a string which identifies the *user* performing the action. This is **mandatory**.
* `logPrepare` - Used to log that a user has started a configuration change operation.
* `logSuccess` - Used to log that the operation has been completed successfully.
* `logFailure` - Used to log that the operation has not been completed successfully.

### General security messages

```js
auditLog.securityMessage('%d unsuccessful login attempts', 3)
  .by('John Doe')
  .externalIP('127.0.0.1')
  .tenant('tenantId')
  .log(function (err) {

  });
```

* `securityMessage` - takes a formatted string as in [util.format()](https://nodejs.org/api/util.html#util_util_format_format_args).
* `externalIP` - states the IP of the machine that contacts the system. Specifying it is optional, but if provided, should be either IPv4 or IPv6.
* `by` - takes a string which identifies the *user* performing the action. This is **mandatory**.
* `tenant` - takes a string which specifies the tenant id.
* `log` - logs the message.


## Local development

### Without Audit log service

```js
var credentials = {
  logToConsole: true
};
var auditLog = require('@sap/audit-logging')(credentials);

// or

require('@sap/audit-logging').v2(credentials, function (err, auditLog) {

});
```

When `logToConsole` is `true` the library will ignore other credential properties and will not use the Audit log service,
but will write the messages to the console.

**Hint:** If you use the *@sap/xsenv* package, you can pass the credentials through the *default-services.json* file
or `VCAP_SERVICES` environment variable.

### With Audit log service

If your application is not deployed in Cloud Foundry or XS Advanced,
but you have a running Audit log service somewhere, you should set the `VCAP_APPLICATION` environment variable to a string like
`{ "application_name" : "my-app", "organization_name" : "my-org", "space_name" : "my-space" }`

**Hint:** If you use the *@sap/xsenv* package, you can set environment variables like this:

```js
var xsenv = require('@sap/xsenv');

xsenv.loadEnv();
var credentials = xsenv.getServices({ auditlog: 'auditlog-instance-name' }).auditlog;
var auditLog = require('@sap/audit-logging')(credentials);
```

*default-env.json* file:

```json
{
  "VCAP_APPLICATION": {
    "application_name" : "my-app",
    "organization_name" : "my-org",
    "space_name" : "my-space"
  },

  "VCAP_SERVICES" : {
    "auditlog" : [ {
      "name" : "auditlog-instance-name",
      "credentials" : {
        "password" : "password",
        "user" : "user",
        "url" : "https://host:port"
      }
    } ]
  }
}
```
