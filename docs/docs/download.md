## Download and Installation

 There are two supported ways to install the multi-target application archive builder (MBT) tool:

 **Manually**:
 
  1. [Download](https://github.com/SAP/cloud-mta-build-tool/releases) the **latest** binary file according to your operating system.
  2. Extract the archive file to the folder where you want to install the tool.
  3. Add the binary file to your `~/bin` path according to your operating system:     
        - Darwin / Linux
          - Copy the binary file to the `~/usr/local/bin/` folder, for example: `cp mbt /usr/local/bin/`
        - Windows
          -  Copy the binary file `mbt.exe` to the `C:/Windows/` folder.

**Use npm**:
```
npm install -g mbt
```