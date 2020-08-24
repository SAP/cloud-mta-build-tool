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
2. Run `go mod vendor` in the repository root.
3. If you’ve fixed a bug or added code that should be tested, add tests!
4. See commit prefix section
5. Ensure the test suite passes via `go test -v ./... ` Tip: you can use command `make tests`.
6. You can test the binary by using command `make` which will build the binary for each target OS.
7. If you change some config file you should run `go generate` command, this will create equivalent byte content file. 
8. Format your code with `go fmt` and run [linter](https://github.com/golang/lint) or better use `make tools` `make lint` on your changes.


## Contribution Prerequisites

1. You have [go](https://golang.org/dl/) installed at v1.13+
2. You are familiar with [GIT](https://git-scm.com/) 

## Commit Prefix

- [feat] (new feature for the user, not a new feature for build script)
- [fix] (bug fix for the user, not a fix to a build script)
- [docs] (changes to the documentation)
- [style] (formatting, missing semi colons, etc; no production code change)
- [refactor] (refactoring production code, eg. renaming a variable)
- [test] (adding missing tests, refactoring tests; no production code change)
- [chore] (updating grunt tasks etc; no production code change)
