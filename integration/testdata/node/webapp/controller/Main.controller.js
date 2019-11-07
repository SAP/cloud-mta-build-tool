sap.ui.define([
    "sap/ui/core/mvc/Controller"
], function(Controller) {
    "use strict";

    return Controller.extend("com.sap.teched.teched.controller.Main", {
        onInit: function() {
            var oView = this.getView();
            var oModel800 = this.getOwnerComponent().getModel("800");
            var oTable = oView.byId("table0");
            oTable.setModel(oModel800);
        },

        onSelectionChange: function(oEvent) {
            var oKey = oEvent.getSource().getSelectedKey();
            switch (oKey) {
                case "400":
                    var oView = this.getView();
                    var oModel400 = this.getOwnerComponent().getModel("400");
                    var oTable = oView.byId("table0");
                    oTable.setModel(oModel400);
                    break;
                case "800":
                    var oView = this.getView();
                    var oModel800 = this.getOwnerComponent().getModel("800");
                    var oTable = oView.byId("table0");
                    oTable.setModel(oModel800);
                    break;
            }

        }

    });
});