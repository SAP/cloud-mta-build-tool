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

