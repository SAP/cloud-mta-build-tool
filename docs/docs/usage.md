
### How to get help on the tools commands

| Command | Usage &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;       | Description                                                    
| ------  | --------       |  ----------                                                
| help    | `mbt -h`     | Prints all the available commands.                           
| help    | `mbt [command] --help` or<br> `mbt [command] -h`    | Prints detailed information about the specified command.|

&nbsp;
### How to find out the version of the installed tool

| Command | Usage &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;       | Description                                                    
| ------  | --------       |  ----------                                                
| version | `mbt -v`     | Prints the current Cloud MTA Build Tool version.                                        <br>

&nbsp;
### How to build an MTA archive from the project sources

#### Overview

You can use one of the following two approaches for building your MTA project:
- One-step build using the `mbt build` command
- Two-step build using a combination of the `mbt init` and `make` commands.


Both methods leverage the `GNU Make` technology for the actual build.  <br>
If you are using the one-step approach, the tool generates a temporary build configuration file and automatically invokes the `make` command. The generated `Makefile` is then deleted at the end of the build.  <br>
The second approach allows you to generate the `Makefile` using the `mbt init` command. You can adjust the generated file according to your project needs and then build the MTA archive using the `make` command. In this case, we recommend that you include the generated `Makefile` in the project's source control management system to ensure that the same build process is applied across all the project's contributors, regardless of the build environment. 


#### Prerequisites
* `GNU Make 4.2.1` is installed in your build environment. 
* Module build tools are installed in your build environment.

For more information, see the corresponding [`Download` and `Installation` sections](download.md).

#### One-step build

<b> Quick start example:</b>

```go
// Executes the MTA project build for the Cloud Foundry target environment.
mbt build

```

<b>`mbt build`</b>
Generates a temporary `Makefile` according to the MTA descriptor and runs the `make` command to package the MTA project into the MTA archive.

<b>Usage:</b> `mbt build <flags>`

<b>Flags:</b>

| Flag        | Mandatory&nbsp;/<br>Optional        | Description&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;                 | Examples&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;                                    
| -----------  | -------       |  ----------                          |  -----------------------------
| `-p (--platform)`   | Optional  | The name of the target deployment platform. <br>The supported deployment platforms are: <ul><li>`cf` for SAP Cloud Platform, Cloud Foundry environment  <li>`neo` for the SAP Cloud Platform, Neo environment <li>`xsa` for the SAP HANA XS advanced model</ul> If this parameter is not provided, the project is built for the SAP Cloud Platform, Cloud Foundry environment                             | `mbt build -p=cf`
| `-s (--source)`   | Optional  | The path to the MTA project; the current path is set as the default.                              | `mbt build -p=cf -s=C:/TestProject`
| `-t (--target)`   | Optional  | The folder for the generated `MTAR` file. If this parameter is not provided, the `MTAR` file is saved in the `mta_archives` subfolder of the current folder. If the parameter is provided, the `MTAR` file is saved in the root of the folder provided by the argument.  | `mbt build -p=cf -t=C:/TestProject`
| `--mtar`   | Optional  | The file name of the generated archive file. If this parameter is omitted, the file name is created according to the following naming convention: <br><br> `<mta_application_ID>_<mta_application_version>.mtar` <br><br> If the parameter is provided, but does not include an extension, the `.mtar` extension is added.  | `mbt build -p=cf --mtar=TestProject.mtar`
| `-e (--extensions)`   | Optional  | The path or paths to multitarget application extension files (.mtaext). Several extension files separated by commas can be passed with a single flag, or each extension file can be specified with its own flag.  |`mbt build -e=test1.mtaext,test2.mtaext`<br>or<br>`mbt build -e=test1.mtaext -e=test2.mtaext`
| `--strict`   | Optional  | The default value is `true`. If set to `true`, the duplicated fields and fields that are not defined in the `mta.yaml` schema are reported as errors. If set to `false`, they are reported as warnings.  | `mbt build -p=cf --strict=true`
| BETA &nbsp;&nbsp;`-m (--mode)`   | Optional  | The possible value is `verbose`. If run with this option, the temporary `Makefile` is generated in a way that allows the parallel execution of `Make` jobs to make the build process faster.   | `mbt build -m=verbose`
| BETA  &nbsp;&nbsp;`-j (--jobs)`   | Optional  | Used only with the `--mode` parameter. This option configures the number of `Make` jobs that can run simultaneously. If omitted or if the value is less than or equal to zero, the number of jobs is defined by the number of available CPUs (maximum 8).    | `mbt build -m=verbose -j=8`


