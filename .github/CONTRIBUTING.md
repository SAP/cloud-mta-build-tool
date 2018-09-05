## Open Development

All work on MBT happens directly on GitHub. 
Both core team members and external contributors send pull requests which go through the same review process.


## Branch Organization

We will do our best to keep the master branch in a good shape, with tests passing at all times. 
But in order to move fast, we will make API changes that your application might not be compatible with. 
We recommend that you use the latest stable version of MBT.

If you send a pull request, please do it against the `master` branch. 
We maintain stable branches for major versions separately but we don’t accept pull requests to them directly. 
Instead, we cherry-pick non-breaking changes from master to the latest stable major version.


## Semantic Versioning

MBT follows semantic versioning. 
We release patch versions for bug-fixes, minor versions for new features, and major versions for any breaking changes. 
When we make breaking changes, we also introduce deprecation warnings in a minor version 
so that our users learn about the upcoming changes and migrate their code in advance.
Every significant change will be documented in the changelog file.


## Sending a Pull Request

The team is monitoring for pull requests. We will review your pull request and either merge it, 
request changes to it, or close it with an explanation. 


## Before submitting a pull request, please make sure the following is done:

1. Fork the repository and create your branch from master.
2. Run [dep](https://github.com/golang/dep) ensure in the repository root.
3. If you’ve fixed a bug or added code that should be tested, add tests!
4. Ensure the test suite passes via `go test -v ./... ` Tip: you can use command `make test`.
5. Format your code with `go fmt` and run [linter](https://github.com/golang/lint) on your changes.
6. If you haven’t already, complete the CLA.


## Contribution Prerequisites

1. You have [go](https://golang.org/dl/) installed at v1.11+
2. You have [dep](https://github.com/golang/dep) installed at v0.5.0+
3. You are familiar with [GIT](https://git-scm.com/) 