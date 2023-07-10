
 You can install the Cloud MTA Build Tool (MBT) using either of these methods below:

 **Install manually**

   - [Download](https://github.com/SAP/cloud-mta-build-tool/releases) the latest binary file according to your operating system.

```
    // Example for Darwin/Linux: 
    wget https://github.com/SAP/cloud-mta-build-tool/releases/download/<LATEST>/cloud-mta-build-tool_<LATEST>_Linux_amd64.tar.gz
```
    
   - Extract the archive file to the folder where you want to install the tool.


```
//Example for Darwin/Linux:
  tar xvzf cloud-mta-build-tool_LATEST_Linux_amd64.tar.gz
```

   - Add the binary file to your `~/bin` path according to your operating system:  

     * In Darwin / Linux, copy the binary file to the `~/usr/local/bin/` folder, for example: `cp mbt /usr/local/bin/`

     * In Windows, copy the `mbt.exe` binary file to the `C:/Windows/` folder.

**Install using npm**

Run the command below.

```
npm install -g mbt
```

From MBT 1.2.25 version, the `build-parameters ignore` attribute will be upgraded to support [full glob patterns](https://en.wikipedia.org/wiki/Glob_(programming)). You must install the micromatch wrapper to use the glob patterns.

You can install micromatch-wrapp by instruction [Here](https://github.com/SAP/cloud-mta-build-tool/tree/master/docs/docs/micromatch-wrapper.md)