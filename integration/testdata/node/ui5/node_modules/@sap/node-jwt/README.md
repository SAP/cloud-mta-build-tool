@sap/node-jwt
=============

JSON Web Token (JWT) offline validation for Node with current binaries

This project contains the JWT binding for Node.js. It also includes the native libraries to run on Windows/Linux.
If you need another platforms, please write to the author.

# Platforms

Supported platforms: **Windows** | **Linux** | **MacOS**
Supported architectures: x64 on all platforms.
Please see also section dependencies for Node.js version.

#### Hello world

This standard example is from http://jwt.io 

```javascript
// you can either load a HMAC key for signatures with HSxxx
v.setSecret("secret"); // load HMAC key
v.setBase64Secret("c2VjcmV0"); // load a Base64 encoded HMAC key
// or you can load a PEM encoded X509 certificate for signatures with RSxxx
v.loadPEM("MIICozCCAYsCCAogFQcmCUcJMA0GCSqGSIb3DQEBCwUAMBQ...."); // load X509 public certificate OR public key for RSA signature validation
// check the token signature and validity
v.checkToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ");
console.log("Test JWT for Node.js");
if (v.getErrorDescription() !== "") {
   // error handling
   console.log("Error in JWT: " + v.getErrorDescription());
} else {
   // in case of success, retrieve the payload
   console.log("JWT Payload : " + v.getPayload());
}
```

# Getting started

From your project directory, run (see below for requirements):

```
$ var jwt = require('@sap/node-jwt');
```

Released versions
```
npm config set @sap:registry=https://npm.sap.com
npm install @sap/node-jwt
```


# Dependencies

* NodeJS v0.10.x is the minimum version, NodeJS 5.x is the current maximum, however due to missing binaries, there might be errors in using this project
* You dont need `node-gyp` or any compiler (e.g. Visual Studio on Windows). The source code and binding.gyp is part of this project in case of errors.
* If you run in error with generic node exceptions, please inform the author. The root cause can be missing jwt.node modules.


# Error situations

The standard error for signature operations is the situation, that the signature is not valid. This error is typical and you should handle
it carefully! and not as fatal error or assert.
If you think, it must work, but it does not, then you can trace the native functions.
SAPSSOEXT library allows you to set the environment variables:
* SAP_EXT_TRC to define a trace file in your file system
* SAP_EXT_TRL an integer 0 to 3

```
set SAP_EXT_TRC=stdout
set SAP_EXT_TRL=3
```

If you run your application in CloudFoundry or XSA then you can define environment variables with client command tool cf / xs, see
https://docs.run.pivotal.io/devguide/deploy-apps/manifest.html#env-block 

In cf landscapes you can then cf logs <your-app> and you will see trace from JWT validation

# Install via NPM

In order to configure the sap NPM registry you need to issue the following command:

```
npm config set @sap:registry=https://npm.sap.com
npm install @sap/node-jwt
```
