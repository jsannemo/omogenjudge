-- TODO: document database fields

-- File tables
CREATE TYPE hash AS (hash VARCHAR(256));

CREATE TABLE stored_file(
  file_hash hash PRIMARY KEY,
  url bytea NOT NULL
);

GRANT ALL ON stored_file TO omogenjudge;

CREATE FUNCTION file_url(in hash, out bytea)
   AS $$ SELECT url FROM stored_file WHERE file_hash = $1 $$
   LANGUAGE SQL;

CREATE FUNCTION file_hash(in hash, out VARCHAR(256))
   AS $$ SELECT $1.hash $$
   LANGUAGE SQL;

-- Account tables
CREATE TABLE account(
  account_id SERIAL PRIMARY KEY,
  username TEXT NOT NULL,
  password_hash TEXT NOT NULL
);

CREATE UNIQUE INDEX account_username ON account(username);

GRANT ALL ON account TO omogenjudge;
GRANT ALL ON account_account_id_seq TO omogenjudge;

-- Problem tables
CREATE TABLE problem(
  problem_id SERIAL PRIMARY KEY,
  short_name TEXT NOT NULL
);

CREATE UNIQUE INDEX problem_shortname ON problem(short_name);

GRANT ALL ON problem TO omogenjudge;
GRANT ALL ON problem_problem_id_seq TO omogenjudge;

CREATE TABLE problem_statement(
  problem_id INTEGER,
  language TEXT NOT NULL,
  title TEXT NOT NULL,
  html TEXT NOT NULL,
  PRIMARY KEY(problem_id, language)
);

GRANT ALL ON problem_statement TO omogenjudge;

CREATE TABLE problem_testgroup(
  problem_testgroup_id SERIAL PRIMARY KEY,
  problem_id INTEGER NOT NULL REFERENCES problem,
  testgroup_name TEXT NOT NULL,
  public_visibility BOOLEAN NOT NULL
);

CREATE INDEX problem_testgroup_problem_id ON problem_testgroup(problem_id);

GRANT ALL ON problem_testgroup TO omogenjudge;
GRANT ALL ON problem_testgroup_problem_testgroup_id_seq TO omogenjudge;

CREATE TABLE problem_testcase(
  problem_testcase_id SERIAL PRIMARY KEY,
  problem_testgroup_id INTEGER NOT NULL REFERENCES problem_testgroup,
  testcase_name TEXT NOT NULL,
  input_file_hash hash NOT NULL REFERENCES stored_file,
  output_file_hash hash NOT NULL REFERENCES stored_file
);

GRANT ALL ON problem_testcase TO omogenjudge;
GRANT ALL ON problem_testcase_problem_testcase_id_seq TO omogenjudge;

CREATE INDEX problem_testcase_problem_testgroup_id ON problem_testcase(problem_testgroup_id);

CREATE TYPE status AS ENUM('new', 'compiling', 'running', 'successful');
CREATE TYPE verdict AS ENUM('AC', 'WA', 'TLE', 'RTE');

-- Submission tables
CREATE TABLE submission(
  submission_id SERIAL PRIMARY KEY,
  account_id INTEGER NOT NULL REFERENCES account,
  problem_id INTEGER NOT NULL REFERENCES problem,
  status status NOT NULL
);

GRANT ALL ON submission TO omogenjudge;
GRANT ALL ON submission_submission_id_seq TO omogenjudge;

CREATE INDEX submission_status ON submission(status);

CREATE FUNCTION notify_submission() RETURNS TRIGGER AS $$
BEGIN
  PERFORM pg_notify('new_submission', (NEW.submission_id)::text);
  RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER "new_submission"
AFTER INSERT ON submission
FOR EACH ROW EXECUTE PROCEDURE notify_submission();

CREATE TABLE submission_file(
  submission_id INTEGER NOT NULL REFERENCES submission,
  file_path TEXT NOT NULL,
  file_contents bytea NOT NULL
);

GRANT ALL ON submission_file TO omogenjudge;

CREATE INDEX submission_file_submission ON submission_file(submission_id);
