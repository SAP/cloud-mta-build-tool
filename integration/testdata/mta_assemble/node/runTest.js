'use strict';

// start script for starting jest tests

const jest = require('jest');

// tests only
jest.run("--config jest.json");

// tests + coverage
// jest.run("--coverage --config jest.json");