dropdb omogenjudge || true
dropuser omogenjudge || true
psql -c "CREATE USER omogenjudge WITH PASSWORD 'omogenjudge';"
createdb omogenjudge
psql omogenjudge -f schema/schema.sql
