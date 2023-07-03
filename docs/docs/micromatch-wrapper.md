#### Overview
[Micromatch](https://github.com/micromatch/micromatch) support various matching features, such as Glob pattern, Advanced globbing.

Micromatch is a nodejs application, and MBT is a go application. In order to enable MBT to use the functionalities provided by micromatch, we create [micromath wrapper](https://github.com/SAP/cloud-mta-build-tool/tree/master/micromatch) to package it.

#### Install
You can install the Micromatch Wrapper using either of these methods below:

 **Download and install manually**

   - [Download](https://github.com/SAP/cloud-mta-build-tool/releases) the latest binary file according to your operating system.
    
   - Extract the archive file to the folder where you want to install the tool.

   - Add the binary file to your `~/bin` path according to your operating system:  

     * In Darwin / Linux, copy binary file `micromatch-wrapper` to the `~/usr/local/bin/` folder, for example: `cp micromatch-wrapper /usr/local/bin/`

     * In Windows, copy the `micromatch-wrapper.exe` binary file to the `C:/Windows/` folder.

**Build and install manually**

   - [Download](https://github.com/SAP/cloud-mta-build-tool/releases) the latest source code.

   - Install [pkg](https://github.com/vercel/pkg/) on your operating system
```
npm install -g pkg
```

   - Build micromatch wrapper under `micromatch` subfolder of source code
```
cd micromatch
npm install
pkg ./
```
   - Add the binary file to your `~/bin` path according to your operating system:  

     * In Darwin, copy binary file `micromatch-wrapper` to the `~/usr/local/bin/` folder, for example: `cp micromatch-wrapper-macos /usr/local/bin/micromatch-wrapper`

     * In Linux, copy binary file `micromatch-wrapper` to the `~/usr/local/bin/` folder, for example: `cp micromatch-wrapper-linux /usr/local/bin/micromatch-wrapper`

     * In Windows, copy the `micromatch-wrapper.exe` binary file to the `C:/Windows/` folder.

**Install using npm**

Run the command below, it will install mbt and micromatch-wrapper together

```
npm install -g mbt
mbt --version
micromatch-wrapper --version
```