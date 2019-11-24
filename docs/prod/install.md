# Installation guide

## System requirements
The judge "officially" supports only the latest Ubuntu LTS release, but should work on most recent Ubuntu and Debian releases.
Other than that, we generally recommend two separate servers to be used -- one for running submissions, and one for serving the web frontend and running the judging coordinator.
This reduces run-time noise and allows for better security by physically separating the host running untrusted code and the hosts with e.g. database access.

## Prerequisites
Before installing the judge, you need to

- [enable quota](quota.md) on the system where submissions will run
- make sure you have the packages needed to install `.deb` packages
- set up Postgres on the non-submission host
- install Bazel
- check out the repository

## Installing the judge
To build the installable packages, run `scripts/build/prod-packages.sh`, which builds production version of all packages.

Then, install the packages `omogenjudge-sandbox.deb` and `omogenjudge-local.deb` on the machine that should run submissions (in that order).
Now the judging machine is set up and ready to run submissions.

The other host should first create a database and execute the schema from `schema/schema.sql` in it.
Then, the packages `omogenjudge-master.deb` and `omogenjudge-frontend.deb` should be installed (in that order).

## Configuring the judge
TODO(jsannemo): add this section once judge can be configured

## Starting the judge
To start the servers, run:
- `systemctl start omogenjudge-frontend.service` for the frontend
- `systemctl start omogenjudge-local.service` for the local judge
- `systemctl start omogenjudge-master.service` for the judging coordinator
