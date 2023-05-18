#!/usr/bin/env node
const micromatch = require('micromatch');
const { program } = require('commander');
const path = require('path');
const util = require('util');
const fs = require('fs');

program
  .version('1.0.0')
  .description('Micromatch CLI Wrapper')

const matchCommand = program
  .command('match')
  .description('The main function takes a list of strings and one or more glob patterns to use for matching.')
  .option('-f, --files <files...>', 'Specify file paths')
  .option('-p, --patterns <patterns...>', 'Specify match patterns')
  .action((options) => {
    const result = micromatch(options.files, options.patterns)
    // options.files.forEach(file => process.stdout.write(util.format("File: %s\n", file)));
    // options.patterns.forEach(pattern => process.stdout.write(util.format("Pattern: %s\n", pattern)));
    if (result.length == 0) {
      process.stdout.write("Not Match");
    }
    else {
      process.stdout.write("Match Files: " + result.toString());
    }
  });

const isMatchCommand = program
  .command('ismatch')
  .description('Returns true if the specified string matches the given glob patterns.')
  .option('-f, --file <file>', 'Specify file paths')
  .option('-p, --patterns <patterns...>', 'Specify match patterns')
  .action((options) => {
    // process.stdout.write(util.format("File: %s\n", options.file));
    // options.patterns.forEach(pattern => process.stdout.write(util.format("Pattern: %s\n", pattern)));
    const result = micromatch(options.file, options.patterns)
    if (result.length == 0) {
      process.stdout.write("false");
    }
    else {
      process.stdout.write("true");
    }
  });
  
  function walk(rootPath, parentPath, patterns, outputFilePath) {
    return new Promise((resolve, reject) => {
      fs.readdir(parentPath, function(err, files) {
        if (err) reject(err);
        let promises = [];
        files.forEach(function(file) {
          const filepath = path.join(parentPath, file);
          promises.push(
            new Promise((resolve, reject) => {
              fs.stat(filepath, function(err, stats) {
                if (err) reject(err);
                const relativeFilePath = path.normalize(path.relative(rootPath, filepath)).replace(/\\/g, '/');
                // console.log("root path: " + rootPath)
                // console.log("parent path: " + parentPath)
                // console.log("file path: " + filepath)
                // console.log("relative path: " + relativeFilePath)
                // console.log("patterns:" + patterns)
                const files = [];
                files.push(relativeFilePath);
                const result = micromatch(files, patterns)
                // console.log("result:" + result)
                if (result.length == 0) {
                  // console.log("file path: " + relativeFilePath + " is not match")
                  outputFilePath(relativeFilePath)
                }
                else {
                  // console.log("file path: " + relativeFilePath + " is match")
                }
                // console.log("")
                if (stats.isDirectory()) {
                  walk(rootPath, filepath, patterns, outputFilePath).then(resolve).catch(reject);
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
  
  const getNotIgnoredFilesCommand = program
    .command('getNotIgnoreFiles')
    .description('Return all files and sub directories which are not matche ignore pattern.')
    .option('-s, --source <source>', 'Source path')
    .option('-t, --target <target>', 'Target file path')
    .option('-p, --patterns <patterns...>', 'Ignore Patterns', [])
    .action((options) => {
      const rootPath = options.source;
      const targetFile = options.target;
      const patterns = options.patterns.map(pattern => pattern.replace(/\\/g, '/'));
  
      if (!rootPath || !targetFile) {
        console.error('Usage: node index.js rootPath targetFile');
        process.exit(1);
      }
  
      // console.log("rootPath: " + rootPath);
      // console.log("targetFile: " + targetFile);
      // console.log("patterns: " + patterns)
  
      const writeStream = fs.createWriteStream(targetFile);
  
      writeStream.on('error', function(err) {
        console.error('Error writing to file:', err);
        process.exit(1);
      });
  
      writeStream.on('finish', function() {
        console.log('Done!');
      });
  
      walk(rootPath, rootPath, patterns, function(filepath) {
        // console.log("current file: " + filepath);
        writeStream.write(filepath + '\n');
      })
      .then(() => {
        writeStream.end();
      })
      .catch((err) => {
        console.error('Error walking through files:', err);
        process.exit(1);
      });
    });



program.parse(process.argv);

