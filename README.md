# MTA Build Tool (UNDER DEVELOPMENT!)


- The mta command-line tool provides a convenient way to build an MTA project into an MTAR (MTA Archive). 

This includes:
- Generate a Makefile as manifest that describe the build process.
- Building each of the modules in the MTA project.
- Build META-INF folder with the following content:
  - Translating the mta.yaml source file into the mtad.yaml deployment descriptor.
  - Create META-INFO file which describe the build artifacts structure.
- Provide atomic command's that can be executed as isolated process
- Packaging the results into an MTAR file




### Commands

Following is the command which the MBT support:


| Command | usage      | description                                            |
| ------  | ------     |  ----------                                            |
| version | mbt -v     | Prints the MBT version                                 |
| help    | mbt -h     | Prints all the available commands                      | 
| init    | mbt init   | Generate Makefile according to the mta.yaml            |
| ...     | ....       | Will be provided soon                                  |



#### Todos

 - Support MVP scenario 
 - Add comprehensive tests
 - Release process
 - Usage
 - Add concrete limitations
 
 
 #### Limitations
 
   - Under development!
 
 