## Troubleshooting

### Installation

**Indicator:**
Installation fails with the following error:
```Error: EACCES: permission denied```
 
**Solution:**
Grant  the user "admin" permissions to run the command for the installation process using the following command:

```sudo npm install -g mbt --unsafe-perm=true --allow-root```

### Building a Multitarget Application
  
#### Make cannot run on Mac OS
**Indicator:**
The `mbt build` command fails on Mac and the build output contains the following error:
```xcrun: error: invalid active developer path (/Library/Developer/CommandLineTools), missing xcrun at: /Library/Developer/CommandLineTools/usr/bin/xcrun```

**Solution:**
Install the Command-Line Tools:
```xcode-select --install```
The tools can also be downloaded from [https://developer.apple.com/download/more/](https://developer.apple.com/download/more/).

#### Make package collides with Make GNU ####
**Indicator:**
The `mbt build` command fails with the following error:
```make ✖ err missing makefile / bakefile```

**Solution:**
Remove the global make package:
```npm uninstall -g make```
