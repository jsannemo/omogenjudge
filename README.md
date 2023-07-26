# Omogen Judging System The judging system is only tested on Ubuntu LTS (22.04).
Development or production use on any other system may not work.

## Setup
To build the system from source rather than downloading the ready-built packages:
1. Install the build dependencies with `./development/install_builddeps.sh`.
1. Run `./packaging/build-host.sh`, `./packaging/build-queue.sh` and `./packaging/build-web.sh`.

To install The system:

1. Following the setup in the [sandbox README](https://github.com/jsannemo/omogenexec) on all machines that should judge submissions.
1. Install the built `omogenjudge-web.deb` and `omogenjudge-queue.deb` packages on the machine hosting the frontend.
1. Install the built `omogenjudge-host.deb` package on all machines that should judge submissions. 

## Configuration
All configuration lives in `/etc/omogen/`.

## Frontend Development Setup
First, follow the setup and configuration sections to set up the backend, except that you shouldn't install `omogenjudge-web`.

First, you need to set up some development tooling:

1. Ensure that `$HOME/.local/bin` is in your PATH (for example by adding `PATH=$PATH:HOME/.local/bin` at the end of your `~/.bashrc`, if using bash). `poetry` is installed there later on.
1. Install the development dependencies by running `./development/install_devdeps.sh`
1. Create your local Django configuration by copying `omogenjudge/settings/local_development.example.py` to `omogenjudge/settings/local_development.py`
1. Setup a new database by running `./development/new_db.sh`.

To start the frontend server, you need to perform two steps:

1. Start the frontend asset compiler to build CSS and JavaScript upon changes by running `./frontend_assets/watch_assets.sh`.
1. Start the webserver by running `poetry run python manage.py runserver`.

After pulling in new changes, you might need to do two things:
1. Re-run `./development/install_devdeps.sh` since dependencies might have been updated.
1. Run `poetry run python manage.py migrate` to apply any database schema changes.

## Backend Development Setup
First, follow the setup and configuration sections, except that you don't need to install `omogenjudge-web` if you have your local development installation.

1. The evaluator library is included during compilation by the judgehost. If you need to make changes to it, check out the `omogenexec` repository as a sibling to the `omogenjudge` repository.  Update the `WORKSPACE` file to point to your local `omogenexec` copy instead (search for `EVALUATOR LIB` to find the right place) and follow the next point to run your own judgehost build.
1. Kill the auto-started judgehost with `sudo systemctl stop omogenjudge-host`. Enter the `judgehost` directory and run the annoying command `bazel build //judgehost:omogenjudge-host && sudo cp ./bazel-bin/judgehost/omogenjudge-host_/omogenjudge-host . && sudo -u omogenjudge-host omogenjudge-host` to start the judgehost.
1. Kill the auto-started one with `sudo systemctl stop omogenjudge-queue`. Enter the `judgehost` directory and run `bazel run //queue:omogenjudge-queue`.

