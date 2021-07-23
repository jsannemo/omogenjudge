# Omogen Judging System

## Frontend development Setup
First, you need to install some dependencies:
- [NPM](https://github.com/nodesource/distributions/blob/master/README.md)
- [Poetry](https://python-poetry.org/docs/)
- Bazelisk (`npm install -g @bazel/bazelisk`)
- PostgreSQL (`sudo apt install postgresql`)

Then, some configuration:
- Copy `omogenhost/webapi/dev/webapi.toml' to /etc/omogenjudge/webapi.toml`
- Update it with your database port (if different from `5432`)
- Update `omogenjudge/settings/local_development.py` with your database port (if different from `5432`)

Now, the database must be setup:
- `./admin/new_db.sh`
- `poetry install`
- `poetry run python manage.py migrate`

You're ready to go!
- `ibazel run //:serve` to start the frontend
- `ibazel run //webapi:webapi` to run the frontend API:s
