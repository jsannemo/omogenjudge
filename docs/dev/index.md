# Developing OmogenJudge
OmogenJudge is mainly built to work on Linux systems that follow certain conventions.
So far, it has only been tested on recent Ubuntu and Debian distributions, with only Ubuntu LTS versions
officially supported for running production-wise.

If you don't have such a system, you should probably develop OmogenJudge in a virtual machine instead.

## Building
OmogenJudge is built using the [Bazel](http://bazel.build) build system.
A recent version of Bazel needs to be installed to run OmogenJudge.

Once installed, the Debian packages for installation can be built using the script `scripts/build/packages.sh`.
This command builds development versions of all packages.

One can also use `scripts/build/all.sh` to build all targets in the repository to verify there are no broken targets.
This is also a pre-commit check.

## Running
To run development versions of everything, you can use the `scripts/deploy/{all,master,local,frontend,sandbox}.sh` scripts depending on what you want to deploy.
Note that packages depend on each other in that the frontend depends on the master, which depends on local, which depends on sandbox.

Once development versions are installed, be sure to follow the [production install guide](/docs/prod/install.md) to set up your system properly (like enabling quota).

## Making commits
Commits should never break the build -- this is also verified by a Travis CI pre-commit hook.

## Adding dependencies
Dependencies need to be added explicitly to the `WORKSPACE` root file.
For Golang dependencies, this includes indirect dependencies as well.
When adding dependencies, make sure to use checksums and point to exact commits (not tags) to keep the build entirely reproducible.

## Style guide
For any language for which there is a Google style guide, we aim to follow it.
