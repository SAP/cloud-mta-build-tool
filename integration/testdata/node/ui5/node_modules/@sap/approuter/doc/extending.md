Extending Application Router
============================

<!-- toc -->

- [Basics](#basics)
- [Inject Custom Middleware](#inject-custom-middleware)
- [Application Router Extensions](#application-router-extensions)
- [Customize Command Line](#customize-command-line)
- [Dynamic Routing](#dynamic-routing)
- [State synchronization](#state-synchronization)
- [API Reference](#api-reference)
  * [approuter](#approuter)
    + [`approuter()`](#approuter)
    + [Event: 'login'](#event-login)
    + [Event: 'logout'](#event-logout)
    + [`first`](#first)
    + [`beforeRequestHandler`](#beforerequesthandler)
    + [`beforeErrorHandler`](#beforeerrorhandler)
    + [`start(options, callback)`](#startoptions-callback)
    + [`close(callback)`](#closecallback)
    + [`createRouterConfig(options, callback)`](#createrouterconfigoptions-callback)
    + [`resolveUaaConfig(request, uaaOptions, callback)`](#resolveuaaconfigrequest-uaaoptions-callback)
  * [Middleware Slot](#middleware-slot)
    + [`use(path, handler)`](#usepath-handler)

<!-- tocstop -->

## Basics
Insead of starting the application router directly, your application can have its own start script.
You can use the application router as a regular Node.js package.
```js
var approuter = require('@sap/approuter');

var ar = approuter();
ar.start();
```

## Inject Custom Middleware

The application router uses the [connect](https://github.com/senchalabs/connect)
framework.
You can reuse all _connect_ middlewares within the application router directly.
You can do this directly in your start script:
```js
var approuter = require('@sap/approuter');

var ar = approuter();

ar.beforeRequestHandler.use('/my-ext', function myMiddleware(req, res, next) {
  res.end('Request handled by my extension!');
});
ar.start();
```
__Tip:__ Name your middleware to improve troubleshooting.

The path argument is optional. You can also chain `use` calls.
```js
var approuter = require('@sap/approuter');
var morgan = require('morgan');

var ar = approuter();

ar.beforeRequestHandler
  .use(morgan('combined'))
  .use('/my-ext', function myMiddleware(req, res, next) {
    res.end('Request handled by my extension!');
  });
ar.start();
```

The application router defines the following slots where you can insert custom middleware:
* `first` - right after the _connect_ application is created, and before any
application router middleware.
At this point security checks are not performed yet.
__Tip:__ This is a good place for infrastructure logic like logging and monitoring.
* `beforeRequestHandler` - before standard application router request handling,
that is static resource serving or forwarding to destinations.
__Tip:__ This is a good place for custom REST API handling.
* `beforeErrorHandler` - before standard application router error handling.
__Tip:__ This is a good place to capture or customize error handling.

If your middleware does not complete the request processing, call `next`
to return control to the application router middleware:
```js
ar.beforeRequestHandler.use('/my-ext', function myMiddleware(req, res, next) {
  res.setHeader('x-my-ext', 'passed');
  next();
});
```

## Application Router Extensions

You can use application router extensions.

An extension is defined by an object with the following properties:
* `insertMiddleware` - describes the middleware provided by this extension
  * `first`, `beforeRequestHandler`, `beforeErrorHandler` - an array of middleware, where each one is either
    * a middleware function (invoked on all requests), or
    * an object with properties:
      * `path` - handle requests only for this path
      * `handler` - middleware function to invoke

Here is an example (my-ext.js):
```js
module.exports = {
  insertMiddleware: {
    first: [
      function logRequest(req, res, next) {
        console.log('Got request %s %s', req.method, req.url);
      }
    ],
    beforeRequestHandler: [
      {
        path: '/my-ext',
        handler: function myMiddleware(req, res, next) {
          res.end('Request handled by my extension!');
        }
      }
    ]
  }
};
```
You can use it in your start script like this:
```js
var approuter = require('@sap/approuter');

var ar = approuter();
ar.start({
  extensions: [
    require('./my-ext.js')
  ]
});
```

## Customize Command Line

By default the application router handles its command line parameters, but you can
customize that too.

An _approuter_ instance provides the property `cmdParser` that is a
[commander](https://github.com/tj/commander.js/) instance.
It is configured with the standard application router command line options.
There you can add custom options like this:
```js
var approuter = require('@sap/approuter');

var ar = approuter();

var params = ar.cmdParser
  // add here custom command line options if needed
  .option('-d, --dummy', 'A dummy option')
  .parse(process.argv);

console.log('Dummy option:', params.dummy);
```
To completely disable the command line option handling in the application router,
reset the following property:
```js
ar.cmdParser = false;
```

## Dynamic Routing

The application router can use a custom routing configuration for each request.

Here is an example:
```js
var approuter = require('@sap/approuter');

var ar = approuter();
ar.start({
  getRouterConfig: getRouterConfig
});

var customRouterConfig;
var options = {
  xsappConfig: {
    routes: [
      {
        source: '/service',
        destination: 'backend',
        scope: '$XSAPPNAME.viewer',
      }
    ]
  },
  destinations: [
    {
      name: 'backend',
      url: 'https://my.app.com',
      forwardAuthToken: true
    }
  ],
  xsappname: 'MYAPP'
};
ar.createRouterConfig(options, function(err, routerConfig) {
  if (err) {
    console.error(err);
  } else {
    customRouterConfig = routerConfig;
  }
});

function getRouterConfig(request, callback) {
  if (/\?custom-query/.test(request.url)) {
    callback(null, customRouterConfig);
  } else {
    callback(null, null); // use default router config
  }
}
```

## State synchronization

The application router can be scaled to run with multiple instances like any other application on Cloud Foundry.
Still application router instances are not aware of each other and there is no communication among them.
So if extensions introduce some state, they should take care to synchronize it across application router instances.

## API Reference

### approuter

#### `approuter()`

Creates a new instance of the application router.

#### Event: 'login'
Parameters:
* `session`
  * `id` - session id as a string

Emitted when a new user session is created.

#### Event: 'logout'
Parameters:
* `session`
  * `id` - session id as a string

Emitted when a user session has expired or a user has requested to log out.

#### `first`
A [Middleware Slot](#middleware-slot) before the first application router middleware

#### `beforeRequestHandler`
A [Middleware Slot](#middleware-slot) before the standard application router request handling

#### `beforeErrorHandler`
A [Middleware Slot](#middleware-slot) before the standard application router error handling

#### `start(options, callback)`

Starts the application router with the given options.

* `options` this argument is optional. If provided, it should be an object which can have any of the following properties:
  * `port` - a TCP port the application router will listen to (string, optional)
  * `workingDir` - the working directory for the application router,
  should contain the _xs-app.json_ file (string, optional)
  * `extensions` - an array of extensions, each one is an object as defined in
  [Application Router Extensions](#application-router-extensions) (optional)
  * `xsappConfig` - An object representing the content which is usually put in xs-app.json file.
  If this property is present it will take precedence over the content of xs-app.json. (optional)
  * `httpsOptions` - Options similar to [`https.createServer`](https://nodejs.org/api/https.html#https_https_createserver_options_requestlistener).
  If this property is present, application router will be started as an https server. (optional)
  * `getToken` - `function(request, callback)` Provide custom access token (optional)
    * `request` - Node request object
    * `callback` - `function(error, token)`
      * `error` - Error object in case of error
      * `token` - Access token to use in request to backend
  * `getRouterConfig` - `function(request, callback)` Provide custom routing configuration (optional)
    * `request` - Node request object
    * `callback` - `function(error, routerConfig)`
      * `error` - Error object in case of error
      * `routerConfig` - Custom routing configuration to use for given request.
      This object should be created via `createRouterConfig`.
      If `null` or `undefined`, default configuration will be used.
* `callback` - optional function with signature `callback(err)`.
It is invoked when the application router has started or an error has occurred.
If not provided and an error occurs (e.g. the port is busy), the application will abort.

#### `close(callback)`
Stops the application router.

* `callback` - optional function with signature `callback(err)`.
It is invoked when the application router has stopped or an error has occurred.

#### `createRouterConfig(options, callback)`
Prepares the routing configuration to be used by the application router.
As part of this, the application router validates the given options.
This function can be used at any point in runtime to create additional routing configurations.

**Note:** This function can be called only after `start` function.

* `options`
  * `xsappname` - Value to replace $XSAPPNAME placeholder in scope names.
  If not provided, it will be taken from UAA service binding. (optional)
  * `xsappConfig` - An object representing the content which is usually put in xs-app.json file.
**Note:** Only the following configurations are taken into account from this property (the rest are taken from the xs-app.json file):
`welcomeFile`, `logout.logoutEndpoint`, `logout.logoutPage`, `routes`, `websockets`, `errorPage`.
  * `destinations` - An array containing the configuration of the backend destinations.
  If not provided, it will be taken from `destinations` environment variable. (optional)
* `callback` - `function(error, routerConfig)`
  * `error` - Error object in case of error
  * `routerConfig` - Routing configuration to be passed to the callback of `getRouterConfig`.
  Approuter extensions should not access the content of this object.

#### `resolveUaaConfig(request, uaaOptions, callback)`

Calculates tenant-specific UAA configuration.

* `request` - node request object used to identify the tenant
* `uaaOptions` - UAA options as provided in service binding
* `callback` - `function(error, tenantUaaOptions)`
  * `error` - Error object in case of error
  * `tenantUaaOptions` - new UAA configuration with tenant-specific properties

### Middleware Slot

#### `use(path, handler)`
Inserts a request handling middleware in the current slot.

* `path` - handle only requests starting with this path (string, optional)
* `handler` - a middleware function to invoke (function, mandatory)

Returns `this` for chaining.
