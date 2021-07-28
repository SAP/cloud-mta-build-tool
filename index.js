var binwrap = require('binwrap');
var path = require('path');

var packageInfo = require(path.join(__dirname, 'package.json'));
var version = packageInfo.version;
var root = (process.env.XMAKE_IMPORT_COMMON_0 ? `${process.env.XMAKE_IMPORT_COMMON_0}/com/github/sap/cloud-mta-build-tool/${version}/cloud-mta-build-tool-${version}-` : `https://github.com/SAP/cloud-mta-build-tool/releases/download/v${version}/cloud-mta-build-tool_${version}_`);

module.exports = binwrap({
  dirname: __dirname,
  binaries: [
    'mbt'
  ],
  urls: {
    'darwin-arm64': root + 'Darwin_arm64.tar.gz',
    'darwin-x64': root + 'Darwin_amd64.tar.gz',
    'linux-x64': root + 'Linux_amd64.tar.gz',
    'win32-x64': root + 'Windows_amd64.tar.gz'
  }
});