&nbsp;

#### Two-step build

<b> Quick start example:</b>

```go
// Generates the `Makefile.mta` file.
mbt init 

// Executes the MTA project build for Cloud Foundry target environment.
make -f Makefile.mta p=cf

```

<b>`mbt init`</b>

Generates the `Makefile.mta` file according to the MTA descriptor (mta.yaml file). The `make` command uses the generated `Makefile.mta` file to package the MTA project. <br>


&nbsp;

<b>Usage:</b> `mbt init <flags>`

<b>Flags:</b>

| Flag        | Mandatory&nbsp;/<br>Optional        | Description&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;                 | Examples&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;                                    
| -----------  | -------       |  ----------                          |  -----------------------------
| `-s (--source)`   | Optional  | The path to the MTA project; the current path is set as the default.                              | `mbt init -s=C:/TestProject`
| `-t (--target)`   | Optional  | The path to the generated `Makefile` folder; the current path is set as the default.   | `mbt init -t=C:/TestFolder`
| `-e (--extensions)`   | Optional  | The path or paths to multitarget application extension files (.mtaext). Several extension files separated by commas can be passed with a single flag, or each extension file can be specified with its own flag.    | `mbt build -e=test1.mtaext,test2.mtaext`<br>or<br>`mbt build -e=test1.mtaext -e=test2.mtaext`



&nbsp;

<b>`make`</b>

Packages the MTA project into the MTA archive according to the `Makefile`.

<b>Usage:</b> `make <parameters>`

<b>Parameters:</b>

| Parameter        | Type | Mandatory&nbsp;/<br>Optional        | Description&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;                 | Examples    
| -----------  | ------ | -------       |  ----------                              |  --------------------------------------
| `-f <path to Makefile.mta>`    | string     | Mandatory  | The path to the `Makefile.mta` file that contains the build configurations.                             | `make -f Makefile.mta p=cf`
| `p`  | string     | Mandatory     | The name of the target deployment platform. <br>The supported deployment platforms are: <ul><li>`cf` for the SAP Cloud Platform, Cloud Foundry environment  <li>`neo` for the SAP Cloud Platform, Neo environment <li>`xsa` for the SAP HANA XS advanced model                                     |`make -f Makefile.mta p=cf`
| `t`    | string     | Optional  | The folder for the generated `MTAR` file. If this parameter is not provided, the `MTAR` file is saved in the `mta_archives` subfolder of the current folder. If the parameter is provided, the `MTAR` file is saved in the root of the folder provided by the argument.                              | `make -f Makefile.mta p=cf t=C:\temp`
| `mtar`    | string     | Optional  | The file name of the generated archive file. If this parameter is omitted, the file name is created according to the following naming convention: <br><br> `<mta_application_ID>_<mta_application_version>.mtar` <br><br> If the parameter is provided, but does not include an extension, the `.mtar` extension is added. | `make -f Makefile.mta p=cf mtar=myMta`<br><br> `make -f Makefile.mta p=cf mtar=myMta.mtar`
| `strict`    | Boolean     | Optional    | The default value is `true`. If set to `true`, the duplicated fields and fields that are not defined in the `mta.yaml` schema are reported as errors. If set to `false`, they are reported as warnings. | `make -f Makefile.mta p=cf strict=false`

&nbsp;
### How to build an MTA archive from the modules' build artifacts 





<b>`mbt assemble`</b>

