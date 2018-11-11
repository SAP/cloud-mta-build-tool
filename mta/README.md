[![CircleCI](https://circleci.com/gh/SAP/cloud-mta-build-tool.svg?style=svg&circle-token=ecedd1dce3592adcd72ee4c61481972c32dcfad7)](https://circleci.com/gh/SAP/cloud-mta-build-tool)
![GitHub license](https://img.shields.io/badge/license-Apache_2.0-blue.svg)


<b>Disclaimer</b>: The MTA build tool is under heavy development and is currently in a `pre-alpha stage’.
                   Some functionality is still missing and the APIs are subject to change, use at your own risk.
# MTA services

This is a MTA service tool for exploring and validating Multi-target application descriptor (mta.yaml).

The tool commands (API's) are used:

   - Explore the structure of the mta.yaml objects, e.g. retrieve a list of resources required by a specific module.
   - Validate mta.yaml against the specified schema version;
   - Ensure semantic correctness of the mta.yaml, e.g. uniqueness of module/resources names, resolution of requires/provides pairs, etc.
   - Validate the descriptor against the project folder structure, e.g. ‘path’ attribute reference existing project folder;
   

### Multi-Target Applications

A multi-Target application is a package comprised of multiple application and resource modules, 
which have been created using different technologies and deployed to different run-times; however, they have a common life-cycle. 
A user can bundle the modules together, describe (using the `mta.yaml` file) them along with their inter-dependencies to other modules, 
services, and interfaces, and package them in an MTA project.
 

 ## License
 
 MTA Services is [Apache License 2.0 licensed](./LICENSE).
