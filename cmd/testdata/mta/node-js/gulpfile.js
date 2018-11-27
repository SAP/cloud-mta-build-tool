/**
 * This is the gulp script to run jasmine tests with coverage analysis using istanbul.
 * For more information, please see also:
 * @see {@link http://jasmine.github.io/2.0/introduction.html}
 * @see {@link https://www.npmjs.com/package/istanbul}
 * @see {@link https://www.npmjs.com/package/gulp}
 * @see {@link https://www.npmjs.com/package/gulp-istanbul}
 * @see {@link https://www.npmjs.com/package/gulp-jasmine}
 * @see {@link https://www.npmjs.com/package/jasmine-reporters}
 * @see {@link https://www.npmjs.com/package/glob}
 **/

var path = require("path");

var gulp = require("gulp");
var jasmine = require("gulp-jasmine");
var istanbul = require("gulp-istanbul");
var reporters = require("jasmine-reporters");

// required to replace absolute paths in results with relative ones
var replace = require("gulp-replace");

// configuration of the test coverage, uses glob syntax
// @see {@link https://www.npmjs.com/package/glob}

// the following files are included in the coverage analysis
// include all javascript files and exclude myExclude.js
// replace or remove files depending on what is to be excluded or included in addition
var includedScripts = ["**/*.js", "!myExclude.js"];

// the following files are part of the framework and are to be excluded from the test coverage
var defaultExclusion = ["!**/*spec.js", "!rungulp.js", "!gulpfile.js", "!**/node_modules/**", "!appcontroller.*/**", "!vendor/**"];

// test results folder for test view history
// assign each test run a unique timestamp, coverage and test results
// are associated via timestamp, so the test results folder fills up with files
// of form:
// 123456_report.xml
// 123456_coverage.json
// 456789_report.xml
// 456789_coverage.json
// ...
var testResultsDir = path.join(__dirname, ".testresults");

var timestamp = Date.now();
var testResultFile = timestamp + "_report.xml";
var coverageResultFile = timestamp + "_coverage.json";

/**
 * Instrument the test/productive code
 */
gulp.task("instrumentation", function() {
    return gulp.src(includedScripts.concat(defaultExclusion))
    // Covering files
        .pipe(istanbul({
            includeUntested: true // instrument all files
        }))
        // Force `require` to return covered files
        .pipe(istanbul.hookRequire());
});

/**
 * Execute tests with coverage information, requires instrumentation
 * before tests are executed
 */
gulp.task("jasmine-istanbul", ["instrumentation"], function() {
    // run all tests ending with "spec", skip all tests that are part of the node_modules
    return gulp.src(["**/*spec.js", "!**/node_modules/**"])
        .pipe(jasmine({
            errorOnFail: false,
            // use the standard junit xml reporter
            reporter: new reporters.JUnitXmlReporter({
                savePath: testResultsDir,
                filePrefix: testResultFile,
                consolidateAll: true
            })
        }))
        .pipe(istanbul.writeReports({
            // generage json report for the coverage
            reporters: ["json"],
            reportOpts: {
                json: {
                    dir: testResultsDir,
                    file: coverageResultFile
                }
            }
        }));
});

/**
 * Execute tests without coverage information
 */
gulp.task("jasmine", function() {
    // run all tests ending with "spec", skip all tests that are part of the node_modules
    return gulp.src(["**/*spec.js", "!**/node_modules/**"])
        .pipe(jasmine({
            errorOnFail: false,
            // use the standard junit xml reporter
            reporter: new reporters.JUnitXmlReporter({
                savePath: testResultsDir,
                filePrefix: testResultFile,
                consolidateAll: true
            })
        }));
});

/**
 * Remove absolute path portions so test view can open the referenced files
 */
function postProcessing() {
    return gulp.src([
        path.join(testResultsDir, testResultFile),
        path.join(testResultsDir, coverageResultFile)
    ], {
        dot: true
    })
        .pipe(replace(__dirname, ""))
        .pipe(gulp.dest(testResultsDir));
}

// run tests with coverage analysis
gulp.task("test-coverage",  ["jasmine-istanbul"], function() {
    return postProcessing();
});

// only run the tests, skip test coverage analysis
gulp.task("test", ["jasmine"], function() {
    return postProcessing();
});

//the default task is to run tests with coverage analysis
gulp.task("default", ["test-coverage"]);
