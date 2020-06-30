FROM alpine:3.12.0

ARG MBT_VERSION=1.0.15

ADD http://aia.pki.co.sap.com/aia/SAP%20Global%20Root%20CA.crt \
    /etc/ssl/certs/SAP_Global_Root_CA.crt

RUN apk add --no-cache --update \
    ca-certificates curl && \
    curl -L "https://github.com/SAP/cloud-mta-build-tool/releases/download/v${MBT_VERSION}/cloud-mta-build-tool_${MBT_VERSION}_Linux_amd64.tar.gz" | tar -zx -C /usr/local/bin && \
    chown root:root /usr/local/bin/mbt


