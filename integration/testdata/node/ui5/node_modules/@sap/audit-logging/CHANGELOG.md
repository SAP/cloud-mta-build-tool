# Change Log
All notable changes to this project will be documented in this file.

This project adheres to [Semantic Versioning](http://semver.org/).

The format is based on [Keep a Changelog](http://keepachangelog.com/).

## 2.2.4 - 2018-08-14

### Fixed
- Update dependencies.

## 2.2.3 - 2018-07-17

### Fixed
- Update request to v2.87.0.

## 2.2.2 - 2018-05-18

### Fixed
- Update request to v2.86.0.

## 2.2.1 - 2018-04-05

### Fixed
- Update dependencies.

## 2.2.0 - 2018-01-23

### Added
- npm-shrinkwrap.json.
- Audit log service v2 support.

### Fixed
- Documentation for data access messages.

## 2.1.1 - 2017-10-13

### Changed
- Dependencies' versions

## 2.1.0 - 2017-08-08

### Added
- Support for old and new values in data modification messages.
- Support for setting a tenant.
- Support for Node.js v8.
- Documentation improvements.

## 2.0.1 - 2017-05-30

### Fixed
- Improved debug messages.

## 2.0.0 - 2017-03-13

### Added
- The library requires the application to be bound to an instance of the **Audit log service**.
- Instantiating the library now requires a configuration object containing the `credentials` from the service binding to be provided (see README.md).
- `configurationChange` method is added to the API (see README.md).
- `updateConfigurationChange` method is added to the API (see README.md).

### Changed
- `read`, `update` and `securityMessage` methods have some changes to their argument lists (see README.md).

### Removed
- `delete` method is removed from the API.
- `create` method is removed from the API.
