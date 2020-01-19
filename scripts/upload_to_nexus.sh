MBT_VERSION=$(cat ./VERSION)
echo $MBT_VERSION
curl -L https://github.com/SAP/cloud-mta-build-tool/releases/download/v${MBT_VERSION}/cloud-mta-build-tool_${MBT_VERSION}_Darwin_amd64.tar.gz -o cloud-mta-build-tool_${MBT_VERSION}_Darwin_amd64.tar.gz
curl -L https://github.com/SAP/cloud-mta-build-tool/releases/download/v${MBT_VERSION}/cloud-mta-build-tool_${MBT_VERSION}_Linux_amd64.tar.gz -o cloud-mta-build-tool_${MBT_VERSION}_Linux_amd64.tar.gz
curl -L https://github.com/SAP/cloud-mta-build-tool/releases/download/v${MBT_VERSION}/cloud-mta-build-tool_${MBT_VERSION}_Windows_amd64.tar.gz -o cloud-mta-build-tool_${MBT_VERSION}_Windows_amd64.tar.gz
artifactdeployer pack --script-file configpack --package-file cloud-mta-build-tool-pack -D mbtVersion=$MBT_VERSION
artifactdeployer deploy --package-file cloud-mta-build-tool-pack --artifact-version $MBT_VERSION --repo-url http://nexusmil.wdf.sap.corp:8081/nexus/content/repositories/sap.milestones.manual-uploads.hosted --repo-user $MVN_REPO_USER --repo-passwd $MVN_REPO_PASSWD
