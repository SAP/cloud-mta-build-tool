#!/usr/bin/env node
const micromatch = require('micromatch');
const { program } = require('commander');
const path = require('path');
const util = require('util');
const fs = require('fs');
const { log } = require('console');

program
  .version('1.0.0')
  .description('Micromatch CLI Wrapper')

const matchCommand = program
  .command('match')
  .description('The main function takes a list of strings and one or more glob patterns to use for matching.')
  .option('-f, --files <files...>', 'Specify file paths')
  .option('-p, --patterns <patterns...>', 'Specify match patterns')
  .action((options) => {
    const matchedFiles = micromatch(options.files, options.patterns)
    if (matchedFiles.length == 0) {
      process.stdout.write("Not Match");
    }
    else {
      process.stdout.write("Match Files: " + matchedFiles.toString());
    }
  });

const isMatchCommand = program
  .command('ismatch')
  .description('Returns true if the specified string matches the given glob patterns.')
  .option('-f, --file <file>', 'Specify file paths')
  .option('-p, --patterns <patterns...>', 'Specify match patterns')
  .action((options) => {
    const matchedFiles = micromatch(options.file, options.patterns)
    if (matchedFiles.length == 0) {
      process.stdout.write("false");
    }
    else {
      process.stdout.write("true");
    }
  });

  function exportFilePath(writeStream, filePath) {
    writeStream.write(filePath + '\n');
  }

  function checkFilePathExists(filePath) {
    return new Promise((resolve, reject) => {
      fs.access(filePath, fs.constants.F_OK, (err) => {
        if (err) {
          resolve(false);
        } else {
          resolve(true);
        }
      });
    });
  }

  function checkFilePathIsDir(filePath) {
    return new Promise((resolve, reject) => {
      fs.stat(filePath, (err, stats) => {
        if (err) {
          reject(err);
        } else {
          resolve(stats.isDirectory());
        }
      });
    });
  }

  function walk(rootPath, currentPath, patterns, writeStream, exportFilePath, visitedPaths = new Set()) {
    return new Promise((resolve, reject) => {
      fs.readdir(currentPath, function(err, files) {
        if (err) reject(err);
        let promises = [];
        files.forEach(function(file) {
          const filePath = path.join(currentPath, file);
          promises.push(
            new Promise((resolve, reject) => {
              fs.stat(filePath, function(err, stats) {
                if (err) reject(err);
                // (1) check symbolic link recursive 
                fs.lstat(filePath, function(err, linkstats) {
                  if (linkstats.isSymbolicLink()) {
                    const resolvedPath = fs.realpathSync(filePath);
                    if (visitedPaths.has(resolvedPath)) {
                      console.log('Recursive symbolic link detected:', filePath);
                      process.exit(1);
                    }
                    visitedPaths.add(resolvedPath);
                  }
                });
                // (2) if not match ignore pattern, export to file
                const relativeFilePath = path.normalize(path.relative(rootPath, filePath)).replace(/\\/g, '/');
                const files = [];
                files.push(relativeFilePath);
                const matchedFiles = micromatch(files, patterns);
                if (matchedFiles.length == 0) {
                  exportFilePath(writeStream, relativeFilePath);
                }
                // (3) walk dir 
                if (stats.isDirectory()) {
                  walk(rootPath, filePath, patterns, writeStream, exportFilePath, visitedPaths).then(resolve).catch(reject);
                } else {
                  resolve();
                }
              });
            })
          );
        });
        Promise.all(promises).then(resolve).catch(reject);
      });
    });
  }
    
  const getPackagedFilesCommand = program
    .command('getPackagedFiles')
    .description('Get files and directories which will be packaged, and export file path to target file')
    .option('-s, --source <source>', 'Source dir')
    .option('-t, --target <target>', 'Target file path')
    .option('-p, --patterns <patterns...>', 'Ignore Patterns', [])
    .action((options) => {
      // (1) check if parameters are null
      if (!options.source || !options.target) {
        console.error('source or target paramerter should not be empty!');
        process.exit(1);
      }

      const sourcePath = options.source;
      const targetFile = options.target;
      const targetPath = path.dirname(targetFile)
      const patterns = options.patterns.map(pattern => pattern.replace(/\\/g, '/'));
  
      // (2) check if source path is exist
      checkFilePathExists(sourcePath)
      .then(exists => {
        if (!exists) {
          console.error('Source Path ', sourcePath, ' is not exist.')
          process.exit(1);
        }
      })
      .catch(err => {
        console.error('Error checking file path:', err);
        process.exit(1);
      });

      // (3) check if target path is exist
      checkFilePathExists(targetPath)
      .then(exists => {
        if (!exists) {
          console.error('Target Path ', targetPath, ' is not exist.')
          process.exit(1);
        }
      })
      .catch(err => {
        console.error('Error checking file path:', err);
        process.exit(1);
      });

      // (4) check if source path is dir
      checkFilePathIsDir(sourcePath)
      .then(isDir => {
        if (!isDir) {
          console.log('Source Path ', sourcePath, ' is not a directory.');
          process.exit(1);
        }
      })
      .catch(err => {
        console.error('Error checking file path:', err);
        process.exit(1);
      });
      
      // (5) create file stream
      const writeStream = fs.createWriteStream(targetFile);
      writeStream.on('error', function(err) {
        console.error('Error writing to file:', err);
        process.exit(1);
      });
      writeStream.on('finish', function() {
        console.log('Done!');
      });

      // (6) walk file tree from source path, write packaged filepath to file stream
      walk(sourcePath, sourcePath, patterns, writeStream, exportFilePath)
      .then(() => {
        writeStream.end();
      })
      .catch((err) => {
        console.error('Error walking through files:', err);
        process.exit(1);
      });
    });

program.parse(process.argv);

