/* global QUnit */
QUnit.config.autostart = false;

sap.ui.getCore().attachInit(function () {
	"use strict";

	sap.ui.require([
		"ns/my_nouaa_app/test/unit/AllTests"
	], function () {
		QUnit.start();
	});
});