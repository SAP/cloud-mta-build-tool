# Change Log
All notable changes to this project will be documented in this file.

## 2.1.15 - 2018-08-13

- Update version of module request

## 2.1.14 - 2018-07-24

- Evaluate SAP_JWT_TRUST_ACL if trustedclientidsuffix is present but not matching

## 2.1.13 - 2018-07-18

- Update version of module request

## 2.1.12 - 2018-06-01

- Support for API methods getSubaccountId and getOrigin
- Mark API method getIdentityZone as deprecated

## 2.1.11 - 2018-05-18

- Update version of module request

## 2.1.10 - 2018-04-20

- Fixes for keycache

## 2.1.9 - 2018-04-18

- Update version of module @sap/node-jwt (1.4.8)
- Fixes for keycache
- Update version of module request

## 2.1.8 - 2018-03-14

- Support for API method getAppToken

## 2.1.7 - 2018-03-05

- Support for API method requestToken

## 2.1.6 - 2018-02-19

- Update version of module @sap/node-jwt

## 2.1.5 - 2018-02-07

- Update version of module request

## 2.1.4 - 2017-12-04

- Support new JWT structure (attribute location ext_cxt)
- First implementation for keycache

## 2.1.3 - 2017-11-29

- Support for API method getClientId

## 2.1.2 - 2017-10-23

- Support for API method getSubdomain

## 2.1.1 - 2017-10-09

- Update version of modules @sap/node-jwt, @sap/xsenv and debug

## 2.1.0 - 2017-07-06

- Support of API method requestTokenForClient
- Update version of module @sap/node-jwt

## 2.0.0 - 2017-06-26

- Removal of deprecated constructor method createSecurityContextCc
- Removal of API method method getUserInfo

## 1.3.0 - 2017-06-23

- Revert removal of API method method getUserInfo

## 1.2.0 - 2017-06-22

- Support for API methods getLogonName, getGivenName, getFamilyName, getEmail
- Removal of API method method getUserInfo
- Fix identity zone validation (only relevant for tenants created with SAP Cloud Cockpit)

## 1.1.1 - 2017-05-30
- Update version of dependent modules

## 1.1.0 - 2017-05-22
- Mark API method createSecurityContextCC as deprecated

## 1.0.4 - 2017-05-17

- Support for validation of XSUAA broker plan tokens
- Support for API methods getCloneServiceInstanceId and getAdditionalAuthAttribute
- Support for validation of XSUAA application plan tokens in arbitrary identity zones

## 1.0.3 - 2017-03-29

- Update version of dependent modules

## 1.0.2 - 2017-02-22

- Support for validation of SAML Bearer tokens

## 1.0.1 - 2017-02-02

- Support for client credentials tokens in JWT strategy

## 1.0.0 - 2017-01-25

- Introduction of scopeing, module name changed to @sap/xssec
