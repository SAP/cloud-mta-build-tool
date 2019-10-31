@sap/xssec: XS Advanced Container Security API for node.js
==========================================================

## XS Advanced Authentication Primer

Authentication for node applications in XS Advanced relies on a special usage of the OAuth 2.0 protocol, which is based on central authentication at the UAA server that then vouches for the authenticated user's identity via a so-called OAuth Access Token. The current implementation uses as access token a JSON web token (JWT), which is a signed text-based token following the JSON syntax.

Normally, your node application will consist of several parts, that appear as separate applications in your manifest file, e.g. one application part that is responsible for the HANA database content, one application part for your application logic written e.g. in node.js (this is the one that can make use of this XS Advanced Container Security API for node.js), and finally one application part that is responsible for the UI layer (this is the one that may make use of the application router functionality). The latter two applications (the application logic in node.js and the application router) should be bound to one and the same UAA service instance. This has the effect, that these two parts can use the same OAuth client credentials.

When your business users access your application UI with their broser, the application router redirects the browser to the UAA where your business users need to authenticate. After successful authentication, the UAA sends - again via the business user's browser - an OAuth authorization code back to the application router. Now the application router sends this authorization code directly (not via the browser) to the UAA to exchange it into an OAuth access token. If the access token is obtained successfully, the business user has logged on to the UI part of your application already. In order to enable your UI to pass this authentication on to the node.js application part, you need to ensure that the destination to your node.js application part is configured such that the access token is actually sent to the node.js part ("forwardAuthToken": true).

In order to authenticate this request, which arrives at the node.js backend, sap-xssec offers two mechanisms: Firstly, you can use the XS Advanced Container Security API directly to validate the access token. Secondly, you can make use of the passport strategy that is contained in module sap-xssec as another convenient way how to handle the access token. In the following, both options are described followed by the sap-xssec API description.

sap-xssec offers an offline validation of the access token, which requires no additional call to the UAA. The trust for this offline validation is created by binding the XS UAA service instance to your application. Inside the credentials section in the environment variable VCAP_SERVICES, the key for validation of tokens is included. By default, the offline validation check will only accept tokens intended for the same OAuth2 client in the same UAA identity zone. This makes sense and will cover the vast majority of use cases. However, if an application absolutely wants to consume token that were issued for either different OAuth2 clients or different identity zones, an Access Control List (ACL) entry for this can be specified in an environment variable named SAP_JWT_TRUST_ACL. The name of the OAuth client is sb-<xsappname from xs-security.json>
The content is a JSON String, containing an array of identity zones and OAuth2 clients. To trust any OAuth2 client and/or identity zones, an * can be used. For OP, identity zones are not used and value for the identity zone is uaa.

```JSON
SAP_JWT_TRUST_ACL: [ {"clientid":"<client-id of the OAuth2 client>","identityzone":"<identity zone>"},...]
```

If you want to enable another (foreign) application to use some of your application's scopes, you can add a ```granted-apps``` marker to your scope in the ```xs-security.json``` file (as in the following example). The value of the marker is a list of applications that is allowed to request a token with the denoted scope.

```JSON
{
  "xsappname"     : "sample-leave-request-app",
  "description"   : "This sample application demos leave requests",
  "scopes"        : [ { "name"                : "$XSAPPNAME.createLR",
                        "description"         : "create leave requests" },
                      { "name"                : "$XSAPPNAME.approveLR",
                        "description"         : "approve leave requests",
                        "granted-apps"        : ["MobileApprovals"] }
                    ],
  "attributes"    : [ { "name"                : "costcenter",
                        "description"         : "costcenter",
                        "valueType"           : "string"
                    } ],
  "role-templates": [ { "name"                : "employee",
                        "description"         : "Role for creating leave requests",
                        "scope-references"    : [ "$XSAPPNAME.createLR","JobScheduler.scheduleJobs" ],
                        "attribute-references": [ "costcenter"] },
                      { "name"                : "manager",
                        "description"         : "Role for creating and approving leave requests",
                        "scope-references"    : [ "$XSAPPNAME.createLR","$XSAPPNAME.approveLR","JobScheduler.scheduleJobs" ],
                        "attribute-references": [ "costcenter" ] }
                    ]
}
```

## Usage of the XS Advanced Container Security API in your node.js Application

In order to use the capabilities of the XS Advanced container security API,  add the module "sap-xssec" to the dependencies section of your application's package.json.

To enable tracing, you can set the environment variable DEBUG as follows: `DEBUG=xssec:*`.

### Direct Usage with existing Access Token

For the usage of the XS Advanced Container Security API it is necessary to pass a JWT token. If you have such a token, you may use the API as follows. The examples below rely on users and credentials that you should substitute with the ones in your context. The code below is based on version v0.0.9 (if you use another version, the coding might differ).

The typical use case for calling this API lies from within a container when an HTTP request is received. In an authorization header (with keyword `bearer`) an access token is contained already. You can remove the prefix `bearer` and pass the remaining string (just as in the following example as `access_token`) to the API.

