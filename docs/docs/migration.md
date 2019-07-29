### Overview
If you have previously used the [Multitarget Application Archive Builder](https://help.sap.com/viewer/58746c584026430a890170ac4d87d03b/Cloud/en-US/ba7dd5a47b7a4858a652d15f9673c28d.html) for building your MTA projects, you should be aware of the differences between the tools.


#### Features that are handled differently in the Cloud MTA Build Tool

* The Cloud MTA build Tool uses `GNU Make` technology for building an MTA project. Therefore, you should have `GNU Make` installed in your build environmnet. 

For more information, see sections [`GNU Make` installation](makefile.md) and [commands for building a project](usage.md#how-to-build-an-mta-archive-from-the-project-sources). 
&nbsp;
* Packaging of HTML5 modules in `deploy_mode=html5-repo`
You need to update your `mta.yaml` file to exclude `html5` modules from the resulting MTA archive and configure the build result folder. In order to do that, add the following to the `build-parameters` section for each  module of this type:

```yaml

   build-parameters:
      supported-platforms: []
      build-result: dist
```
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; For more information about the `supported-platforms` build parameter, see [Configuring and Packaging Modules According to Target Platforms](configuration.md#configuring-and-packaging-modules-according-to-target-platforms).


* The following `build-parameters` are not supported by the Cloud MTA Build Tool: <ul><li>`npm-opts`<li>`grunt-opt`<li>`maven-opts`</ul>

  If you need to change the default build behavior defined for the corresponding builder, see [configure `custom` builder](configuration.md#configuring-the-custom-builder). For a complete list of available builders and their default behaviors, see [Builders execution commands](https://github.com/SAP/cloud-mta-build-tool/blob/master/configs/builder_type_cfg.yaml).
  
---
**NOTE**
If you try to build the project without these changes, the build tool will identify the cases and will fail the build with the corresponding errors.

---

#### New features in the Cloud MTA Build Tool

* In addition to configuring the build behaviour in `mta.yaml`, you can configure build process of the specific module or the whole project in the `Makefile.mta` file, which you can generate the file using [`mbt init` command](usage.md#cloud-mta-build-tool-commands). The generated file contains the default configurations for buidling the MTA project according to our best practices.
&nbsp;&nbsp;
* If you want to run a builder process before running builders of the specific modules, define it using [global `before-all` build parameters](configuration.md#configuring-global-build).  
&nbsp; 
* You can define your own build commands as described here: [configuring `custom` builder](configuration.md#configuring-the-custom-builder). 
&nbsp; 

#### Features that are temporarily not available in the Cloud MTA Build Tool

The following features are supported by the [Multitarget Application Archive Builder](https://help.sap.com/viewer/58746c584026430a890170ac4d87d03b/Cloud/en-US/ba7dd5a47b7a4858a652d15f9673c28d.html) and will be provided in the Cloud MTA Build Tool soon:


* Running MTA builds with extension files

