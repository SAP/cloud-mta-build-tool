For Software Bill of Materials (SBOM) file generation, some SBOM tools are required to be installed in your environment:

#### Install CycloneDX Gomod
[cyclonedx-gomod](https://github.com/CycloneDX/cyclonedx-gomod) creates CycloneDX Software Bill of Materials (SBOM) from Go modules.

Run the following command to install CycloneDX Gomod:

```yaml
curl -fsSLO --compressed https://github.com/CycloneDX/cyclonedx-gomod/releases/download/v1.4.1/cyclonedx-gomod_1.4.0_linux_amd64.tar.gz

tar -xzf cyclonedx-gomod_1.4.0_linux_amd64.tar.gz

chmod a+rx cyclonedx-gomod

mv cyclonedx-gomod /usr/local/bin/cyclonedx-gomod
```

#### Install the CycloneDX Node Module
[cyclonedx-node-module](https://github.com/eoftedal/cyclonedx-node-module) creates CycloneDX BOMs from Node.js projects.

Run the following command to install the CycloneDX Node Module:

```yaml
npm install -g cyclonedx-bom
```

#### Install the CycloneDX CLI
The [CycloneDX CLI tool](https://github.com/CycloneDX/cyclonedx-cli) currently supports BOM analysis, modification, diffing, merging, format conversion, signing, and verification.

Run the following command to install the CycloneDX CLI tool:

```yaml
curl -fsSLO --compressed https://github.com/CycloneDX/cyclonedx-cli/releases/download/v0.24.2/cyclonedx-linux-x64

chmod a+rx cyclonedx-linux-x64

mv cyclonedx-linux-x64 /usr/local/bin/cyclonedx
```