```js
xssec.createSecurityContext(access_token, xsenv.getServices({ uaa: 'uaa' }).uaa, function(error, securityContext) {
    if (error) {
        console.log('Security Context creation failed');
        return;
    }
    console.log('Security Context created successfully');
});
```

Note that the example above uses module `xsenv` to retrieve the configuration of the default services (which are read from environment variable `VCAP_SERVICES` or if not set, from the default configuration file). However, it passes only the required `uaa` configuration to the method `createSecurityContext`. As default the UAA configuration is searched with tag `xsuaa` by `xsenv`. For details we refer to module @sap/xsenv. The `xsenv` documentation also helps if you want to provide the credentials from e.g. a user provided service.

The creation function `xssec.createSecurityContext` is to be used for an end-user token (e.g. for grant_type `password` or grant_type `authorization_code`) where user information is expected to be available within the token and thus within the security context.

`createSecurityContext` also accepts a token of grant_type `client_credentials`. This leads to the creation of a limited SecurityContext where certain functions are not available. For more details please consult the API description below or your documentation.


### Usage with Passport Strategy

If you use [express](https://www.npmjs.com/package/express) and [passport](https://www.npmjs.com/package/passport), you can easily plug a ready-made authentication strategy.

```js
var express = require('express');
var passport = require('passport');
var JWTStrategy = require('@sap/xssec').JWTStrategy;
var xsenv = require('@sap/xsenv');

...

var app = express();

...

passport.use(new JWTStrategy(xsenv.getServices({uaa:{tag:'xsuaa'}}).uaa));

app.use(passport.initialize());
app.use(passport.authenticate('JWT', { session: false }));
```

If JWT token is present in the request and it is successfully verified, following objects are created:
* request.user - according to [User Profile](http://passportjs.org/docs/profile) convention
  * id
  * name
    * givenName
    * familyName
  * emails `[ { value: <email> } ]`
* request.authInfo - the [Security Context](#api-description)

If the `client_credentials` JWT token is present in the request and it is successfully verified, following objects are created:
* request.user - empty object
* request.authInfo - the [Security Context](#api-description)

#### Session

It is recommended to _disable the session_ as in the example above.
In XSA each request comes with a JWT token so it is authenticated explicitly and identifies the user.
If you still need the session, you can enable it but then you should also implement [user serialization/deserialization](http://passportjs.org/guide/configure/) and some sort of [session persistency](https://github.com/expressjs/session).

### Test Usage without having an Access Token

For test purposes, you may retrieve the token for a certain user (whose credentials you know) from the UAA as in the following code-snippet.

```js
var http = require("http");
var xssec = require("@sap/xssec");
var xsenv = require('@sap/xsenv');
var request = require('request');

var uaaService = xsenv.getServices( { uaa: 'uaa' } ).uaa;
var testService = xsenv.getServices( { test : { label : 'test' } } ).test;
process.env.XSAPPNAME = testService.test.xsappname;

var options = {
    url : uaaService.url + '/oauth/token?client_id=' + uaaService.clientid
            + '&grant_type=password&username=' + testService.userid + '&password='
            + testService.usersecret
};
request.get(
    options,
    function(error, response, body) {
        if (error || response.statusCode !== 200) {
            console.log('Token request failed');
            return;
        }
        var json = null;
        try {
            json = JSON.parse(body);
        } catch (e) {
        	return callback(e);
        }
        xssec.createSecurityContext(json.access_token, uaaService, function(error, securityContext) {
            if (error) {
                console.log('Security Context creation failed');
                return;
            }
            console.log('Security Context created successfully');
        });
    }
).auth(uaaService.clientid, uaaService.clientsecret, false);
```
Note that this example assumes additional test configuration in the file `default-services.json`.

```json
{
  "uaa": {
    "url"             : "<UAA URL>",
    "clientid"        : "<your application's OAuth client id>",
    "clientsecret"    : "<your application's OAuth client secret>",
    "xsappname"       : "<your application's name>",
    "identityzone"    : "<desired UAA identity zone>",
    "tags"            : ["xsuaa"],
    "verificationkey" : "<verification key for offline validation>"
  },
  "test": {
    "userid"          : "marissa",
    "usersecret"      : "koala"
  }
}
```

## API Description

### createSecurityContext

This function creates the Security Context by validating the received access token against credentials put into the application's environment via the UAA service binding.

Usually, the received token must be intended for the current application. More clearly, the OAuth client id in the access token needs to be equal to the OAuth client id of the application (from the application's environment).

However, there are some use cases, when a "foreign" token could be accepted although it was not intended for the current application. If you want to enable other applications calling your application backend directly, you can specify in your xs-security.json file an access control list (ACL) entry and declare which OAuth client from which Identity Zone may call your backend.

Parameters:

* `access token` ... the access token as received from UAA in the "authorization Bearer" HTTP header
* `config` ... a structure with mandatory elements url, clientid and clientsecret
* `callback(error, securityContext)`

### getLogonName

not available for tokens of grant_type `client_credentials`, returns the logon name

### getGivenName

not available for tokens of grant_type `client_credentials`, returns the given name

### getFamilyName

not available for tokens of grant_type `client_credentials`, returns the family name

### getEmail

not available for tokens of grant_type `client_credentials`, returns the email

### getOrigin

* returns the user origin. The origin is an alias that refers to a user store in which the user is persisted. For example, users that are authenticated by the UAA itself with a username/password combination have their origin set to the value uaa.

### checkLocalScope

checks a scope that is published by the current application in the xs-security.json file.

Parameters:

* `scope` ... the scope whose existence is checked against the available scopes of the current user. Here, no prefix is required.
* returns `true` if the scope is contained in the user's scopes, `false` otherwise

### checkScope

checks a scope that is published by an application.

Parameters:

* `scope` ... the scope whose existence is checked against the available scopes of the current user.  Here, the prefix is required, thus the scope string is "globally unique".
* returns `true` if the scope is contained in the user's scopes, `false` otherwise

### getToken (obsolete, use getHdbToken or getAppToken)

Parameters:

* `namespace` ... Tokens can eventually be used in different contexts, e.g. to access the HANA database, to access another XS2-based service such as the Job Scheduler, or even to access other applications/containers. To differentiate between these use cases, the `namespace` is used. In `lib/constants.js` we define supported namespaces (e.g. `SYSTEM`).
* `name` ... The name is used to differentiate between tokens in a given namespace, e.g. `HDB` for HANA database or `JOBSCHEDULER` for the job scheduler. These names are also defined in the file `lib/constants.js`.
* returns a token that can be used e.g. for contacting the HANA database. If the token, that the security context has been instantiated with, is a foreign token (meaning that the OAuth client contained in the token and the OAuth client of the current application do not match), `null` is returned instead of a token.

### getAppToken

* returns the token of the application that can be used e.g. for token forwarding to another app.

### getHdbToken

* returns a token that can be used for contacting the HANA database. If the token, that the security context has been instantiated with, is a foreign token (meaning that the OAuth client contained in the token and the OAuth client of the current application do not match), `null` is returned instead of a token.

### requestToken

Requests a token based on the given type. The type can be `constants.TYPE_USER_TOKEN` or `constants.TYPE_CLIENT_CREDENTIALS_TOKEN`. Prerequisite for the former is that the requesting client has `grant_type=user_token` and that the current user token includes the scope `uaa.user`.

* `serviceCredentials` ... the credentials of the service as JSON object. The attributes `clientid`, `clientsecret` and `url` (UAA) are mandatory. Note that the subdomain of the `url` will be adapted to the subdomain of the application token if necessary.
* `type` ... allowed values are `constants.TYPE_USER_TOKEN` and `constants.TYPE_CLIENT_CREDENTIALS_TOKEN`
* `additionalAttributes` ... the attributes that should be included into the JWT token as JSON object (key-value list), e.g. `{"attr1" : "value1", "attr2" : "value2"}` 
* `cb(error, token)` ... callback function

### requestTokenForClient (obsolete, use requestToken instead)

Requests a token with `grant_type=user_token` from another client. Prerequisite is that the requesting client has `grant_type=user_token` and that the current user token includes the scope `uaa.user`.

Parameters:

* `serviceCredentials` ... the credentials of the service as JSON object. The attributes `clientid`, `clientsecret` and `url` (UAA) are mandatory.
* `scopes` ... comma-separated list of requested scopes for the token, e.g. `app.scope1,app.scope2`. If null, all scopes are granted. Note that $XSAPPNAME is not supported as part of the scope names.
* `cb(error, token)` ... callback function

### hasAttributes

not available for tokens of grant_type `client_credentials`.

* returns `true` if the token contains any xs user attributes, `false` otherwise.

### getAttribute

not available for tokens of grant_type `client_credentials`.

Parameters:

* `name` ... The name of the attribute that is requested.
* returns the attribute exactly as it is contained in the access token. If no attribute with the given name is contained in the access token, `null` is returned. If the token, that the security context has been instantiated with, is a foreign token (meaning that the OAuth client contained in the token and the OAuth client of the current application do not match), `null` is returned regardless of whether the requested attribute is contained in the token or not.

### getAdditionalAuthAttribute

Parameters:

* `name` ... The name of the additional authentication attribute that is requested.
* returns the additional authentication attribute exactly as it is contained in the access token. If no attribute with the given name is contained in the access token, `null` is returned. Note that additional authentication attributes are also returned in foreign mode (in contrast to getAttribute).

### isInForeignMode

* returns `true` if the token, that the security context has been instantiated with, is a foreign token that was not originally issued for the current application, `false` otherwise.

### getSubdomain

* returns the subdomain that the access token has been issued for.

### getClientId

* returns the client id that the access token has been issued for.

### getIdentityZone (obsolete, use getSubaccountId instead)

* returns the identity zone that the access token has been issued for.

### getSubaccountId

* returns the subaccount id that the access token has been issued for.

### getExpirationDate

* returns the expiration date of the access token as javascript Date object.

### getCloneServiceInstanceId

* returns the service instance id of the clone if the XSUAA broker plan is used.

### getGrantType

* returns the grant type of the JWT token, e.g. `authorization_code`, `password`, `client_credentials` or `urn:ietf:params:oauth:grant-type:saml2-bearer`.
