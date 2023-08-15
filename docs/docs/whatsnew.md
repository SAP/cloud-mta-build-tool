# <b>Important updates</b>

## v1.2.25

### Configuration of `maven` builder has changed. 
As of version 1.2.25, the `mvn -B clean package` command is used where the `maven` builder is configured for building a module or in a global build step. The `-B` parameter is added to `maven` builder's build command which will start maven in a interactive mode.

As a build tool, MBT packages all modules’s build process internal, all native builders, such as golang, maven and nodejs, run build command in batch mode is resonable.

<b>NOTE:</b>  The `maven` builder is configured implicitly for the `java` module type.

If you want to keep the previous behavior, that is, to apply the `mvn clean package` command, you need to change the build parameters of the relevant module in one of the following ways:

<b>Examples:</b> Configure the `custom` builder.
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
