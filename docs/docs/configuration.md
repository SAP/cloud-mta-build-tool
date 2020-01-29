
###Overview

Optionally, you can define the builder behavior by configuring the parameters in the `build-parameters` section in the `mta.yaml` file for each module or globally.

For each module, the Cloud MTA Build Tool performs a build operation using a default technology specific for the module type, for example, the `npm` builder.

For the complete list of available builders and their default behaviors, see [Builders Execution Commands](https://github.com/SAP/cloud-mta-build-tool/blob/master/configs/builder_type_cfg.yaml).


To find out which builder is applied to a module if the `builder` parameters are not explicitly set in the MTA development descriptor file (`mta.yaml`), see [Modules' default builder configuration](https://github.com/SAP/cloud-mta-build-tool/blob/master/configs/module_type_cfg.yaml). 

The following sections describe in detail how to configure each of the supported builders.

#### Configuring a builder for a module 
To use a non-default builder for a module, specify its name in the builder parameter.

```yaml

- name: module1
   type: java
   build-parameters:     
     builder: zip
     
```

#### Configuring build order 
You can define the dependencies between modules to ensure that they are built in a certain order. In addition, you can use the build results of one module to build another module.
Also you can can configure build steps that are performed before building any of the modules.

##### Defining module build order
To define the build order of modules, in the `mta.yaml` file, add a `requires` section under the `build-parameters` section of the dependent module; that is, a module that should be built only after building another required module. Then specify the required module name as a property of the dependent module’s `requires` section.

The example `mta.yaml` file below demonstrates how to build module A after module B:


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

> **_NOTE:_** The order of modules in the `mta.yaml` file is not important as long as the `requires` sections have been coded with the relevant dependencies.

> **_CAUTION:_** Cyclical dependencies are not allowed. If such dependencies are found, the build fails.

<br>
##### Using build results for building another module

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

#### Configuring a global build

If you want to run addtional build steps before running builders of the specific modules, define it by using the `build-parameters` section at global level in the `mta.yaml` file as follows:


```yaml

ID: myproject
_schema-version: '3.1'
version: 0.0.1

build-parameters:
  before-all:
    - builder: npm

modules:
 ...
     
```

You can configure any of the supported builders. 
Also, you can list several builders and they will run in the specified order.

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
> **_NOTE:_** Only the following build parameters are considered when configuring a builder at the global level: <li>`timeout`<li>`commands` (for the `custom` builder only)<li> `fetcher-opts` (for the `fetcher` builder only)
<br>

#### Configuring the `fetcher` builder
The `fetcher` builder can be used for packaging already available Maven artifacts into a multitarget application archive as the module’s build results.

```yaml

- name: my_module
  build-parameters:
      builder: fetcher
      fetcher-opts:
         repo-type: maven
         repo-coordinates: mygroup:myart:1.0.0
     
```

The `repo-coordinates` parameter is a string of the form `groupId:artifactId:version[:packaging[:classifier]]`.

This is equivalent to a Maven dependency command, delegating the specified coordinates to the standard Maven dependency mechanism: 
`mvn dependency:copy -Dartifact=mygroup:myart:1.0.0`

#### Configuring the `custom` builder
You can define your own build commands by configuring a `custom` builder as follows:

```yaml

- name: my_module
  build-parameters:
      builder: custom
      commands:
        - npm install
        - grunt
        - npm prune --production

     
```

When used for configuring a module builder, you can leave the `commands` parameter list empty:

```yaml

- name: my_module
  build-parameters:
      builder: custom
      commands: []
```
In this case, no build is performed and the module's root folder or a folder specified in the `build-result` is packaged into the MTA archive.

#### Configuring module build artifacts to package into MTA archive
You can configure the following build parameters to define artifacts to package into the MTA archive for the specific module:

| Name | Default value        | Description                                                    
| ------  | --------       |  ----------                                                
| `build-result`    | For the `maven` builder: `<module's folder>/target/*.war` <br><br>  For the `fetcher` builder: `<module's folder>/target/*.*` <br><br> For other builders: `<module's folder>`     | A path to the build results that should be packaged.
| `ignore`    | None     | Files and/or subfolders to exclude from the package. 



For example:

```yaml

- name: module1
   type: java
   build-parameters:     
     builder: zip
     build-result: myfolder
     ignore: ["*.txt", "mtaignore/"]

     
```

> **_NOTE:_** These parameters are not considered for the `fetcher` builder.



#### Configuring and packaging modules according to target platforms

If you want to control which modules should be packaged into a multitarget application archive for a specific platform, use the `supported-platforms` build parameter as follows:

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
   path: pathtomoduleC
   build-parameters:
      supported-platforms: []

     
```

If you build an application with the multitarget application descriptor above for a `CF` (Cloud Foundry) target, all its modules are built, but the resulting multitarget application archive and deployment descriptor include module A only.

If you build the same application for a `NEO` target, the multitarget application archive includes module A and module B.

If you have a module C with an empty supported platform list, it is never included in any multitarget application archive, regardless of the target platform of the specific multitarget application build.

If the `supported-platforms` build parameter is not used, the module is packaged for any target platform.

You can use these values or any combination of these values for the `supported-platforms` build parameter: <ul><li>`CF` for the SAP Cloud Platform Cloud Foundry environment  <li>`NEO` for the SAP Cloud Platform Neo environment <li>`XSA` for the SAP HANA XS advanced model 

#### Configuring timeout sessions
When you build a specific module, there is a default 5-minute timeout allowance. After this time, the build will fail. You can configure the time allowed for timeout when performing a build by adding the `timeout` property to the module build parameters. The timeout property uses the `<number of hours>h<number of minutes>m<number of seconds>s` format.
<br>

For example:

```yaml

- name: module1
   type: java
   build-parameters: 
     timeout: 6m30s
```
Also, you can use this parameter to define timeout for the [global `before-all` build](configuration.md#configuring-global-build).

#### Configuring the build artifact name
The module build results are by default packaged into the resulting archive under the name “data”. You can change this name as needed using the `build-artifact-name` build parameter:  &nbsp;
&nbsp;



```yaml

modules:
  - name: db
    type: hdb
    path: db
    build-parameters:
      build-artifact-name: myfileName
     
```

&nbsp;

> **_NOTE:_** The file extension is not configurable; it is predefined by the module type (.zip or .war).
