
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
mbt build -p=cf

```

<b>`mbt build`</b>
Generates a temporary `Makefile` according to the MTA descriptor and runs the `make` command to package the MTA project into the MTA archive.

<b>Usage:</b> `mbt build <flags>`

<b>Flags:</b>

| Flag        | Mandatory&nbsp;/<br>Optional        | Description&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;                 | Examples&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;                                    
| -----------  | -------       |  ----------                          |  -----------------------------
| `-p (--platform)`   | Mandatory  | The name of the target deployment platform. <br>The supported deployment platforms are: <ul><li>`cf` for SAP Cloud Platform Cloud Foundry environment  <li>`neo` for the SAP Cloud Platform Neo environment <li>`xsa` for the SAP HANA XS advanced model                              | `mbt build -p=cf`
| `-s (--source)`   | Optional  | The path to the MTA project; the current path is set as the default.                              | `mbt build -p=cf -s=C:/TestProject`
| `-t (--target)`   | Optional  | The folder for the generated `MTAR` file. If this parameter is not provided, the `MTAR` file is saved in the `mta_archives` subfolder of the current folder. If the parameter is provided, the `MTAR` file is saved in the root of the folder provided by the argument.  | `mbt build -p=cf -t=C:/TestProject`
| `--mtar`   | Optional  | The file name of the generated archive file. If this parameter is omitted, the file name is created according to the following naming convention: <br><br> `<mta_application_ID>_<mta_application_version>.mtar` <br><br> If the parameter is provided, but does not include an extension, the `.mtar` extension is added.  | `mbt build -p=cf --mtar=TestProject.mtar`
| `--strict`   | Optional  | The default value is `true`. If set to `true`, the duplicated fields and fields that are not defined in the `mta.yaml` schema are reported as errors. If set to `false`, they are reported as warnings.  | `mbt build -p=cf --strict=true`


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



&nbsp;

<b>`make`</b>

Packages the MTA project into the MTA archive according to the `Makefile`.

<b>Usage:</b> `make <parameters>`

<b>Parameters:</b>

| Parameter        | Type | Mandatory&nbsp;/<br>Optional        | Description&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;                 | Examples    
| -----------  | ------ | -------       |  ----------                              |  --------------------------------------
| `-f <path to Makefile.mta>`    | string     | Mandatory  | The path to the `Makefile.mta` file that contains the build configurations.                             | `make -f Makefile.mta p=cf`
| `p`  | string     | Mandatory     | The name of the target deployment platform. <br>The supported deployment platforms are: <ul><li>`cf` for the SAP Cloud Platform Cloud Foundry environment  <li>`neo` for the SAP Cloud Platform Neo environment <li>`xsa` for the SAP HANA XS advanced model                                     |`make -f Makefile.mta p=cf`
| `t`    | string     | Optional  | The folder for the generated `MTAR` file. If this parameter is not provided, the `MTAR` file is saved in the `mta_archives` subfolder of the current folder. If the parameter is provided, the `MTAR` file is saved in the root of the folder provided by the argument.                              | `make -f Makefile.mta p=cf t=C:\temp`
| `mtar`    | string     | Optional  | The file name of the generated archive file. If this parameter is omitted, the file name is created according to the following naming convention: <br><br> `<mta_application_ID>_<mta_application_version>.mtar` <br><br> If the parameter is provided, but does not include an extension, the `.mtar` extension is added. | `make -f Makefile.mta p=cf mtar=myMta`<br><br> `make -f Makefile.mta p=cf mtar=myMta.mtar`
| `strict`    | Boolean     | Optional    | The default value is `true`. If set to `true`, the duplicated fields and fields that are not defined in the `mta.yaml` schema are reported as errors. If set to `false`, they are reported as warnings. | `make -f Makefile.mta p=cf strict=false`

&nbsp;
### How to build an MTA archive from the modules' build artifacts 

| Command | Usage &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;       | Description                                                    
| ------  | --------       |  ----------                                                
| assemble    | `mbt assemble`     | Creates an MTA archive `.mtar` file from the module build artifacts according to the MTA deployment descriptor (`mtad.yaml` file). Runs the command in the directory where the `mtad.yaml` file is located. <br>**Note:** Make sure the path property of each module's `mtad.yaml` file points to the module's build artifacts that you want to package into the target MTA archive. 
