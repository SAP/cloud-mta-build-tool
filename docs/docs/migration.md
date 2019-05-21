### Overview
If you have previously used the [Multitarget Application Archive Builder](https://help.sap.com/viewer/58746c584026430a890170ac4d87d03b/Cloud/en-US/ba7dd5a47b7a4858a652d15f9673c28d.html) for building your MTA projects, you should be aware of the differences between the tools.


#### Features that handled differently in the Cloud MTA Build Tool

* Packaging of HTML5 modules in `deploy_mode=html5-repo`
You need to update your `mta.yaml` file to exclude `html5` modules from the resulting MTA archive. In order to do that, add the following to the `build-parameters` section for each  module of this type:

```yaml

   build-parameters:
      supported-platforms: []
```
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; For more information about the `supported-platforms` build parameter, see [Configuring and Packaging Modules According to Target Platforms](configuration.md#configuring-and-packaging-modules-according-to-target-platforms)
  
* The following `build-parameters` are not supported by the Cloud MTA Build Tool: <ul><li>`npm-opts`<li>`grunt-opt`<li>`maven-opts`</ul>

  If you need to change the default build behavior defined for the correponding builder, see [configure `custom` builder](configuration.md#configuring-the-custom-builder).
  For a complete list of available builders and their default behaviors, see [Builders execution commands](https://github.com/SAP/cloud-mta-build-tool/blob/master/configs/builder_type_cfg.yaml).
  <br>

#### New features in the Cloud MTA Build Tool

* If you want to run a builder process before running builders of the specific modules, define it according to [global `before-all` build parameters](configuration.md#configuring_before_all_build).
* You can define your own build commands as described here: [configuring `custom` builder](configuration.md#configuring-the-custom-builder).


#### Features that are temporarily not available in the Cloud MTA Build Tool

The following features are supported by the [Multitarget Application Archive Builder](https://help.sap.com/viewer/58746c584026430a890170ac4d87d03b/Cloud/en-US/ba7dd5a47b7a4858a652d15f9673c28d.html) and will be provided in the Cloud MTA Build Tool soon:

* Generation of MTA builds with extension files
* Configuration of timeout sessions
* Configuration of build artifact names
