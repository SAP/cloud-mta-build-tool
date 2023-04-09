#!/bin/sh

go build

# ./cloud-mta-build-tool.exe build --mode verbose -s /c/Workspace/MTA/cloud-mta-build-tool-bom-support-test-app/go-test-case -b sbom-gen-result/merged.bom.xml --keep-makefile

./cloud-mta-build-tool.exe build -s /c/Workspace/MTA/cloud-mta-build-tool-bom-support-test-app/go-test-case -b sbom-gen-result/merged.bom.xml

./cloud-mta-build-tool.exe build -s /c/Workspace/MTA/cloud-mta-build-tool-bom-support-test-app/nodejs-test-case -b sbom-gen-result/merged.bom.xml

./cloud-mta-build-tool.exe build -s /c/Workspace/MTA/cloud-mta-build-tool-bom-support-test-app/java-test-case -b sbom-gen-result/merged.bom.xml

# ./cloud-mta-build-tool.exe sbom-gen -s /c/Workspace/MTA/cloud-mta-build-tool-bom-support-test-app/go-test-case -b sbom-gen-result/merged.bom.xml

# ./cloud-mta-build-tool.exe sbom-gen -s /c/Workspace/MTA/cloud-mta-build-tool-bom-support-test-app/nodejs-test-case -b sbom-gen-result/merged.bom.xml

# ./cloud-mta-build-tool.exe sbom-gen -s /c/Workspace/MTA/cloud-mta-build-tool-bom-support-test-app/java-test-case -b sbom-gen-result/merged.bom.xml

# ./cloud-mta-build-tool.exe build -s /c/Workspace/MTA/cf-mta-examples/blue-green-deploy-strategy/hello-blue -b sbom-gen-result/merged.bom.xml 

# ./cloud-mta-build-tool.exe sbom-gen -s /c/Workspace/MTA/cf-mta-examples/blue-green-deploy-strategy/hello-blue -b sbom-gen-result/merged.bom.xml