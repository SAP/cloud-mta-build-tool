# base64-url

Base64 encode, decode, escape and unescape for URL applications.

<a href="https://nodei.co/npm/base64-url/"><img src="https://nodei.co/npm/base64-url.png?downloads=true"></a>


[![Build Status](https://travis-ci.org/joaquimserafim/base64-url.svg?branch=master)](https://travis-ci.org/joaquimserafim/base64-url)[![Coverage Status](https://coveralls.io/repos/github/joaquimserafim/base64-url/badge.svg)](https://coveralls.io/github/joaquimserafim/base64-url)[![ISC License](https://img.shields.io/badge/license-ISC-blue.svg?style=flat-square)](https://github.com/joaquimserafim/base64-url/blob/master/LICENSE)[![NodeJS](https://img.shields.io/badge/node-6.x.x-brightgreen.svg?style=flat-square)](https://github.com/joaquimserafim/base64-url/blob/master/package.json#L43)

[![JavaScript Style Guide](https://cdn.rawgit.com/feross/standard/master/badge.svg)](https://github.com/feross/standard)


## API

`const base64url = require('base64-url')`


### examples

```js

base64url.encode('Node.js is awesome.')
// returns Tm9kZS5qcyBpcyBhd2Vzb21lLg

base64url.decode('Tm9kZS5qcyBpcyBhd2Vzb21lLg')
// returns Node.js is awesome.

base64url.escape('This+is/goingto+escape==')
// returns This-is_goingto-escape

base64url.unescape('This-is_goingto-escape')
// returns This+is/goingto+escape==

//
// setting a different econding 
//

base64url.encode(string to encode, encoding)
base64url.decode(string to decode, encoding)

```


#### ISC License (ISC)

# Alternatives

- [base64url](https://github.com/brianloveswords/base64url)
