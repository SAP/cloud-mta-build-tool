var binwrap = require('binwrap');
var path = require('path');

var packageInfo = require(path.join(__dirname, 'package.json'));
var version = '0.2.7'; //packageInfo.version;
var root = `https://github.com/SAP/cloud-mta-build-tool/releases/download/v${version}/cloud-mta-build-tool_${version}_`;

module.exports = binwrap({
  dirname: __dirname,
  binaries: [
    'mbt'
  ],
  urls: {
    'darwin-x64': root + 'Darwin_amd64.tar.gz',
    'linux-x64': root + 'Linux_amd64.tar.gz',
    'win32-x64': root + 'Windows_amd64.tar.gz'
  }
});
