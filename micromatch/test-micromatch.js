#!/usr/bin/env node
const micromatch = require('micromatch');

// Test 1
console.log("Test 1")
console.log(micromatch(['a/b/3.js'], ['a/b/**', '!a/b/3.js']))
console.log("")

// Test 2
console.log("Test 2")
console.log(micromatch(['a/b/3.js', 'a/b/4.js', 'a/b/5.js', 'a/b/6.js'], ['a/b/**', '!a/b/3.js']))
console.log(micromatch(['a/b/3.js'], ['a/b/**', '!a/b/3.js']))
console.log(micromatch(['a/b/4.js'], ['a/b/**', '!a/b/3.js']))
console.log(micromatch(['a/b/5.js'], ['a/b/**', '!a/b/3.js']))
console.log(micromatch(['a/b/6.js'], ['a/b/**', '!a/b/3.js']))
console.log("")


// Test 3
console.log("Test 3")
console.log(micromatch(['a/b/3.js'], ['a/b/**', '!a/b/3.js', '!a/b/4.js']))
console.log("")

// Test 4
console.log("Test 4")
console.log(micromatch(['a/b/3.js', 'a/b/4.js', 'a/b/5.js', 'a/b/6.js'], ['a/b/**', '!a/b/3.js', '!a/b/4.js']))
console.log(micromatch(['a/b/3.js'], ['a/b/**', '!a/b/3.js', '!a/b/4.js']))
console.log(micromatch(['a/b/4.js'], ['a/b/**', '!a/b/3.js', '!a/b/4.js']))
console.log(micromatch(['a/b/5.js'], ['a/b/**', '!a/b/3.js', '!a/b/4.js']))
console.log(micromatch(['a/b/6.js'], ['a/b/**', '!a/b/3.js', '!a/b/4.js']))
console.log("")

// Test 5
console.log("Test 5")
console.log(micromatch(['a/b/3.js', 'a/b/4.js', 'a/b/5.js', 'a/b/6.js'], ['a/b/!(3.js)', 'a/b/!(4.js)']))
console.log(micromatch(['a/b/3.js'], ['a/b/!(3.js)', 'a/b/!(4.js)']))
console.log(micromatch(['a/b/4.js'], ['a/b/!(3.js)', 'a/b/!(4.js)']))
console.log(micromatch(['a/b/5.js'], ['a/b/!(3.js)', 'a/b/!(4.js)']))
console.log(micromatch(['a/b/6.js'], ['a/b/!(3.js)', 'a/b/!(4.js)']))
console.log("")

// Test 6
console.log("Test 6")
console.log(micromatch(['a/b/3.js', 'a/b/4.js', 'a/b/5.js', 'a/b/6.js'], ['a/b/**', 'a/b/!(3.js)', '!a/b/3.js', '!a/b/4.js']))
console.log(micromatch(['a/b/3.js'], ['a/b/**', 'a/b/!(3.js)', '!a/b/3.js', '!a/b/4.js']))
console.log(micromatch(['a/b/4.js'], ['a/b/**', 'a/b/!(3.js)', '!a/b/3.js', '!a/b/4.js']))
console.log(micromatch(['a/b/5.js'], ['a/b/**', 'a/b/!(3.js)', '!a/b/3.js', '!a/b/4.js']))
console.log(micromatch(['a/b/6.js'], ['a/b/**', 'a/b/!(3.js)', '!a/b/3.js', '!a/b/4.js']))
console.log("")

// Test 7
console.log("Test 7")
console.log(micromatch(['a/b/3.js'], ['!a/b/3.js', 'a/b/**']))
console.log("")

// Test 8
console.log("Test 8")
console.log(micromatch(['a/b/3.js', 'a/b/4.js', 'a/b/5.js', 'a/b/6.js'], ['!a/b/4.js', 'a/b/**', '!a/b/3.js']))
console.log(micromatch(['a/b/3.js'], ['!a/b/4.js', 'a/b/**', '!a/b/3.js']))
console.log(micromatch(['a/b/4.js'], ['!a/b/4.js', 'a/b/**', '!a/b/3.js']))
console.log(micromatch(['a/b/5.js'], ['!a/b/4.js', 'a/b/**', '!a/b/3.js']))
console.log(micromatch(['a/b/6.js'], ['!a/b/4.js', 'a/b/**', '!a/b/3.js']))
console.log("")

