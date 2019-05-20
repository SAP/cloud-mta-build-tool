# <b>Welcome to the Cloud MTA Build Tool (a.k.a MBT) </b>

## <b>Introduction</b>


The Cloud MTA Build Tool is a standalone command-line tool that builds a deployment-ready
multi-target application (MTA) archive `.mtar` file from the artifacts of an MTA project according to the project’s MTA
development descriptor (`mta.yaml` file) or from module build artifacts according to the MTA deployment descriptor (`mtad.yaml` file).

The archive builder is used on a file system independently of the development environment in which the application project has been created.

The archive builder is supported on Windows and Linux.

If you previously used [Multitarget Application Archive Builder](https://help.sap.com/viewer/58746c584026430a890170ac4d87d03b/Cloud/en-US/ba7dd5a47b7a4858a652d15f9673c28d.html) for building your MTA projects please learn about [differences between the tools](migration.md).

### <b> Multitarget Application</b>

Before using this package, be sure you are familiar with the multi-target application concept and terminology.
For background and detailed information, see the [Multi-Target Application Model](https://www.sap.com/documents/2016/06/e2f618e4-757c-0010-82c7-eda71af511fa.html) guide.   

### <b>Creating MTA archive according to MTA development descriptor (`mta.yaml` file)</b>

The build process and the resulting multitarget application archive depend on the target platform on which the archive will be deployed. The currently supported target platforms are SAP Cloud Platform (both the Neo and Cloud Foundry environments), and SAP HANA XS advanced model.

Each module is built using the build technology that is either associated with the module type by default or is configured explicitly. Then the results of each module's build are packaged into an archive together with the multitarget application deployment descriptor. The deployment descriptor is created based on the structure of the multitarget application project, the multitarget application development descriptor (mta.yaml), the extension descriptors, if they are supplied (the builder merges them with mta.yaml), and the command line options. In the deployment descriptor, the design-time modules are converted into the module types that are recognized by the respective target platforms. The resulting multitarget application archive can then be deployed on the required target platform.
See [Module types conversion](https://github.com/SAP/cloud-mta-build-tool/blob/master/configs/platform_cfg.yaml) for more details on how design time module types in `mta.yaml` are converted into deployment time module types in `mtad.yaml` according to the target platform. 

If the MTA Build Tool encounters a module type in the `mta.yaml` that is not listed in the configuration file, its definition is copied to `mtad.yaml` as is and the `zip` builder is applied to the module (i.e. archived module sources will be packaged as the build result).

All resource definitions are passed to the `mtad.yaml` file as is without mapping and validations (e.g., if the resource is supported in the target platform or specified parameters match the type).

###<b>Creating MTA archive according to the MTA deploymnet descriptor (`mtad.yaml` file)</b>

TBD
