### Supported commands

| Command | Usage &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;       | Description                                                    
| ------  | --------       |  ----------                                                    
| version | `mbt -v`     | Prints the Multi-Target Application Archive Builder tool version.                                        | x
| help    | `mbt -h`     | Prints all the available commands.                             
| assemble    | `mbt assemble`     | Creates an MTA archive `.mtar` file from the module build artifacts according to the MTA deployment descriptor (`mtad.yaml` file). Runs the command in the directory where the `mtad.yaml` file is located. <br>**Note:** Make sure the path property of each module's `mtad.yaml` file points to the module's build artifacts that you want to package into the target MTA archive. 
| init    | `mbt init`     | Generates the `Makefile.mta` file according to the MTA descriptor (`mta.yaml` file or `mtad.yaml` file). <br> The `make` command uses the generated `Makefile.mta` file to package the MTA project. 
<br>
For more information, see the command help output available via either of the following:

- `mbt [command] --help` 
- `mbt [command] -h`
