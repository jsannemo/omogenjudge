# Omogen Judging System

## Frontend development Setup
First, you need to install some dependencies:
- [NPM](https://github.com/nodesource/distributions/blob/master/README.md)
- [Poetry](https://python-poetry.org/docs/)
- Bazelisk (`npm install -g @bazel/bazelisk`)
- PostgreSQL (`sudo apt install postgresql`)

Then, some configuration:
- Update `omogenjudge/settings/local_development.py` with your database port (if different from `5432`)

Now, the database must be setup:
- `./admin/new_db.sh`
- `poetry install`
- `poetry run python manage.py migrate`

You're ready to go!
- `poetry run python manage.py runserver
