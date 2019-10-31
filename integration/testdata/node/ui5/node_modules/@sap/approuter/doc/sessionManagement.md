Extended Session Management
===========================

<!-- toc -->

- [Abstract](#abstract)
- [Session Lifecycle](#session-lifecycle)
- [Security](#security)
- [Data Privacy](#data-privacy)
- [API Reference](#api-reference)
- [Example](#example)
- [Performance](#performance)

<!-- tocstop -->

## Abstract

The application router uses a memory store as a session repository to provide
the best runtime performance. However, it is not persisted and it is not shared
across multiple instances of the application router. 

*__Note:__ The Limitations above do not prevent the application router from
being scaled out, since session stickiness is in place by default.* 

While it is good enough for most of the cases, it may be required to
provide a highly-available solution, which may be achieved by storing
the state (session) of the application router outside - in durable shared
storage.
To allow implementing these qualities, the application router exposes the
*extended session management* API described below.

## Session Lifecycle

The application router stores user agent sessions as JavaScript objects 
serialized to strings. It also stores the session timeout associated with
each session, which indicates the amount of time left until session
invalidation. 

### Initial Data

During the start of the application router, the internal session store is initiated.
It contains an empty list of sessions and their timeouts. The internal session
store is not available right after the application router instance is
created, but is available in the callback of `approuter.start` and all
the time afterwards until the application router is stopped.

In case an external session storage is used, the application router extension
should perform the following actions to synchronize the internal session
store with the external one:

- Load existing sessions from external storage 
- Start the application router
- Populate the application router's internal session store

### Read

A session identifier may be obtained from the request object `req.sessionID`.

On each request, the application router executes registered middlewares
in a certain order and the session is not available to all of them.

- First it passes the request to `approuter.first` middleware. 
  At this point, there is no session associated with
  the incoming request. 
- Afterwards, the application router checks if the user is authenticated, reads
  the relevant session from the internal session store and puts it into the request 
  context.
- Next, the application router passes a request to 
  `approuter.broforeRequestHandler`. At this point, the session object is
  available and associated with the incoming request.
- `approuter.beforeErrorHandler` also has access to session.  

### Login

When a user agent requests a resource, served via a route that requires
authentication, the application router will request the user agent to
pass authentication first (usually via redirect to XSUAA). At this point,
the application router does not create any session. Only after
the authentication process is finished, the application router creates a session,
stores it in the internal session storage and emits a `login` event.

### Update Session

Any changes made to the session are not stored in the internal session store 
immediately, but are accumulated to make a bulk update after the end of the response.
While the request is passed through the chain of middlewares, the session object
may be modified. Also, when the backend access token is close to expire,
the application router may trigger the refresh backend token flow. In both cases,
the actual update of the internal session store is done later on, outside of
the request context.

### Timeout

There is a time-based job in the application router that basis outside 
the request context and destroys sessions with an elapsed timeout.

Each time the application router reads a session from the session store,
the timeout of this session is reset to the initial value that may be retrieved
using the [`getDefaultSessionTimeout()`](#sessionstoregetdefaultsessiontimeout)
API.

### Logout

When a user agent requests a URL defined as the `logoutEndpoint` in the 
`xs-app.json` file, a central logout process takes place. As part of this
process, the application router emits a `logout` event. More detailed
information about the central logout may be found in 
[README.md](../README.md) 

## Security

The application router uses session secret to sign session cookies and
prevent tampering. The session secret, by default, is generated using
a random sequence of bytes at the startup of the application router. It is
different for each instance and changed on each restart of the same
instance.

Using the default session secret generation mechanism for highly available
application routers may cause issues in the following scenarios:

- The user agent is authenticated and the session is stored in a session store.
  The application router is restarted (due to internal error or triggered
  by platform) and a new session secret is generated. The authenticated user
  agent makes a request, which contains the session cookies. However, the cookies are 
  signed using another secret and the application router ignores them.
- The user agent is authenticated and the session is stored in the session store.
  The application router instance is unavailable. The authenticated user agent 
  makes a request to the application router and the request contains the session
  cookies. The load balancer forwards the request to another instance of 
  the application router. However, cookies are signed using another secret and
  the application router ignores them.

In both scenarios, the session in the store is no longer accessible, the cookies
sent by the user agent are redundant, and the user agent will be requested to
pass authentication once again.

To avoid the issues described above, the extension that implements the extended session
management mechanism, should make sure to implement the `getSessionSecret` hook.

```js
var ar = AppRouter();

ar.start({
  getSessionSecret: function () {
    return 'CUSTOM_PERSISTED_SESSION_SECRET';
  },
  ...
});
```

It is recommended to have at least 128 characters in the string that replaces 
`CUSTOM_PERSISTED_SESSION_SECRET`.

## Data Privacy

The user agent session potentially contains personal data. By implementing
the custom session management behaviour, you take the responsibility to be
compliant with all personal data protection laws and regulations
(e.g. [GDPR](https://en.wikipedia.org/wiki/General_Data_Protection_Regulation))
that may be applied in the regions, where the application will be used.

## API Reference

### Methods

#### approuter.start(options)

* `options`
  * `getSessionSecret` - returns the session secret to be used
    by the application router for the signing of the session cookies. 

#### approuter.getSessionStore()

returns `SessionStore` instance.

#### sessionStore.getDefaultSessionTimeout()

returns the default session timeout in minutes.

#### sessionStore.getSessionTimeout(sessionId, callback)

* `sessionId` - an unsigned session identifier
* `callback` - `function(error, session)` a function that is called 
   when the session object is retrieved from the internal session 
   storage of the application router.
  * `error` - an error object in case of an error, otherwise `null`
  * `timeout` - time, in minutes, until the session times out

#### sessionStore.get(sessionId, callback)

* `sessionId` - an unsigned session identifier
* `callback` - `function(error, session)` a function that is called 
   when the session object is retrieved from the internal session 
   storage of the application router.
  * `error` - an error object in case of an error, otherwise `null`
  * `session` - the session object
    * `id` - session identifier, immutable

#### sessionStore.set(sessionId, sessionString, timeout, callback)

* `sessionId` - an unsigned session identifier
* `sessionString` - a session object serialized to string
* `timeout` - a timestamp in milliseconds, after which the session should be 
   automatically invalidated
* `callback` - a function that is called after the session is saved in the
   internal session storage of the application router 

#### sessionStore.update(sessionId, callback, resetTimeout)

* `sessionId` - an unsigned session identifier
* `callback` - `function(currentSession)` function, which returns 
  session object. Callback function may modify and return current 
  session object or create and return brand new session object
  * `currentSession` - current session object
* `resetTimeout` - a boolean that indicates whether to reset the session timeout


#### sessionStore.destroy(sessionId, callback)

* `sessionId` - an unsigned session identifier
* `callback` - a function that is called after the session is destroyed in
  the internal session storage of the application router 

### Events

Extension may subscribe to application router events using the standard
[`EventEmitter`](https://nodejs.org/api/events.html) API.

```js
var ar = AppRouter();

ar.on('someEvent', function handler() {
  // Handle event
});
```

#### `login`

Emitted when user agent is authenticated.

Parameters:
* `session` - session object
  * `id` - session identifier, immutable

#### `logout`

Emitted when a user agent session is going to be terminated in
the internal session store of the application router. Emitted either when
the user agent session is timed-out or when `logoutEndpoint` was requested. 

*__Note:__ Central logout is an asynchronous process. The order in which
the backend and the application router sessions are invalidated, is not
guaranteed.*

Parameters:
* `session` - session object
  * `id` - session identifier, immutable

## Example

There may be many various options, how the application router extension
decides to store sessions exposed via the session management API. The example
below assumes a `SessionDataAccessObject` to be implemented by the extension
developer and to have the following API:

### Methods:

* `sessionDataAccessObject.create` - `function(session, timeout)`
* `sessionDataAccessObject.update` - `function(sessionId, timeout)`
* `sessionDataAccessObject.delete` - `function(sessionId)`
* `sessionDataAccessObject.load` - `function()`

### Events:

#### `create`

Parameters: 

  * `sessionId` - session identifier
  * `session` - session object serialized to string
  * `timeout` - timestamp, when session should expire
  * `callback` - function to be called after session is stored in
  internal session storage

#### `update`

Parameters:

  * `sessionId` - session identifier
  * `session` - session object serialized to string
  * `timeout` - timestamp, when session should be expired
  * `callback` - function to be called after session is stored in
  internal session storage

#### `delete`

Parameters:

  * `sessionId` - session identifier

#### `load`

Parameters:

  * `sessions[]` - array of objects
    * `id` - session identifier
    * `session` - session object serialized to string
    * `timeout` - timestamp, when session should expire

```js
var ar = new require('@sap/approuter')();
var dao = new SessionDataAccessObject();

dao.on('load', function (data) {

    ar.start({
        getSessionSecret: function getSessionSecret() {
            return process.env.SESSION_SECRET;
        }
    }, function() {
        var store = ar.getSessionStore();
        var defaultTimeout = store.getDefaultSessionTimeout();

        // AppRouter -> Persistence
        ar.on('login', function(session) {
            dao.create(session, defaultTimeout);
        });
        ar.on('update', function(sessionId, timeout) {
            dao.update(sessionId, timeout);
        });
        ar.on('logout', function(sessionId) {
            dao.delete(sessionId);
        });

        // Load Initial Data
        data.forEach(function(item) {
            store.set(item.id, item.session, item.timeout);
        });

        // Persistence -> AppRouter
        dao.on('create', store.set);
        dao.on('update', store.set);
        dao.on('delete', store.destroy);
    });

});

dao.load();
```

## Performance

*__Note:__ The `update` event of the application router may be potentially
triggered thousands of times a second. It is recommended to throttle or
debounce calls to the external storage to reduce network and CPU
consumption.*

Here is an example of a throttled `dao.update()`, where the latest change
will be persisted in the external storage no more than once in `500ms` for
the same session.

```js
// Throttled update
update(sessionId, timeout) {
    var dao = this;
    var sessionStore = this._sessionStore;
    if(typeof timeout === 'undefined') {
        if (!this.updateTimers[sessionId]) {
            this.updateTimers[sessionId] = setTimeout(function() {
                dao.updateTimers[sessionId] = null;
            }, 500);
            sessionStore.get(sessionId, function(err, session) {
                dao._saveSession(sessionId, session)
            });
        }
    } else {
        if (!this.timeoutTimers[sessionId]) {
            this.timeoutTimers[sessionId] = setTimeout(function() {
                dao.timeoutTimers[sessionId] = null;
            }, 500);
            sessionStore.getSessionTimeout(sessionId, function(err, timeout) {
                dao._saveTimeout(sessionId, timeout)
            });
        }
    }
}
```

And here is an example of a debounced `dao.update()`, where the latest 
change will be persisted in the external storage only if there were no other
changes during the last `500ms` for the same session.

```js
// Debounced update
update(sessionId, timeout) {
    var dao = this;
    var sessionStore = this._sessionStore;
    if(typeof timeout === 'undefined') {
        if (this.updateTimers[sessionId]) {
            clearTimeout(this.updateTimers[sessionId]);
        }
        this.updateTimers[sessionId] = setTimeout(function() {
            sessionStore.get(sessionId, function(err, session) {
                dao._saveSession(sessionId, session)
            });
        }, 500);
    } else {
        if (this.timeoutTimers[sessionId]) {
            clearTimeout(this.timeoutTimers[sessionId]);
        }
        this.timeoutTimers[sessionId] = setTimeout(function() {
            sessionStore.getSessionTimeout(sessionId, function(err, timeout) {
                dao._saveTimeout(sessionId, timeout)
            });
        }, 500);
    }
}
```

To understand the difference between throttling and debouncing, let's
consider an example, where requests for the same session come every 
`100ms` for `1sec`. In case of `500ms` debouncing, changes will be 
persisted one time. In case of `500ms` throttling, changes will be
persisted two times. Without any optimisation, changes will be
persisted ten times.
