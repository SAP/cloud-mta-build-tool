[![CircleCI](https://circleci.com/gh/SAP/cloud-mta-build-tool.svg?style=svg&circle-token=ecedd1dce3592adcd72ee4c61481972c32dcfad7)](https://circleci.com/gh/SAP/cloud-mta-build-tool)
![GitHub license](https://img.shields.io/badge/license-Apache_2.0-blue.svg)

<b>Disclaimer</b>: This repository is under development  
# MTA Build Tool


 The MTA command-line tool provides a convenient way to bundle an MTA project into an MTAR (MTA Archive). 


### Multi-Target Applications

A Multi-Target Application is a package comprised of multiple application and resource modules, 
which have been created using different technologies and deployed to different runtime, however, have a common life-cycle. 
user bundle the modules together, describe (via mta.yaml) them along with their inter-dependencies to other modules, 
services, and interfaces, and package them in an MTA.
 

#### MTA Build Tool 

The MTA build tool will provide a clear separation between the generation process and the build process as follows:

##### CLI 

- The CLI tool will parse and analyze the development descriptor a.k.a mta.yaml file and generate a Makefile accordingly 
- Provide atomic command's that can be executed as isolated process
- Build META-INF folder with the following content:
  - Translating the mta.yaml source file into the mtad.yaml deployment descriptor.
  - Create META-INFO file which describe the build artifacts structure.
  
  
##### [Makefile](https://www.gnu.org/software/make/)

The generated `Makefile` (GNU Make) will describe and execute the build process with two flavors:
- default - provide a generic build process that can be modified according to the project needs
- verbose - provide verbose build file as manifest which describe each step in separate target (experimental)

During the build process the generated Makefile is responsible on the following:
- building each of the modules in the MTA project.
- invoking the CLI commands in the right order 
- provide mta archive ready for deployment

### Commands <a id='commands'></a>

Following is the command which the MBT support:


| Command | usage        | description                                            |
| ------  | ------       |  ----------                                            |
| version | `mbt -v`     | Prints the MBT version                                 |
| help    | `mbt -h`     | Prints all the available commands                      | 
| init    | `mbt init`   | Generate Makefile according to the mta.yaml            |
| TBD     |              | Additional commands when available to use 



## What is an MTA Project, in a nutshell

An MTA(Multi-target-application) project is defined by a project file in the root folder called `mta.yaml` that contains three different sections:

   * [General information section](#general) - including the MTA ID and version
   * [Modules section](#modules) - describing the content that is actually delivered within the MTA
   * [Resources section](#resources) - describing the resources that the modules expect to have prepared during deployment

Resources can vary between platform - an MTA module intended to run on CloudFoundry might expect a resource
that maps to a _Service Instance_.

### The General Information Section <a id='general'></a>

This section contains information that is relevant for the entire multi-targeted-application. It is used by the MTA deployer
in order to identify the MTA being deployed.

|key|mandatory|constraints|description|
| --- | --- | --- | --- |
|`_schema_version`|[x]|must be a supported by the target deployer|Specifies the version of the MTA spec that is being targeted|
|`ID`|[x]|Must conform to this regexp: `/\A[A-Za-z0-9_\-\.]+\z/`|Identifies the MTA application to the MTA deployer|
|`version`|[x]|Must be a valid _semver_ string|Identifies the version of the MTA application|
|`provider`|[ ]|Free text|The provider of vendor of this software|
|`copyright`|[ ]|Free text|A copyright statement from the provider|
|`parameters`|[ ]|Map element|Provides global deployment parameters to the MTA deployer - not supported on all MTA platforms|


### The Modules Section <a id='modules'></a>

This section contains an entry for each module that is part of the multi-target-application. A module defined in the `mta.yaml` file
is an atomic unit-of-execution that is eventually deployed to a runtime platform such as CloudFoundry. In the source-time `mta.yaml`
the module definition contains the information that ties each module to the entire multi-target-application, as well as information on
how to build the module and what to package into the application.

|*key*|*mandatory*|*constraints*|*description*|
| --- | --- | --- | --- |
|name| [x] |Must conform to this regexp: `/\A[A-Za-z0-9_\-\.]+\z/` and be unique in the `mta.yaml`|An MTA internally unique name.|
|type| [x] |[Supported Values](#module-types)|The expected runtime type of the module.|
|path| [x] |Path expression|The path to the module being built. Must be contained inside the MTA project folder. Use Unix-style path separators ('/').|
|description|[ ]|Free text|A description of the module.|
|provides|[ ]|[Map Element](#provides)|A specification of sets of name-value pairs of configuration data that is provided by this module.|
|requires|[ ]|[Map Element](#requires)|A specification of required configuration that is required by this the module.|
|properties|[ ]|Map element|Properties that will be provided to the deployed application at runtime. Each platform may provide them using a different method. For example, cloud foundry uses environment variables. Valid content depends on the module itself. |
|parameters|[ ]|Map element|Deployment parameters that will be used during the deployment of the module in the target platform. Valid content depends on the targeted platform.|
|build-parameters|[ ]|[Map element](#builders)|Build time information used by the various build tools.|


#### _Provides_ Section <a id='provides'></a>

The _provides_ section allows a module to define sets of name-value pairs as configurations that are made available to other modules in the same MTAR or different MTAR's.

#### _Requires_ Section <a id='requires'></a>

The _requires_ section allows a module to define which configuration sets it needs to receive at deployment time. Dependencies can be provided to an application either via a resource in the _resources_ section or via a _provides_ section from a different module.

### The Resources Section <a id='resources'></a>

This section contains an entry for each resource that must be setup by the MTA deployer for consumption by one of the modules contained in the MTA archive.


|*key*|*mandatory*|*constraints*|*description*|
| --- | --- | --- | --- |
|name|[x]|Must conform to this regexp: `/\A[A-Za-z0-9\_\-\.]+\z/ ` and be unique in the `mta.yaml`|An MTA internaly unique name.|
|type|[ ]|
|parameters|[ ]| |
|properties|[ ]| |

#### Todo's

 - [ ] Support first MVP scenarios such as:
 
   - [ ] Feature build
   - [ ] XMake integration 
   - [ ] Partial build
   
 - [ ] Release process
 - [ ] Usage
 - [ ] Add concrete limitations per release

 
 #### Limitations
 
   - TBD
 
 
 ### License
 
 MTA Build Tool is [Apache License 2.0 licensed](./LICENSE).