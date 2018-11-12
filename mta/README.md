[![CircleCI](https://circleci.com/gh/SAP/cloud-mta-build-tool.svg?style=svg&circle-token=ecedd1dce3592adcd72ee4c61481972c32dcfad7)](https://circleci.com/gh/SAP/cloud-mta-build-tool)
![GitHub license](https://img.shields.io/badge/license-Apache_2.0-blue.svg)


<b>Disclaimer</b>: The MTA explorer services is under heavy development and is currently in a `pre-alpha stage’.
                   Some functionality is still missing and the APIs are subject to change, use at your own risk.
                   
# MTA Explorer Services

MTA service tool for exploring, validating Multi-target application descriptor (mta.yaml).

The tool commands (API's) are used:

   - Explore the structure of the mta.yaml objects, e.g. retrieve a list of resources required by a specific module.
   - Validate mta.yaml against the specified schema version;
   - Ensure semantic correctness of the mta.yaml, e.g. uniqueness of module/resources names, resolution of requires/provides pairs, etc.
   - Validate the descriptor against the project folder structure, e.g. ‘path’ attribute reference existing project folder;
   - Get data for constructing deployment MTA descriptor, e.g. deployment module types
   

### Multi-Target Applications

A multi-Target application is a package comprised of multiple application and resource modules, 
which have been created using different technologies and deployed to different run-times; however, they have a common life-cycle. 
A user can bundle the modules together, describe (using the `mta.yaml` file) them along with their inter-dependencies to other modules, 
services, and interfaces, and package them in an MTA project.
 
## Roadmap 

### Milestone 1 
 
 - [ ] mta parser 
 - [ ] 2.1 Schema validations 
 - [ ] Semantic validations (mta->project)
‘path’ validation
 
### Milestone 2
 
- [ ] Semantic validations (mta)
- [ ] Uniqueness of module and resource names
- [ ] Multiple schema support
- [ ] Advanced mta.yaml (3.1, 3.2) schemas support
 
### Milestone 3
- [ ] update scenarios, e.g. add module/resource, add module property, add dependency, etc
- [ ] placeholder resolution

## Installation

To install the package, you need to install Go and set your Go workspace first.

1. Download and install it:

```sh
$ go get -u github.com/mta-explorer/mta
```

2. Import it in your code:

```go
import "github.com/mta-explorer/mta"
```

 ## License
 
 MTA Services is [Apache License 2.0 licensed](./LICENSE).
