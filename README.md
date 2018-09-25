[![CircleCI](https://circleci.com/gh/SAP/cloud-mta-build-tool.svg?style=svg&circle-token=ecedd1dce3592adcd72ee4c61481972c32dcfad7)](https://circleci.com/gh/SAP/cloud-mta-build-tool)
![GitHub license](https://img.shields.io/badge/license-Apache_2.0-blue.svg)

<b>Disclaimer</b>:  This repository is <b>under development</b> and currently unstable, don't use until there is an official release.  
# Multi-target Application (MTA) Archive Builder


The MTA command-line archive builder provides a convenient way to bundle an MTA project into an MTAR file (MTA archive).

### Multi-Target Applications

A multi-Target application is a package comprised of multiple application and resource modules, 
which have been created using different technologies and deployed to different runtimes; however, they have a common life-cycle. 
A user can bundle the modules together, describe (using the `mta.yaml` file) them along with their inter-dependencies to other modules, 
services, and interfaces, and package them in an MTA project.
 

## MTA Archive Builder Tool 

The MTA archive builder tool will provide a clear separation between the generation process and the build process as follows:

### CLI 

The CLI tool will:
- Parse and analyze the development descriptor a.k.a mta.yaml file and generate a Makefile accordingly. 
- Provide atomic commands that can be executed as isolated process.
- Build a `META-INF` folder with the following content:
  - Translation of the `mta.yaml` source file into the `mtad.yaml` deployment descriptor.
  - `META-INFO` file that describe the build artifact structure.
  
  
#### [Makefile](https://www.gnu.org/software/make/)

The generated `Makefile` (GNU Make) will describe and execute the build process with two flavors:
- default - provide a generic build process that can be modified according to the project needs.
- verbose - provide verbose build file as manifest which describe each step in separate target (experimental).

During the build process the generated `makefile` is responsible for the following:
- Building each of the modules in the MTA project.
- Invoking the CLI commands in the right order. 
- Providing an MTA archive that is ready for deployment.

## Commands <a id='commands'></a>

The MBT supports the following commands:


| Command | usage        | description                                            |
| ------  | ------       |  ----------                                            |
| version | `mbt -v`     | Prints the MBT version.                                 |
| help    | `mbt -h`     | Prints all the available commands.                     | 
| init    | `mbt init`   | Generates the `makefile` according to the `mta.yaml` file.             |
| TBD     |              | Additional commands will be added as they become available.



## What is an MTA Project

For background and detailed information, see The [Multi-Target Application Model](http://help.sap.com/disclaimer?site=http://www.sap.com/documents/2016/06/e2f618e4-757c-0010-82c7-eda71af511fa.html) Information published on SAP site.

## Todo's

 - [ ] Move MTA & validations package to separate GitHub repository 
 - [ ] Support first MVP scenarios such as:
 
   - [ ] Feature build
   - [ ] XMake integration 
   - [ ] Partial build
   
 - [ ] Release process
 - [ ] Usage
 - [ ] Add concrete limitations per release

 
 ## Limitations
 
   - TBD
 
 
 ## License
 
 MTA Build Tool is [Apache License 2.0 licensed](./LICENSE).