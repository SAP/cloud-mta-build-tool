# AGENTS.md

## Project Overview
The Cloud MTA Build Tool (mbt) is a CLI tool that builds deployment-ready Multi-Target Application (MTA) archives (.mtar) from MTA project artifacts or module build artifacts. Primary tech stack: Go (core CLI and build logic) and Node.js (small npm wrapper/installer scripts). Key build tooling is driven from the repository Makefile and Go modules (go.mod).

## Dev Environment Tips
- Required runtimes
  - Go >= 1.13 (go.mod specifies go 1.13)
  - Node.js (used for the npm wrapper/installer; package.json present)
- Package manager commands
  - Node/npm: npm install (install JS wrapper dependencies listed in package.json)
  - Go modules: go mod download (download Go dependencies listed in go.mod)
- How to install dependencies
