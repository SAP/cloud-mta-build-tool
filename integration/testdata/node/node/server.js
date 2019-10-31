/*eslint no-console: 0*/
"use strict";

const express = require('express')
const app = express()
const port = 3000

app.get('/', (req, res) => res.send('Hello World!'))

app.listen(port, () => console.log(`Example app listening on port ${port}!`))

// var http = require("http");
// var port = process.env.PORT || 3000;
//
// http.createServer(function (req, res) {
//   res.writeHead(200, {"Content-Type": "text/plain"});
//   res.end("Hello World\n");
// }).listen(port);
//
// console.log("Server listening on port %d", port);
