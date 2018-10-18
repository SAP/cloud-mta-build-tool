module.exports = function(grunt) {
	"use strict";
	grunt.loadNpmTasks("@sap/grunt-sapui5module-bestpractice-build");
	grunt.registerTask("default", [
		"lint"
	]);
};