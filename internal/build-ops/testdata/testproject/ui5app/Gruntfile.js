module.exports = function(grunt) {
	"use strict";
	grunt.loadNpmTasks("@company/grunt-companyui5module-bestpractice-build");
	grunt.registerTask("default", [
		"lint"
	]);
};