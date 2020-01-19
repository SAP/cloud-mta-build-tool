## Troubleshooting

### Installation

**Indicator:**
Installation fails with the following error:
```Error: EACCES: permission denied```
 
**Solution:**
Grant  the user "admin" permissions to run the command for the installation process using the following command:

```sudo npm install -g mbt --unsafe-perm=true --allow-root```

### Building a Multitarget Application

#### Timeouts
##### Missing npm registry configuration
**Indicator:**
Sometimes, when executing the `mbt build [args]` command, the build process hangs and ends with a timeout error.
This may be due to a network error or because of missing registry configurations.
When using packages provided by SAP, you should verify that your `npm config` file refers to the SAP registry.
   
**Solution:**
Add an `.npmrc` file to the module (that hangs) as a sibling to the 'package.json' file as follows:
```
  $ cat .npmrc

  @sap:registry=https://npm.sap.com/
```

For more details, see the npm [documentation](https://docs.npmjs.com/files/npmrc).
  
#### Make cannot run on Mac OS
**Indicator:**
The `mbt build` command fails on Mac and the build output contains the following error:
```xcrun: error: invalid active developer path (/Library/Developer/CommandLineTools), missing xcrun at: /Library/Developer/CommandLineTools/usr/bin/xcrun```

**Solution:**
Install the Command-Line Tools:
```xcode-select --install```
The tools can also be downloaded from [https://developer.apple.com/download/more/](https://developer.apple.com/download/more/).