// Test 9
console.log("Test 9")
console.log(micromatch(['a/b/3.js', 'a/b/4.js', 'a/b/5.js', 'a/b/6.js'], ['!a/b/3.js', '!a/b/4.js', 'a/b/**']))
console.log(micromatch(['a/b/3.js'], ['!a/b/3.js', '!a/b/4.js', 'a/b/**']))
console.log(micromatch(['a/b/4.js'], ['!a/b/3.js', '!a/b/4.js', 'a/b/**']))
console.log(micromatch(['a/b/5.js'], ['!a/b/3.js', '!a/b/4.js', 'a/b/**']))
console.log(micromatch(['a/b/6.js'], ['!a/b/3.js', '!a/b/4.js', 'a/b/**']))
console.log("")

// Test 10
console.log("Test 10")
console.log(micromatch(['a/b/3.js', 'a/b/4.js', 'a/b/5.js', 'a/b/6.js'], ['a/b/!(4.js)', 'a/b/!(3.js)']))
console.log(micromatch(['a/b/3.js'], ['a/b/!(4.js)', 'a/b/!(3.js)']))
console.log(micromatch(['a/b/4.js'], ['a/b/!(4.js)', 'a/b/!(3.js)']))
console.log(micromatch(['a/b/5.js'], ['a/b/!(4.js)', 'a/b/!(3.js)']))
console.log(micromatch(['a/b/6.js'], ['a/b/!(4.js)', 'a/b/!(3.js)']))
console.log("")


// Test 11
console.log("Test 11")
console.log(micromatch(['a/b/3.js', 'a/b/4.js', 'a/b/5.js', 'a/b/6.js'], ['a/b/!(3.js)', 'a/b/**', '!a/b/3.js', '!a/b/4.js']))
console.log(micromatch(['a/b/3.js'], ['a/b/!(3.js)', 'a/b/**', '!a/b/3.js', '!a/b/4.js']))
console.log(micromatch(['a/b/4.js'], ['a/b/!(3.js)', 'a/b/**', '!a/b/3.js', '!a/b/4.js']))
console.log(micromatch(['a/b/5.js'], ['a/b/!(3.js)', 'a/b/**', '!a/b/3.js', '!a/b/4.js']))
console.log(micromatch(['a/b/6.js'], ['a/b/!(3.js)', 'a/b/**', '!a/b/3.js', '!a/b/4.js']))
console.log("")

// Test 12
console.log("Test 12")
console.log(micromatch(['a/b/3.js', 'a/b/4.js', 'a/b/5.js', 'a/b/6.js'], ['a/b/**', '!a/b/3.js', '!a/b/4.js', 'a/b/!(3.js)']))
console.log(micromatch(['a/b/3.js'], ['a/b/**', '!a/b/3.js', '!a/b/4.js', 'a/b/!(3.js)']))
console.log(micromatch(['a/b/4.js'], ['a/b/**', '!a/b/3.js', '!a/b/4.js', 'a/b/!(3.js)']))
console.log(micromatch(['a/b/5.js'], ['a/b/**', '!a/b/3.js', '!a/b/4.js', 'a/b/!(3.js)']))
console.log(micromatch(['a/b/6.js'], ['a/b/**', '!a/b/3.js', '!a/b/4.js', 'a/b/!(3.js)']))
console.log("")

// Test 13
console.log("Test 13")
console.log(micromatch(['a/b/6.js'], []))
console.log("")

// Test 14
console.log("Test 14")
console.log(micromatch(['a/b/6.js', 'a/b/.ignorefile'], ['**/.*']))
console.log("")

// Test 15
console.log("Test 15")
console.log(micromatch(['a/b/6.js', 'a/b/.ignorefile'], ['.*']))
console.log("")

// Test 16
console.log("Test 16")
console.log(micromatch(['a/b/6.js', 'a/b/.ignorefile'], ['.??*']))
console.log("")

// Test 17
console.log("Test 17")
console.log(micromatch(['a/b/6.js', '.ignorefile'], ['.*']))
console.log("")

// Test 18
console.log("Test 18")
console.log(micromatch(['a/b/6.js', '.ignorefile'], ['.??*']))
console.log("")

// Test 19
console.log("Test 19")
console.log(micromatch(['.invisible-dir'], ['.??*']))
console.log("")

// Test 20
console.log("Test 20")
console.log(micromatch(['.invisible-dir', '.invisible-dir/visible-file'], ['.??*']))
console.log("")

// Test 21
console.log("Test 21")
console.log(micromatch(['.invisible-dir', '.invisible-dir/.invisible-file'], ['.??*']))
console.log("")

// Test 22
console.log("Test 22")
console.log(micromatch(['a/b/c.js'], ['*.js']))
console.log("")

// Test 23
console.log("Test 23")
console.log(micromatch(['a/b/c.js'], ['**/*.js']))
console.log("")