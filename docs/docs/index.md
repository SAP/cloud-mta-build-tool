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
For background and detailed information, see the [Multi-Target Application Model](https://www.sap.com/documents/2021/09/66d96898-fa7d-0010-bca6-c68f7e60039b.html) guide.   

### <b>Creating an MTA Archive According to the MTA Development Descriptor (`mta.yaml` file)</b>

The build process and the resulting multitarget application archive depend on the target platform where the archive will be deployed. The currently supported target platforms are SAP Cloud Platform (both the Neo and Cloud Foundry environments) and the SAP HANA XS advanced model.

Each module is built using the build technology that is either associated with the module type by default or is configured explicitly. Then the results of each module's build are packaged into an archive together with the multitarget application deployment descriptor. The deployment descriptor is created based on the structure of the multitarget application project, the multitarget application development descriptor (`mta.yaml` file), the extension descriptors, if they are supplied (the builder merges them with the `mta.yaml` file), and the command-line options. In the deployment descriptor, the design-time modules are converted into the module types that are recognized by the respective target platforms. The resulting multitarget application archive can then be deployed on the required target platform.

For more details about how the design-time module types in the `mta.yaml` file are converted into deployment time module types in the `mtad.yaml` file according to the target platform, see [Module types conversion](https://github.com/SAP/cloud-mta-build-tool/blob/master/configs/platform_cfg.yaml). 

If the MTA Build Tool encounters a module type in the `mta.yaml` file that is not listed in the configuration file, its definition is copied to the `mtad.yaml` file as is and the `zip` builder is applied to the module (in other words, the archived module sources are packaged as the build result).

All resource definitions are passed to the `mtad.yaml` file as is without mapping and validations (for example, if the resource is supported in the target platform or if the specified parameters match the type).

### <b>Creating an MTA Archive According to the MTA Deployment Descriptor (`mtad.yaml` file)</b>

The function is provided to assemble multiple modules, which are developed and prebuilt as separate projects, into one MTAR for deployment. The function is implemented by the `mbt assemble` command, which copies all modules, module required dependencies, module build artifacts, and resource definitions into an MTAR, according to the MTA Deployment Descriptor (mtad.yaml file).

If a module doesn’t have the `path` attribute, the assemble process skips it. If a module has the `path` attribute, but its value is invalid, the assemble process fails with the following error: `ERROR could not copy MTA artifacts to assemble: the "<path value>" path does not exist in the MTA project location`.

If a required dependency or resource doesn’t have the `parameters.path` attribute, the assemble process skips it. If a required dependency or resource has the `parameters.path` attribute, but its value is invalid, the assemble process fails with the following error: `ERROR could not copy MTA artifacts to assemble: the "<parameters.path value>" path does not exist in the MTA project location`.

For more details about how to use the MBT assemble command, see [How to build an MTA archive from the modules' build artifacts](https://github.com/SAP/cloud-mta-build-tool/blob/master/docs/docs/usage.md#how-to-build-an-mta-archive-from-the-modules-build-artifacts).

