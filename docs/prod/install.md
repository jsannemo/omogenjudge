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
Some parts of the judge now needs to be configured:

- the addresses for the submission server, so the judging coordinator can find it
- the database details for the postgres database for the judging coordinator and the frontend server.

The first config is in `/etc/omogen/master/local.conf`, and the second in `/etc/omogen/master/db.conf` and `/etc/omogen/frontend/db.conf`.

TODO: once the servers use authentication, this step will include copying encryption keys between hosts.

## Starting the judge
Now, the judging services can be started using 
