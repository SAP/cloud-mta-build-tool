
###Overview

Optionally, you can define the builder behavior by configuring parameters in build-parameters section in mta.yaml for each module or globaly.

For each module, the MTA Build Tool performs a build operation using a default technology specific for the module type, for example, npm builder.

For complete list of available builders and their default behaviour, see [Builders execution commands](https://github.com/SAP/cloud-mta-build-tool/blob/master/configs/builder_type_cfg.yaml).


See [Modules' default builder configuration](https://github.com/SAP/cloud-mta-build-tool/blob/master/configs/module_type_cfg.yaml) to know what builder is applied for a module if `builder` parameters is not explicily set in the MTA development descriptor file (`mta.yaml`). 

The following sections describe in detail how to configure each of the supported builders.

####Configuring a builder for a module 
To use a non-default builder for a module, specify its name in the builder parameter.

```yaml

- name: module1
   type: java
   build-parameters:     
     builder: zip
     
```

#### Configuring Build order 
You can define dependencies between modules to ensure that they are built in a certain order. In addition, you can use the build results of one module to build another module.
Also you can can configure a build step that is performed before building any of the modules.

##### Defining module build order
To define the build order of modules, in the mta.yaml file, add a `requires` section under the `build-parameters` section of the dependent module; that is, a module that should be built only after building another required module. Then specify the required module name as a property of the dependent module’s `requires` section.

The example mta.yaml file below demonstrates how to build module A after module B:


```yaml

ID: myproject
_schema-version: '3.1'
version: 0.0.1

modules:
 
 - name: A
   type: html5 
   path: pathtomoduleA
   build-parameters:
      requires:
        - name: B
          
 - name: B
   type: html5 
   path: pathtomoduleB
     
```
> **_NOTE:_** The order of modules in the mta.yaml file is not important as long as the requires sections have been coded with the relevant dependencies.

> **_CAUTION:_** Cyclical dependencies are not allowed. If such dependencies are found, the build fails.

<br>
##### Using build results for building another modulecopy artifact

If you want to use the build results of module B for building module A, modify the `requires` section as follows:

```yaml

ID: myproject
_schema-version: '3.1'
version: 0.0.1

modules:
 
 - name: A
   type: html5 
   path: pathtomoduleA
   build-parameters:
      requires:
        - name: B
          artifacts: ["*"]
          target-path: "newfolder"
          
 - name: B
   type: html5 
   path: pathtomoduleB
     
```
<br>

##### Configuring "before-all" build
If you would like to run some builder process before running builders of the specific modules, define it in the `build-parameters` section at global level in mta.yaml as follows:

```yaml

ID: myproject
_schema-version: '3.1'
version: 0.0.1

build-parameters:
  before-all:
    builders:
      - builder: npm

modules:
 ...
     
```
You can configure any of the supported buiders. 
Also you can list several builders and they will run in the specified order.

```yaml

ID: myproject
_schema-version: '3.1'
version: 0.0.1

build-parameters:
  before-all:
    builders:
      - builder: custom
        commands: 
          - npm run script1
      - builder: grunt
        

modules:
 ...
     
```
> **_NOTE:_** The following build parameters only are considered when configuing builder at "before-all" level: <li>timeout<li>commands (for the `custom` builder only)<li> fetcher-opts (for the `fetcher` builder only)
<br>

####Configuring the Fetcher Builder
The fetcher can be used for packaging already available Maven artifacts into a multitarget application archive as the module’s build results.

```yaml

- name: my_module
  build-parameters:
      builder: fetcher
      fetcher-opts:
         repo-type: maven
         repo-coordinates: mygroup:myart:1.0.0
     
```
repo-coordinates is a string of the form `groupId:artifactId:version[:packaging[:classifier]]`.

This is equivalent to a Maven dependency command, delegating the specified coordinates to the standard Maven dependency mechanism: `mvn dependency:copy -Dartifact=mygroup:myart:1.0.0`

####Configuring the Custom Builder
You can define your own build commands by configuring `custom` builder as follows:

```yaml

- name: my_module
  build-parameters:
      builder: custom
      commands:
        - npm install
        - grunt
        - npm prune --production

     
```
When used for configuring a module builder, you can omit `commands` parameters; then no build is performed and the module's root folder or a folder specified in `build-result` will be packaged into the MTA archive.

####Configuring module build artifacts to package into MTA archive
You can configure the following build parameters to define arctifacts to package into MTA archive for the specific module:

| Name | Default value        | Description                                                    
| ------  | --------       |  ----------                                                
| build-result    | The module folder     | A folder with the build results to package.
| ignore    | None     | Files and/or subfolders to exclude from the package. 



For example:

```yaml

- name: module1
   type: java
   build-parameters:     
     builder: zip
     build-result: myfolder
     ignore: ["*.txt", "mtaignore/"]

     
```

> **_NOTE:_** These parameters are not considered for the `fetcher` builder



#### Configuring and Packaging Modules According to Target Platforms

If you want to control which modules should be packaged into a multitarget application archive for a specific platform during a multitarget application build, use the `supported-platforms` build parameter as follows:

```yaml

ID: myproject
_schema-version: '2.0'
version: 0.0.1

modules:
 
 - name: A
   type: html5 
   path: pathtomoduleA
   build-parameters:
      supported-platforms: [NEO, CF]
          
 - name: B
   type: html5 
   path: pathtomoduleB
   build-parameters:
      supported-platforms: [NEO]

 - name: C
   type: html5 
   path: pathtomodulec
   build-parameters:
      supported-platforms: []

     
```
If you build an application with the multitarget application descriptor above for a `CF` target, all its modules will be built but the resulting multitarget application archive and deployment descriptor will include module A only.

If you build the same application for a `NEO` target, the multitarget application archive will include module A and module B.

If you have a module C with an empty supported platform list, it is never included in any multitarget application archive, regardless of the target platform of the specific multitarget application build.

If the `supported-platforms` build parameter is not used, the module is packaged for any target platform.

You can use these values or any combination of these values for the `supported-platforms` build parameter: <ul><li>`CF` for SAP Cloud Platform Cloud Foundry Environment  <li>`NEO` for SAP Cloud Platform Neo Environment <li>`XSA` for SAP HANA XS advanced model 

#### Configuring Timeout Sessions
This feature is not supported yet by the tool.

#### Configuring the Build Artifact Name
This feature is not supported yet by the tool.


