var fs = require("fs");
var axios = require("axios");
var tar = require("tar");
var zlib = require("zlib");
var unzip = require("unzip-stream");
var path = require("path");


var packageInfo = require(path.join(process.cwd(), "package.json"));
var version = packageInfo.version;

var binName = process.argv[2];
var os = process.argv[3] || process.platform;
var arch = process.argv[4] || process.arch;
var root = `https://github.com/SAP/${binName}/releases/download/v${version}/${binName}_${version}_`;


var requested = os + "-" + arch;
var current = process.platform + "-" + process.arch;
if (requested !== current ) {
  console.error("WARNING: Installing binaries for the requested platform (" + requested + ") instead of for the actual platform (" + current + ").")
}

var unpackedBinPath = path.join(process.cwd(), "unpacked_bin");
var config = {
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
};
if (!fs.existsSync("bin")) {
  fs.mkdirSync("bin");
}

var binExt = "";
if (os == "win32") {
  binExt = ".exe";
}

var buildId = os + "-" + arch;
var url = config.urls[buildId];
if (!url) {
  throw new Error("No binaries are available for your platform: " + buildId);
}
function binstall(url, path, options) {
  if (url.endsWith(".zip")) {
    return unzipUrl(url, path, options);
  } else {
    return untgz(url, path, options);
  }
}

function untgz(url, path, options) {
  options = options || {};

  var verbose = options.verbose;
  var verify = options.verify;

  return new Promise(function (resolve, reject) {
    var untar = tar
      .x({ cwd: path })
      .on("error", function (error) {
        reject("Error extracting " + url + " - " + error);
      })
      .on("end", function () {
        var successMessage = "Successfully downloaded and processed " + url;

        if (verify) {
          verifyContents(verify)
            .then(function () {
              resolve(successMessage);
            })
            .catch(reject);
        } else {
          resolve(successMessage);
        }
      });

    var gunzip = zlib.createGunzip().on("error", function (error) {
      reject("Error decompressing " + url + " " + error);
    });

    try {
      fs.mkdirSync(path);
    } catch (error) {
      if (error.code !== "EEXIST") throw error;
    }

    if (verbose) {
      console.log("Downloading binaries from " + url);
    }

    axios
      .get(url, { responseType: "stream" })
      .then((response) => {
        response.data.pipe(gunzip).pipe(untar);
      })
      .catch((error) => {
        if (verbose) {
          console.error(error);
        } else {
          console.error(error.message);
        }
      });
  });
}

function unzipUrl(url, path, options) {
  options = options || {};

  var verbose = options.verbose;
  var verify = options.verify;

  return new Promise(function (resolve, reject) {
    var writeStream = unzip
      .Extract({ path: path })
      .on("error", function (error) {
        reject("Error extracting " + url + " - " + error);
      })
      .on("entry", function (entry) {
        console.log("Entry: " + entry.path);
      })
      .on("close", function () {
        var successMessage = "Successfully downloaded and processed " + url;

        if (verify) {
          verifyContents(verify)
            .then(function () {
              resolve(successMessage);
            })
            .catch(reject);
        } else {
          resolve(successMessage);
        }
      });

    if (verbose) {
      console.log("Downloading binaries from " + url);
    }

    axios
      .get(url, { responseType: "stream" })
      .then((response) => {
        response.data.pipe(writeStream);
      })
      .catch((error) => {
        if (verbose) {
          console.error(error);
        } else {
          console.error(error.message);
        }
      });
  });
}

function verifyContents(files) {
  return Promise.all(
    files.map(function (filePath) {
      return new Promise(function (resolve, reject) {
        fs.stat(filePath, function (err, stats) {
          if (err) {
            reject(filePath + " was not found.");
          } else if (!stats.isFile()) {
            reject(filePath + " was not a file.");
          } else {
            resolve();
          }
        });
      });
    })
  );
}

binstall(url, unpackedBinPath).then(function() {
  config.binaries.forEach(function(bin) {
    fs.chmodSync(path.join(unpackedBinPath, bin + binExt), "755");
  });
}).then(function(result) {
  process.exit(0);
}, function(result) {
  console.error("ERR", result);
  process.exit(1);
});
