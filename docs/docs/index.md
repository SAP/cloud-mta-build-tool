# <b>Welcome to the Cloud MTA Build Tool (MBT) </b>

## <b>Introduction</b>


The Cloud MTA Build Tool is a standalone command-line tool that builds a deployment-ready
multitarget application (MTA) archive `.mtar` file from the artifacts of an MTA project according to the project’s MTA
development descriptor (`mta.yaml` file), or from the module build artifacts according to the MTA deployment descriptor (`mtad.yaml` file).

The archive builder is used on a file system that is independent from the development environment where the application project was created.

The archive builder is supported on Windows and Linux operating systems.

If you previously used the [Multitarget Application Archive Builder](https://help.sap.com/viewer/58746c584026430a890170ac4d87d03b/Cloud/en-US/ba7dd5a47b7a4858a652d15f9673c28d.html) for building your MTA projects, see the topic: [differences between the tools](migration.md).

### <b> Multitarget Application</b>

Before using this package, you should be familiar with the multitarget application concept and terminology.
For background and detailed information, see the [Multi-Target Application Model](https://www.sap.com/documents/2016/06/e2f618e4-757c-0010-82c7-eda71af511fa.html) guide.   

### <b>Creating an MTA Archive According to the MTA Development Descriptor (`mta.yaml` file)</b>

The build process and the resulting multitarget application archive depend on the target platform where the archive will be deployed. The currently supported target platforms are SAP Cloud Platform (both the Neo and Cloud Foundry environments) and the SAP HANA XS advanced model.

Each module is built using the build technology that is either associated with the module type by default or is configured explicitly. Then the results of each module's build are packaged into an archive together with the multitarget application deployment descriptor. The deployment descriptor is created based on the structure of the multitarget application project, the multitarget application development descriptor (`mta.yaml` file), the extension descriptors, if they are supplied (the builder merges them with the `mta.yaml` file), and the command-line options. In the deployment descriptor, the design-time modules are converted into the module types that are recognized by the respective target platforms. The resulting multitarget application archive can then be deployed on the required target platform.

For more details about how the design-time module types in the `mta.yaml` file are converted into deployment time module types in the `mtad.yaml` file according to the target platform, see [Module types conversion](https://github.com/SAP/cloud-mta-build-tool/blob/master/configs/platform_cfg.yaml). 

If the MTA Build Tool encounters a module type in the `mta.yaml` file that is not listed in the configuration file, its definition is copied to the `mtad.yaml` file as is and the `zip` builder is applied to the module (in other words, the archived module sources are packaged as the build result).

All resource definitions are passed to the `mtad.yaml` file as is without mapping and validations (for example, if the resource is supported in the target platform or if the specified parameters match the type).

### <b>Creating an MTA Archive According to the MTA Deploymnet Descriptor (`mtad.yaml` file)</b>

The feature is implemented by MBT assembly command. The generate process will copy all modules, module required dependencies and resouces content to MTAR according to MTA Deploymnet Descriptor(mtad.yaml file). 

For module, if path attribute is empty, it will be skipped; if path is not exist, the generate process will be failed. For required dependencies and resources, if parameters.path attribute is empty, it will be skipped; if parameters.path is not exist, the generate process will be failed.

For more details about how to use MBT assembly command, see [How to build an MTA archive from the modules' build artifacts](https://github.com/SAP/cloud-mta-build-tool/blob/master/docs/docs/usage.md#how-to-build-an-mta-archive-from-the-modules-build-artifacts)
