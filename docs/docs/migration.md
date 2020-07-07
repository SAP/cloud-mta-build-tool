### Overview
If you have previously used the [Multitarget Application Archive Builder](https://help.sap.com/viewer/58746c584026430a890170ac4d87d03b/Cloud/en-US/ba7dd5a47b7a4858a652d15f9673c28d.html) for building your MTA projects, you should be aware of the differences between the tools.


#### Features that are handled differently in the Cloud MTA Build Tool
<ul>
<li>

The Cloud MTA build Tool uses `GNU Make` technology for building an MTA project. Therefore, you should have `GNU Make` installed in your build environmnet. 

For more information, see sections [`GNU Make` installation](makefile.md) and [commands for building a project](usage.md#how-to-build-an-mta-archive-from-the-project-sources). 
&nbsp;
</li>
<li>

Packaging of HTML5 modules in `deploy_mode=html5-repo`
You need to update your `mta.yaml` file to exclude `html5` modules from the resulting MTA archive and configure the build result folder. In order to do that, add the following to the `build-parameters` section for each  module of this type:

```yaml

   build-parameters:
      supported-platforms: []
      build-result: dist
```
For more information about the `supported-platforms` build parameter, see [Configuring and Packaging Modules According to Target Platforms](configuration.md#configuring-and-packaging-modules-according-to-target-platforms). 
&nbsp;
</li>
<li>

The following `build-parameters` are not supported by the Cloud MTA Build Tool: <ul><li>`npm-opts`<li>`grunt-opt`<li>`maven-opts`</ul>

If you need to change the default build behavior defined for the corresponding builder, see [configure `custom` builder](configuration.md#configuring-the-custom-builder). For a complete list of available builders and their default behaviors, see [Builders execution commands](https://github.com/SAP/cloud-mta-build-tool/blob/master/configs/builder_type_cfg.yaml). 
&nbsp;
If you replace the `maven` builder with the `custom` builder, you also need to set the `build-result` parameter (if it was not configured explicitly) as follows:
&nbsp;
```yaml

- name: module1
  type: java
  build-parameters:     
     builder: custom
     ...
     build-result: target/*.war
     ...     
```
&nbsp;
 
</li>
<li>

The Cloud MTA Build tool strictly validates the rule that names of modules, resources, and provided property sets, are unique within the `mta.yaml` file. This ensures that when the name is referenced in the `requires` section, it is unambiguously resolved.  

The [Multitarget Application Archive Builder](https://help.sap.com/viewer/58746c584026430a890170ac4d87d03b/Cloud/en-US/ba7dd5a47b7a4858a652d15f9673c28d.html) allowed the use of the same name for a module and for one of its property sets. For example:

```yaml

- name: SOME_NAME
    type: java
    path: srv
    provides:
      - name: SOME_NAME
        properties:
          url: ${default-url}
```

When migrating to the new build tool, you need to rename either the module or the provided property set. For example:

```yaml

- name: SOME_NAME
    type: java
    path: srv
    provides:
      - name: SOME_NAME_API # New name
        properties:
          url: ${default-url}
```
After renaming, make sure that the places where the name is used refer to the correct artifact. 
&nbsp;
</li>
<li>

The `hdb` builder is not supported by the Cloud MTA Build tool.  You no longer require builder settings for the `hdb` module because the required `npm install --production` command is run by default for this module type.

If you used this builder for other module types, you can repace it with the `npm` builder or use the `custom` builder that runs the `"npm install --production"`command. 
&nbsp;
</li>
&nbsp;

---
**NOTE:**
<b>If you try to build the project without the changes above, the build tool will identify this and will fail the build with the corresponding errors.</b>

---
&nbsp;

<li>

The `hdb` module requires the `package.json` file in the module's folder. If your `hdb` module does not contain the file, you should manually create it with the following content:


```json
{
    "name": "deploy",
    "dependencies": {
        "@sap/hdi-deploy": "^3"
    },
    "scripts": {
        "start": "node node_modules/@sap/hdi-deploy/deploy.js"
    }
}
```

&nbsp;
</li>


---
**NOTE:**
<b>If you try to deploy an MTA archive generated without this change, the operation will fail.</b>

---

&nbsp;
<li>

Service creation or service binding parameters provided in external files and referenced by the `path` property of the correponding entity in the `mta.yaml` are packaged differently by the tools into the result MTA archive.

&nbsp;

Therefore, if you have a `JSON` file with [service creation parameters](https://help.sap.com/viewer/65de2977205c403bbc107264b8eccf4b/Cloud/en-US/a36df26b36484129b482ae20c3eb8004.html) or [service binding parameters](https://help.sap.com/viewer/65de2977205c403bbc107264b8eccf4b/Cloud/en-US/c7b09b79d3bb4d348a720ba27fe9a2d5.html)  and your `JSON` file contains [parameters or placeholders](https://help.sap.com/viewer/65de2977205c403bbc107264b8eccf4b/Cloud/en-US/490c8f71e2b74bc0a59302cada66117c.html) that should be resolved when you deploy the MTA archive, the correponding properties should be moved to the `mta.yaml` file. Otherwise, values assigned to these properties during deployment will be incorrect, because the [parameters or placeholders](https://help.sap.com/viewer/65de2977205c403bbc107264b8eccf4b/Cloud/en-US/490c8f71e2b74bc0a59302cada66117c.html) are resolved only if they are specified within an MTA descriptor, i.e. the `mta.yaml` or `mtad.yaml` files.  &nbsp;

For example, if you provide parameters for creating a UAA service in an `xs-security.json` file:

```yaml

resources:
  - name: my-uaa
    type: com.sap.xs.uaa
    parameters:
      path: ./xs-security.json
```

and your `xs-security.json` file contains a property whose value should be resolved during the MTA archive deployment:

```json

{
  "xsappname": "${default-xsappname}"
}
```

then, you need to modify your `mta.yaml` file as follows:


```yaml

resources:
  - name: my-uaa
    type: com.sap.xs.uaa
    parameters:
      path: ./xs-security.json
      config:
        xsappname: "${default-xsappname}"
```

 In the `xs-security.json` file, you can assign the property a temporary value that does not use the parameter. During deployment, the value specified directly in the MTA descriptor overrides the value specified in the `JSON` file. 

```json

{
  "xsappname": "tmp_appname"
}
```

Alternatively, you can remove the property from your `xs-security.json` file.

&nbsp; &nbsp;


---
**NOTE:**
<b>If you try to deploy the project without the changes above, the deploy service will identify this and will fail the deployment with the corresponding errors.</b>

---


</li>

</ul>

#### Migration of SAP Web IDE Full-Stack Extensions

If you have [custom extensions](https://help.sap.com/viewer/825270ffffe74d9f988a0f0066ad59f0/CF/en-US/cfa254005a63404d98b0820e302729cc.html) for SAP Web IDE Full-Stack and you need to re-build them using the Cloud MTA Build Tool, for example, as part of a new version release process, you should adjust the `mta.yaml` file of the extension MTA project following the steps specified in the comments below:


```yaml

_schema-version: "2.0.0"
ID: sample
version: 0.0.1

modules:
  - name: sample
    type: html5
    path: public
    provides:
      - name: sample   # 1. Change the name of the provided property set so that it is different from the module name.
        public: true
    build-parameters:
      builder: npm
      ignore: [".che/", ".npmrc"]
      timeout: 15m
      requires:
        - name: sample-client
          artifacts: ["dist/*"]
          target-path: "client"
  - name: sample-client
    type: html5
    path: client
    build-parameters:
      builder: npm  # 2. Change the builder type to 'custom'.
      # 3. Add the following 2 lines:
      # commands:
      #   - npm install
      timeout: 15m
      supported-platforms: []
      npm-opts:  # 4. Remove this line and the one below it.
        execute: "install"

```

The final `mta.yaml` file should look like in the example below:

```yaml

_schema-version: "2.0.0"
ID: sample
version: 0.0.1

modules:
  - name: sample
    type: html5
    path: public
    provides:
      - name: sample-provides   
        public: true
    build-parameters:
      builder: npm
      ignore: [".che/", ".npmrc"]
      timeout: 15m
      requires:
        - name: sample-client
          artifacts: ["dist/*"]
          target-path: "client"
  - name: sample-client
    type: html5
    path: client
    build-parameters:
      builder: custom  
      commands:
         - npm install
      timeout: 15m
      supported-platforms: []

```