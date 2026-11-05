const fs = require("node:fs");
const tar = require("tar");
const zlib = require("node:zlib");
const path = require("node:path");
const { Readable } = require("node:stream");

const packageInfo = require(path.join(process.cwd(), "package.json"));
const version = packageInfo.version;

const binName = process.argv[2];
const os = process.argv[3] || process.platform;
const arch = process.argv[4] || process.arch;
const root = `https://github.com/SAP/${binName}/releases/download/v${version}/${binName}_${version}_`;

const requested = os + "-" + arch;
const current = process.platform + "-" + process.arch;
if (requested !== current ) {
  console.error("WARNING: Installing binaries for the requested platform (" + requested + ") instead of for the actual platform (" + current + ").")
}

const unpackedBinPath = path.join(process.cwd(), "unpacked_bin");
const config = {
  dirname: __dirname,
  binaries: [
    'mbt'
  ],
  urls: {
    'darwin-arm64': root + 'Darwin_arm64.tar.gz',
    'darwin-x64': root + 'Darwin_amd64.tar.gz',
    'linux-x64': root + 'Linux_amd64.tar.gz',
    'linux-arm64': root + 'Linux_arm64.tar.gz',
    'win32-x64': root + 'Windows_amd64.tar.gz'
  }
};
if (!fs.existsSync("bin")) {
  fs.mkdirSync("bin");
}

let binExt = "";
if (os == "win32") {
  binExt = ".exe";
}

const buildId = os + "-" + arch;
const url = config.urls[buildId];
if (!url) {
  throw new Error("No binaries are available for your platform: " + buildId);
}

function binstall(url, path, options) {
  return untgz(url, path, options);
}

function untgz(url, path, options) {
  options = options || {};

  const verbose = options.verbose;
  const verify = options.verify;

  return new Promise(function (resolve, reject) {
    const untar = tar
      .x({ cwd: path })
      .on("error", function (error) {
        reject("Error extracting " + url + " - " + error);
      })
      .on("end", function () {
        const successMessage = "Successfully downloaded and processed " + url;

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

    const gunzip = zlib.createGunzip().on("error", function (error) {
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

    fetch(url)
      .then((response) => {
        if (!response.ok) {
          throw new Error("HTTP error! status: " + response.status);
        }
        Readable.fromWeb(response.body).pipe(gunzip).pipe(untar);
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
