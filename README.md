# Cloud MTA Build Tool (Beta)


The Cloud MTA Build Tool is a command-line tool that packages a multitarget application into a deployable archive (MTAR). For full documentation see [Cloud MTA build Tool](https://hub.docker.com/r/devxci/mbtci).


This image can be used to build SAP Multitarget Applications (MTA) containing Java, Node.js, and Golang modules. The image is hosted at hub.docker.com.



## How to use the image



On a Linux/Darwin machine you can run:

```

docker run -it --rm -v "$(pwd)/[proj-releative-path]:/project" devxci/mbtci:0.1.0 mbt build -p=cf -t [target-folder-name]

```


This will build an mtar file for SAP Cloud Platform (Cloud Foundry). The folder containing the project needs to be mounted into the image at /project.


<b>Note:</b> The parameter `-p=cf` can be omitted as the build for cloud foundry is the default build, this is an example of the MBT build parameters, for further commands see MBT docs.


## How to build the image

```

docker build -t devxci/mbtci .

```

## The image provides:


- Cloud MTA Build Tool - 1.0.0

- Nodejs - 12.13.0

- SAP registry (@sap:registry https://npm.sap.com) contained in global node configuration.

- Maven - 3.6.2

- Golang - 1.13.4

- Java - 8



The MTA Archive Builder delegates module builds to other native build tools. This image provides Node.js, Maven, Java, and Golang so the archive builder can delegate to these build technologies. In case other build tools are needed, <b>inherit</b> from this image and add more build tools.



### License


Copyright (c) 2019 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file.
Please note that Docker images can contain other software which may be licensed under different licenses. This License file is also included in the Docker image. For any usage of built Docker images please make sure to check the licenses of the artifacts contained in the images.