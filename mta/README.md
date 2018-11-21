[![CircleCI](https://circleci.com/gh/SAP/cloud-mta-build-tool.svg?style=svg&circle-token=ecedd1dce3592adcd72ee4c61481972c32dcfad7)](https://circleci.com/gh/SAP/cloud-mta-build-tool)
![GitHub license](https://img.shields.io/badge/license-Apache_2.0-blue.svg)
![pre-alpha](https://img.shields.io/badge/Release-pre--alpha-orange.svg)


<b>Disclaimer</b>: The MTA explorer services is under heavy development and is currently in a `pre-alpha` stage.
                   Some functionality is still missing and the APIs are subject to change; use at your own risk.
                   
# MTA Explorer Services

MTA tool for exploring and validating the multi-target application descriptor (`mta.yaml`).

The tool commands (APIs) are used to do the following:

   - Explore the structure of the `mta.yaml` file objects, such as retrieving a list of resources required by a specific module.
   - Validate an `mta.yaml` file against a specified schema version.
   - Ensure semantic correctness of an `mta.yaml` file, such as the uniqueness of module/resources names, the resolution of requires/provides pairs, and so on.
   - Validate the descriptor against the project folder structure, such as the `path` attribute reference in an existing project folder.
   - Get data for constructing a deployment MTA descriptor, such as deployment module types.
   

### Multi-Target Applications

A multi-target application is a package comprised of multiple application and resource modules that have been created using different technologies and deployed to different run-times; however, they have a common life cycle. A user can bundle the modules together using the `mta.yaml` file, describe them along with their inter-dependencies to other modules, services, and interfaces, and package them in an MTA project.
 
## Roadmap 

### Milestone 1 
 
 - [x] Supports the MTA parser 
 - [x] Supports development descriptor schema validations (2.1) 
 - [ ] Supports semantic validations (MTA->project)
 - [ ] Supports ‘path’ validation
 
### Milestone 2
 
- [ ] Supports semantic validations (MTA)
- [ ] Supports uniqueness of module and resource names
- [ ] Supports multiple schema support
- [ ] Supports advanced `mta.yaml` file (3.1, 3.2) schemas support
 
### Milestone 3
- [ ] Supports updating scenarios, such as add module/resource, add module property, add dependency, and so on
- [ ] Supports placeholder resolution

## Installation

To install the package, first you need to install Go and set your Go workspace.

Download and install it:

```sh
$ go get -u github.com/mta-explorer/mta
```

## Usage

 - Import it into your source code:

```go
import "github.com/mta-explorer/mta"
```

 -  Quick start example:

```go
// Set path to the mta project
mp := mta.Loc{SourcePath: “wd/my_mta_project”}
// Returns mta object
m, err := mta.ParseFile(&mp)
if err != nil {
	return err
}
// Returns module properties
module, err := m.GetModuleByName(moduleName)
if err != nil {
	return err
}
```
## License
 
Copyright (c) 2018 SAP SE or an SAP affiliate company. All rights reserved.

This file is licensed under the Apache 2.0 License [except as noted otherwise in the LICENSE file](/LICENSE).
