[![CircleCI](https://circleci.com/gh/SAP/cloud-mta-build-tool.svg?style=svg&circle-token=ecedd1dce3592adcd72ee4c61481972c32dcfad7)](https://circleci.com/gh/SAP/cloud-mta-build-tool)
![GitHub license](https://img.shields.io/badge/license-Apache_2.0-blue.svg)
![pre-alpha](https://img.shields.io/badge/Release-pre--alpha-orange.svg)

<b>Disclaimer</b>: The MTA build tool is under heavy development and is currently in a `pre-alpha` stage.
                   Some functionality is still missing and the APIs are subject to change; use at your own risk.
                   

# Prerequisite

You are familiar with the multi-target application concept and terminology. 
For background and detailed information, see The Multi-Target Application Model  guide.                   
                   

#Description

The multi-target application archive builder is a standalone command-line tool that builds a deployment-ready 
multi-target application (MTA) archive .mtar file from the artifacts of an MTA project according to the project’s MTA 
development descriptor (mta.yaml file) or from module build artifacts according to MTA deployment descriptor (mtad.yaml file)
                  
#Usage

## Commands

| Command | usage        | description                                                    | supported 
| ------  | ------       |  ----------                                                    |  ---------- 
| version | `mbt -v`     | Prints the MBT version.                                        | x
| help    | `mbt -h`     | Prints all the available commands.                             | x
| assemble    | `mbt assemble`     | Creates (MTA) archive .mtar file from module build artifacts according to MTA deployment descriptor (mtad.yaml file). Run the command in the directory where the mtad.yaml file is located. Make sure the path property of each modules in mtad.yaml points to the module build artifacts that should be packaged into the target mta archive. | x
| additional commands  | `tbd`              | `tbd`                                 | 
         
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

Soon.
  
## Contributions

Contributions are greatly appreciated.
See [CONTRIBUTING.md](./.github/CONTRIBUTING.md) for details.

## Known Issues

No known major issues.  To report a new [issue](https://github.com/SAP/cloud-mta-build-tool/issues/new/choose), please use our GitHub bug tracking system.

## Support

Please follow our [issue template](./.github/ISSUE_TEMPLATE/bug_report.md) on how to report an issue.
 
 ## License
 
Copyright (c) 2019 SAP SE or an SAP affiliate company. All rights reserved.

This file is licensed under the Apache 2.0 License [except as noted otherwise in the LICENSE file](/LICENSE).
