# request-stats

[![Build status](https://travis-ci.org/watson/request-stats.svg?branch=master)](https://travis-ci.org/watson/request-stats)
[![js-standard-style](https://img.shields.io/badge/code%20style-standard-brightgreen.svg?style=flat)](https://github.com/feross/standard)

Get stats on your Node.js HTTP server requests.

Emits two events:

- `request` when ever a request starts: Passes a [Request object](#request-object) that can later be used to [query for the progress](#progress) of a long running request
- `complete` when ever a request completes: Passes a [stats object](#oncomplete-callback) containing the overall stats for the entire HTTP request

## Installation

```
npm install request-stats --save
```

## Example usage

Get stats for each completed HTTP request:

```javascript
var requestStats = require('request-stats')
var server = http.createServer(...)

requestStats(server, function (stats) {
  // this function will be called every time a request to the server completes
  console.log(stats)
})
```

Get periodic stats for long running requests:

```javascript
var server = http.createServer(...)

var stats = requestStats(server)

stats.on('request', function (req) {
  // evey second, print stats
  var interval = setInterval(function () {
    var progress = req.progress()
    console.log(progress)
    if (progress.completed) clearInterval(interval)
  }, 1000)
})
```

## API

### Constructor

#### `requestStats(server[, callback])`

Attach request-stats to a HTTP server.

Initialize request-stats with an instance a HTTP server. Returns a
StatsEmitter object. Optionally provide a callback which will be called
for each completed HTTP request with a stats object (see stats object
details below).

If no callback is provided, you can later attach a listener on the
"complete" event.

#### `requestStats(req, res[, callback])`

Attach request-stats to a single HTTP request.

Initialize request-stats with an instance a HTTP request and response.
Returns a StatsEmitter object. Optionally provide a callback which will
be called with a stats object when the HTTP request completes (see stats
object details below).

If no callback is provided, you can later attach a listener on the
"complete" event.

### StatsEmitter object

#### `.on('complete', callback)`

Calls the callback function with a stats object when a HTTP request
completes:

```javascript
{
  ok: true,           // `true` if the connection was closed correctly and `false` otherwise
  time: 0,            // The milliseconds it took to serve the request
  req: {
    bytes: 0,         // Number of bytes sent by the client
    headers: { ... }, // The headers sent by the client
    method: 'POST',   // The HTTP method used by the client
    path: '...',      // The path part of the request URL
    ip: '...',        // The remote ip
    raw: [Object]     // The original `http.IncomingMessage` object
  },
  res: {
    bytes: 0,         // Number of bytes sent back to the client
    headers: { ... }, // The headers sent back to the client
    status: 200,      // The HTTP status code returned to the client
    raw: [Object]     // The original `http.ServerResponse` object
  }
}
```

#### `.on('request', callback)`

Calls the callback function with a special [Request
object](#request-object) when a HTTP request is made to the server.

### Request object

The Request object should not be confused with the Node.js
[http.IncomingMessage](http://nodejs.org/api/http.html#http_http_incomingmessage)
object. The request-stats Request object provides only a single
but powerfull function:

#### `.progress()`

Returns a progress object if called while a HTTP request is in progress.
If called multiple times, the returned progress object will contain the
delta of the previous call to `.progress()`.

```javascript
{
  completed: false, // `false` if the request is still in progress
  time: 0,          // The total time the reuqest have been in progress
  timeDelta: 0,     // The time since previous call to .progress()
  req: {
    bytes: 0,       // Total bytes received
    bytesDelta: 0,  // Bytes received since previous call to .progress()
    speed: 0,       // Bytes per second calculated since previous call to .progress()
    bytesLeft: 0,   // If the request contains a Content-Size header
    timeLeft: 0     // If the request contains a Content-Size header
  },
  res: {
    bytes: 0,       // Total bytes send back to the client
    bytesDelta: 0,  // Bytes sent back to the client since previous call to .progress()
    speed: 0        // Bytes per second calculated since previous call to .progress()
  }
}
```

## Acknowledgement

Thanks to [mafintosh](https://github.com/mafintosh) for coming up with
the initial concept and pointing me in the right direction.

## License

MIT
