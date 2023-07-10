#### Overview
[Micromatch](https://github.com/micromatch/micromatch) support various matching features, such as glob patterns and advanced globbing

Micromatch is a Node.js application, while MBT is a Go application. The micromatch wrapper is a package that enables MBT to use the functionalities provided by micromatch. From MBT 1.2.25 version, we create [micromath wrapper](https://github.com/SAP/cloud-mta-build-tool/tree/master/micromatch) to package it.

#### Install
You can install the micromatch wrapper using one of the methods below:

 **Download and install manually**

   - [Download](https://github.com/SAP/cloud-mta-build-tool/releases) the latest binary file according to your operating system.
    
   - Extract the archive file to the folder where you want to install the tool.

   - Add the binary file to your `~/bin` path according to your operating system:  

     * In Darwin / Linux, copy the micromatch-wrapper binary file to the `~/usr/local/bin/` folder. Here is a sample command: `cp micromatch-wrapper /usr/local/bin/`

     * In Windows, copy the `micromatch-wrapper.exe` binary file to the `C:/Windows/` folder.

**Build and install from the source code**

   - [Download](https://github.com/SAP/cloud-mta-build-tool/releases) the latest source code.

   - Install [pkg](https://github.com/vercel/pkg/) on your operating system:
```
npm install -g pkg
```

   - Build the micromatch wrapper under the `micromatch` subfolder of the source code:
```
cd micromatch
npm install
pkg ./
```
   - Add the binary file to your `~/bin` path according to your operating system:  

     * In Darwin, copy the micromatch-wrapper binary file to the `~/usr/local/bin/` folder. Here is a sample command: `cp micromatch-wrapper-macos /usr/local/bin/micromatch-wrapper`

     * In Linux, copy the micromatch-wrapper binary file to the `~/usr/local/bin/` folder. Here is a sample command: `cp micromatch-wrapper-linux /usr/local/bin/micromatch-wrapper`

     * In Windows, copy the micromatch-wrapper binary file to the `C:/Windows/` folder.

**Install using npm**

Run the command below to install MBT and the micromatch wrapper together:

```
npm install -g mbt@greater_than_or_equal_to_v1.2.25
mbt --version
micromatch-wrapper --version
```