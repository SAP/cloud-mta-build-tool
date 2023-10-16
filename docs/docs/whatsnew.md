# <b>Important updates</b>

## v1.2.26

### Support SBOM generation
The Software Bill of Materials (SBOM) is a list of components, libraries, and module information that are required to build a software, and the supply chain relationships between them. An SBOM also lists the licenses that govern those components, versions of the components used in the codebase, and their patch status.

With SBOM, teams can quickly identify any associated security or license risks of codebase.

The `npm, maven, and golang` native builders and the `mbt build, mbt sbom-gen` commands have been upgraded to support SBOM generation. For `java or nodejs` module types, or if the module's `build-parameters.builder` attribute value is `npm, maven, or golang`, SBOM content will be generated and merged into one file. Currently, only the XML format SBOM file is supported.

The module configuration can be referenced in the [configuration.md](https://github.com/SAP/cloud-mta-build-tool/blob/master/docs/docs/configuration.md) file.

The SBOM generation commands `mbt build` and `mbt sbom-gen` can be referenced in the [usage.md](https://github.com/SAP/cloud-mta-build-tool/blob/master/docs/docs/usage.md) file.

## v1.2.25

### Configuration of `maven` builder has changed. 
As of version 1.2.25, the `mvn -B clean package` command is used where the `maven` builder is configured for building a module or in a global build step. 

By adding the -B parameter, the "maven clean package" command will start in interactive mode.

As a build tool, MBT builds the MTA (Multitarget Application), which contains many types of modules. Each module is built by corresponding builder, such as maven, npm and golang.

When the MTA is built, all build processes for different modules are packaged into the MBT. It wraps all the build processes as internal and it should not run in interactive mode. It is reasonable for MBT to execute the build process in batch mode.

<b>NOTE:</b>  The `maven` builder is configured implicitly for the `java` module type.

If you want to keep the previous behavior, that is, to apply the `mvn clean package` command, you need to change the build parameters of the relevant module by configuring the `custom` builder:
```yaml

- name: mymodule
  ... 
  build-parameters:
      builder: custom
      commands:
        - mvn clean package
      build-result: target/*.war 
```

## v1.1.0 

### Configuration of `maven` builder has changed. 
As of version 1.1.0, the `mvn clean package` command is used where the `maven` builder is configured for building a module or in a global build step.

<b>NOTE:</b>  The `maven` builder is configured implicitly for the `java` module type.

If you want to keep the previous behavior, that is, to apply the `mvn -B package` command, you can use the `maven_deprecated` builder or `custom` builder as shown in the examples below.

<b>NOTE:</b> The `maven_deprecated` builder will be removed on July 2021.

<b>Examples:</b>

If you want to use the old `mvn -B package` command instead of the `maven` builder that now triggers the `mvn clean package` command, you need to change the build parameters of the relevant module in one of the following ways:



<b>Option 1:</b> Set `maven_deprecated` as the module builder parameter.

```yaml

- name: mymodule
  ... 
  build-parameters:
      builder: maven_deprecated
      
```

<b>Option 2:</b> Configure the `custom` builder.
```yaml

- name: mymodule
  ... 
  build-parameters:
      builder: custom
      commands:
        - mvn -B package
      build-result: target/*.war 
```

The same approach can be implemented if the `maven` builder is used in the global build step.
