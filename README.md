[![CircleCI](https://circleci.com/gh/SAP/cloud-mta-build-tool.svg?style=svg&circle-token=ecedd1dce3592adcd72ee4c61481972c32dcfad7)](https://circleci.com/gh/SAP/cloud-mta-build-tool)
![GitHub license](https://img.shields.io/badge/license-Apache_2.0-blue.svg)
![pre-alpha](https://img.shields.io/badge/Release-pre--alpha-orange.svg)

<b>Disclaimer</b>: The MTA build tool is under heavy development and is currently in a `pre-alpha` stage.
                   Some functionality is still missing and the APIs are subject to change; use at your own risk.
                   
# Multi-Target Application (MTA) Archive Builder

The MTA command-line archive builder provides a convenient way to bundle an MTA project into an MTAR file (Multi Target Application aRchive).

### Multi-Target Applications

A multi-target application is a package comprised of multiple application and resource modules that have been created using different technologies and deployed to different runtimes; however, they have a common life cycle. A user can bundle the modules together along with their interdependencies to other modules, services, and interfaces, and package them in an MTA project, describing them using the `mta.yaml` file.
 

## MTA Archive Builder Tool 

The MTA archive builder tool (MBT) will provide a clear separation between the generation process and the build process as follows:

### CLI 

The CLI tool will:
- Parse and analyze the development descriptor, a.k.a `mta.yaml` file, and generate a `Makefile` accordingly. 
- Provide atomic commands that can be executed as an isolated process.
- Build a `META-INF` folder containing the following content:
  - Translation of the `mta.yaml` source file into the `mtad.yaml` deployment descriptor.
  - A `META-INFO` file that describe the build artifact structure.
  
  
#### [Makefile](https://www.gnu.org/software/make/)

The generated `Makefile` (GNU Make) will describe and execute the build process in two flavors:
- default - Provides a generic build process that can be modified according to the project needs.
- verbose - Provides a verbose build file as a manifest that describes each step in a separate target (experimental).

During the build process the generated `Makefile` is responsible for the following:
- Building each of the modules in the MTA project.
- Invoking the CLI commands in the right order. 
- Providing an MTA archive that is ready for deployment.

## Commands <a id='commands'></a>

The MBT supports the following commands:


| Command | usage        | description                                            |
| ------  | ------       |  ----------                                            |
| version | `mbt -v`     | Prints the MBT version.                                 |
| help    | `mbt -h`     | Prints all the available commands.                     | 
| assemble    | `mbt assemble`     | Assemble MTA project according to deployment descriptor.                     | 



## What is an MTA Project

For background and detailed information, see The [Multi-Target Application Model](http://help.sap.com/disclaimer?site=http://www.sap.com/documents/2016/06/e2f618e4-757c-0010-82c7-eda71af511fa.html) information published on the SAP website.


## Roadmap
 
### Milestone 1  - (Q1 - 2019)

 - [x] Supports project assembly based deployment descriptor 
 - [ ] Supports build of HTML5 applications (non repo)
 - [ ] Supports build of node applications
 - [ ] Partial support of build parameters (first phase)
    - [ ] Supports build dependencies
    - [ ] Supports the copying of build results from other modules
    - [ ] Supports the build results from a different location
    - [ ] Supports target platforms
 - [ ] Generates a default `Makefile`
 - [ ] Generates a `mtad.yaml` file from a `mta.yaml` file
 - [ ] Supports builds for `XSA` / `CF` targets 
 
### Milestone 2 - (Q2 - 2019)
 
  - [ ] Generates a verbose `Makefile`
  - [ ] Supports MTA extension
  - [ ] Supports build of Java/Maven applications
  - [ ] Supports ZIP builds
  - [ ] Supports fetcher build 
  - [ ] Supports build parameters
    - [ ] Supports build options
    - [ ] Supports ignore files/folders
    - [ ] Supports the definition of timeouts
    - [ ] Supports build artifact naming
  - [ ] Supports multi-schema
  - [ ] Supports enhancing schema validations
  - [ ] Supports semantic validations
  - [ ] Partial supports advanced `mta.yaml` (3.1, > 3.2) schema
  
 
 ### Milestone 3 - (Q3 - 2019)
 
  - [ ] Supports parallel execution for default `Makefile` 
  - [ ] Supports incremental build (one module at a time)
 
 ### Milestone 4 - (Q3 - 2019)

 - [ ] Supports extensibility framework
 - [ ] Full supports advanced `mta.yaml` (3.1, > 3.2) schema

## Download and Installation

Download the binary according to your operating system, unzip it and add it to your `~/bin` path.
  
## Contributions

Contributions are greatly appreciated.
See [CONTRIBUTING.md](./.github/CONTRIBUTING.md) for details.

## Known Issues

No known major issues.  To report a new issue, please use our GitHub bug tracking system.

## Support

Please follow our [issue template](./.github/ISSUE_TEMPLATE/bug_report.md) on how to report an issue.
 
 ## License
 
Copyright (c) 2019 SAP SE or an SAP affiliate company. All rights reserved.

This file is licensed under the Apache 2.0 License [except as noted otherwise in the LICENSE file](/LICENSE).