Creates an MTA archive `MTAR` file from the module build artifacts according to the MTA deployment descriptor (`mtad.yaml` file). 
> <b>Note</b>: Make sure the path property of each module's `mtad.yaml` file points to the module's build artifacts that you want to package into the target MTA archive.

<b>Usage:</b> `mbt assemble <flags>`

<b>Flags:</b>

| Flag        | Mandatory&nbsp;/<br>Optional        | Description&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;                 | Examples&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;                                    
| -----------  | -------       |  ----------                          |  -----------------------------
| `-s (--source)`   | Optional  | The path to the folder where the project’s `mtad.yaml` file is located; the current path is set as the default.                              | `mbt assemble  -s=C:/TestProject`
| `-t (--target)`   | Optional  | The folder for the generated `MTAR` file. If this parameter is not provided, the `MTAR` file is saved in the `mta_archives` subfolder of the current folder. If the parameter is provided, the `MTAR` file is saved in the root of the folder provided by the argument.  | `mbt assemble  -t=C:/TestFolder`
| `-m (--mtar)`   | Optional  | The name of the generated archive file. If this parameter is omitted, the file name is created according to the following naming convention: <br><br> `<mta_application_ID>_<mta_application_version>.mtar` <br><br> If the parameter is provided, but does not include an extension, the `.mtar` extension is added.  | `mbt assemble  -m=anotherName`
| `-e (--extensions)`   | Optional  | The path or paths to multitarget application extension files (`.mtaext`). Several extension files separated by commas can be passed with a single flag, or each extension file can be specified with its own flag.| `mbt assemble -e=test1.mtaext,test2.mtaext`<br>or<br>`mbt assemble -e=test1.mtaext -e=test2.mtaext`


&nbsp;
### Auxiliary commands  

This section is dedicated for commands that execute specific steps of the MTA build process such as project validation, build for a single module, and generation of the deployment descriptor. These commands are useful if, for example, you want to build and deploy only specific modules for testing purposes, or if you decide to tailor your own build process for packaging MTA archives.  At the moment, only the commands described below are supported. 

<b>`mbt mtad-gen`</b>

Generates the MTA deployment descriptor (`mtad.yaml` file) according to the provided MTA descriptor (`mta.yaml` file) and MTA extensions. 


<b>Usage:</b> `mbt mtad-gen <flags>`

<b>Flags:</b>

| Flag        | Mandatory&nbsp;/<br>Optional        | Description&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;                 | Examples&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;                                    
| -----------  | -------       |  ----------                          |  -----------------------------
| `-p (--platform)`   | Optional  | The name of the target deployment platform. <br>The supported deployment platforms are: <ul><li>`cf` for SAP Cloud Platform, Cloud Foundry environment  <li>`neo` for the SAP Cloud Platform, Neo environment <li>`xsa` for the SAP HANA XS advanced model</ul> If this parameter is not provided, the project is built for the SAP Cloud Platform, Cloud Foundry environment                             | `mbt mtad-gen -p=cf`
| `-s (--source)`   | Optional  | The path to the folder where the project’s `mta.yaml` file is located; the current path is set as default.                              | `mbt mtad-gen  -s=C:/TestProject`
| `-t (--target)`   | Optional  | The folder where the `mtad.yaml` will be generated. If this parameter is not provided, the `mtad.yaml` file is saved in the current folder. If the parameter is provided, the generated file is saved in the root of the folder provided by the argument.  | `mbt mtad-gen  -t=C:/TestFolder`
| `-e (--extensions)`   | Optional  | The path or paths to multitarget application extension files (`.mtaext`). Several extension files separated by commas can be passed with a single flag, or each extension file can be specified with its own flag.| `mbt mtad-gen -e=test1.mtaext,test2.mtaext`<br>or<br>`mbt mtad-gen -e=test1.mtaext -e=test2.mtaext`

<br>
<br>

<b>`mbt module-build`</b>

Triggers the build process of the specified module according to the implicit or explicit build configurations in the MTA descriptor (`mta.yaml` file) and MTA extensions.

<b>Usage:</b> `mbt module-build <flags>`

<b>Flags:</b>

