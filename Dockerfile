FROM openjdk:8-jdk-slim

ARG MTA_USER_HOME=/home/mta
ARG MBT_VERSION=0.3.4
ARG NODE_VERSION=v12.13.0
ARG MAVEN_VERSION=3.6.2

ENV PYTHON /usr/bin/python2.7
ENV M2_HOME=/opt/maven/apache-maven-${MAVEN_VERSION}

RUN apt-get update && \
    apt-get install --yes --no-install-recommends curl git  && \

    # install node
    NODE_HOME=/opt/nodejs; mkdir -p ${NODE_HOME} && \
    curl --fail --silent --output - "http://nodejs.org/dist/${NODE_VERSION}/node-${NODE_VERSION}-linux-x64.tar.gz" \
     | tar -xzv -f - -C "${NODE_HOME}" && \
    ln -s "${NODE_HOME}/node-${NODE_VERSION}-linux-x64/bin/node" /usr/local/bin/node && \
    ln -s "${NODE_HOME}/node-${NODE_VERSION}-linux-x64/bin/npm" /usr/local/bin/npm && \
    ln -s "${NODE_HOME}/node-${NODE_VERSION}-linux-x64/bin/npx" /usr/local/bin/ && \
    # config NPM
    npm config set @sap:registry https://npm.sap.com --global && \
    echo "[INFO] installing maven." && \

    # update maven home
     M2_BASE="$(dirname ${M2_HOME})" && \
    mkdir -p "${M2_BASE}" && \
    curl --fail --silent --output - "https://apache.osuosl.org/maven/maven-3/${MAVEN_VERSION}/binaries/apache-maven-${MAVEN_VERSION}-bin.tar.gz" \
    | tar -xzvf - -C "${M2_BASE}" && \
     ln -s "${M2_HOME}/bin/mvn" /usr/local/bin/mvn && \
     chmod --recursive a+w "${M2_HOME}"/conf/* && \

     # install grunt
     npm install -g grunt-cli && \

     # Download MBT
     curl -L "https://github.com/SAP/cloud-mta-build-tool/releases/download/v${MBT_VERSION}/cloud-mta-build-tool_${MBT_VERSION}_Linux_amd64.tar.gz" | tar -zx -C /usr/local/bin && \
     chown root:root /usr/local/bin/mbt && \

     # handle users permission
     useradd --home-dir "${MTA_USER_HOME}" \
                 --create-home \
                 --shell /bin/bash \
                 --user-group \
                 --uid 1000 \
                 --comment 'Cloud MTA Build Tool' \
                 --password "$(echo weUseMta |openssl passwd -1 -stdin)" mta && \
         # allow anybody to write into the images HOME
         chmod a+w "${MTA_USER_HOME}" && \

    # Install essential build tools and python, required for building db modules
     apt-get install --yes --no-install-recommends \
           build-essential \
           python2.7 && \
    apt-get remove --purge --autoremove --yes \
      curl && \
    rm -rf /var/lib/apt/lists/*

ENV PATH=$PATH:./node_modules/.bin HOME=${MTA_USER_HOME}

WORKDIR /project

USER mta

RUN mbt -v





