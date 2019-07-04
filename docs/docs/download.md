
 You can install the Cloud MTA Build Tool (MBT) using either of these methods below:

 **Install manually**
 &nbsp;
  1. [Download](https://github.com/SAP/cloud-mta-build-tool/releases) the **latest** binary file according to your operating system.
      * Darwin / Linux
``` 
//For example:      
wget https://github.com/SAP/cloud-mta-build-tool/releases/download/<LATEST>/cloud-mta-build-tool_<LATEST>_Linux_amd64.tar.gz 
```

  2. Extract the archive file to the folder where you want to install the tool.
      * Darwin / Linux
```
//For example:
tar xvzf cloud-mta-build-tool_LATEST_Linux_amd64.tar.gz
```
  3. Add the binary file to your `~/bin` path according to your operating system:  &nbsp;   
        * Darwin / Linux
          Copy the binary file to the `~/usr/local/bin/` folder, for example: `cp mbt /usr/local/bin/`
&nbsp;
        * Windows
          Copy the `mbt.exe` binary file to the `C:/Windows/` folder.

**Install using npm**

Run the command below.

```
npm install -g mbt
```