| Flag        | Mandatory&nbsp;/<br>Optional        | Description&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;                 | Examples&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;                                    
| -----------  | -------       |  ----------                          |  -----------------------------
| `m (--module)`   | Mandatory  | The name of the modules to build. Several modules separated by commas can be passed with a single flag, or each module can be specified with its own flag.<br><br> <b>Notes</b>:<ul><li>If you specify several modules with a single flag, make sure there are no spaces between the names. <li>The tool builds the specified modules in the order defined by the build parameters in the `mta.yaml` file. If the specified modules depend on other modules that should be built before, run the command with the `--with-all-dependencies` flag or specify the required modules explicitly.<li>If the `--target` parameter is used, make sure that the names of the specified modules' build artifacts are unique. You can configure the build artifact name using the `build-artifact-name` build parameter as described [here](configuration.md#configuring-the-build-artifact-name).  </ul>                        |`mbt module-build -m=my_module,another_module` <br>or<br>`mbt module-build -m=my_module -m=another_module`
| `-a (--with-all-dependencies)`   | Optional  | If this option is used, the tool builds all modules that the specified modules depend on in the order defined by the build parameters in the `mta.yaml` file before building the selected modules.<br>Without this option, only the selected modules are built.                               |`mbt module-build -m=my_module,another_module  -a`
| `-s (--source)`   | Optional  | The path to the folder where the project’s `mta.yaml` file is located; the current path is set as default.                              |`mbt module-build -m=my_module  -s=C:/TestProject`
| `-t (--target)`   | Optional  | The folder where the module build results will be saved. If this parameter is not provided, the module build results are saved in the `<current folder>/.<projectname>_mta_build_tmp/<module name>` folder. If  the parameter is provided, the build results are saved in the root of the folder provided by the argument. <br><br><b>Note:</b><br>If the `--target` parameter is used when building several modules, make sure that the names of the specified modules' build artifacts are unique. You can configure the build artifact name using the `build-artifact-name` build parameter as described [here](configuration.md#configuring-the-build-artifact-name). | `mbt module-build -m=my_module  -t=C:/TestProject/build_results_tmp`<br> <br>The module’s build results will be saved directly in the `C:/TestProject/build_results_tmp/` folder
| `-e (--extensions)`   | Optional  | The path or paths to multitarget application extension files (`.mtaext`). Several extension files separated by commas can be passed with a single flag, or each extension file can be specified with its own flag.| `mbt module-build -m=my_module -e=test1.mtaext,test2.mtaext`<br>or<br>`mbt module-build -m=my_module -e=test1.mtaext -e=test2.mtaext`
| `-g (--mtad-gen)`   | Optional  | If the parameter is provided, the deployment descriptor `mtad.yaml` is generated by default in the current folder or in the folder configured by the `--target` parameter. <br> A module's `path` property in the generated `mtad.yaml` file points to the module's build results if this module was selected using the `--modules` option. <br><br> <b>Notes</b>:<ul><li>The selected module list specified using the `--module` option, does not affect the list of modules in the resulting `mtad.yaml` file. The `mtad.yaml` file is always generated according to the default Cloud MTA Builder settings, the `build-parameters` configurations in the `mta.yaml` file (e.g. `supported-platforms`), and the selected target platform.<li>By default, the `mtad.yaml` is generated for the `cf` target platform. You can configure a different target plaform using the `--platform` option.  | `mbt module-build -m=my_module1, my_module2 –g`
| `-p (--platform)`   | Optional  |  The name of the target deployment platform. Used only with the `-g (--mtad-gen)` parameter. <br>The supported deployment platforms are: <ul><li>`cf` for SAP Cloud Platform, Cloud Foundry environment  <li>`neo` for the SAP Cloud Platform, Neo environment <li>`xsa` for the SAP HANA XS advanced model</ul> If this parameter is not provided, the `mtad.yaml` file is generated for the SAP Cloud Platform, Cloud Foundry environment.                             | `mbt module-build -m=my_module1, my_module2 –g -p=neo`

