jQuery.sap.declare("ui5.Component");
sap.ui.getCore().loadLibrary("sap.ui.generic.app");
jQuery.sap.require("sap.ui.generic.app.AppComponent");

sap.ui.generic.app.AppComponent.extend("ui5.Component", {
	metadata: {
		"manifest": "json"
	}
});