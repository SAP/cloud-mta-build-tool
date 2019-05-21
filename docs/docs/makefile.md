### Overview


In order to build an MTA archive from a project's source code you should have `GNU Make 4.2.1` installed in your build environmnet. Then you can use the tool's  `mbt init` command that generates the `Makefile.mta` base on the project's MTA development desriptor `mta.yaml` file. The `Makefile.mta` file is used by the `make` command and provides the verbose build manifest, which can be changed according to the project needs. It is responsible for:

- Building each of the modules in the MTA project.
- Invoking the MBT commands in the correct order.
- Archiving the MTA project.<br>

For more details about the commands, see the `Usage` section.


#### Tip

For Windows, use [Chocolatey](https://chocolatey.org/packages/make) to install or upgrade `GNU Make`.
