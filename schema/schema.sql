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
  short_name TEXT NOT NULL,
  author TEXT NOT NULL,
  license TEXT NOT NULL,
  time_limit_ms INTEGER NOT NULL,
  memory_limit_kb INTEGER NOT NULL
);

CREATE UNIQUE INDEX problem_shortname ON problem(short_name);

GRANT ALL ON problem TO omogenjudge;
GRANT ALL ON problem_problem_id_seq TO omogenjudge;

CREATE TABLE problem_output_validator(
  problem_id INTEGER NOT NULL REFERENCES problem ON DELETE CASCADE,
  validator_language_id TEXT NOT NULL,
  validator_source_zip hash NOT NULL,
  UNIQUE(problem_id)
);

GRANT ALL ON problem_output_validator TO omogenjudge;

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
  problem_id INTEGER NOT NULL REFERENCES problem ON DELETE CASCADE,
  testgroup_name TEXT NOT NULL,
  public_visibility BOOLEAN NOT NULL
);

CREATE INDEX problem_testgroup_problem_id ON problem_testgroup(problem_id);

GRANT ALL ON problem_testgroup TO omogenjudge;
GRANT ALL ON problem_testgroup_problem_testgroup_id_seq TO omogenjudge;

CREATE TABLE problem_testcase(
  problem_testcase_id SERIAL PRIMARY KEY,
  problem_testgroup_id INTEGER NOT NULL REFERENCES problem_testgroup ON DELETE CASCADE,
  testcase_name TEXT NOT NULL,
  input_file_hash hash NOT NULL REFERENCES stored_file,
  output_file_hash hash NOT NULL REFERENCES stored_file
);

GRANT ALL ON problem_testcase TO omogenjudge;
GRANT ALL ON problem_testcase_problem_testcase_id_seq TO omogenjudge;

CREATE INDEX problem_testcase_problem_testgroup_id ON problem_testcase(problem_testgroup_id);

CREATE TYPE status AS ENUM('new', 'compiling', 'compilation_failed', 'running', 'successful');
CREATE TYPE verdict AS ENUM('VERDICT_UNSPECIFIED', 'UNJUDGED', 'ACCEPTED', 'WRONG_ANSWER', 'TIME_LIMIT_EXCEEDED', 'RUN_TIME_ERROR');

-- Submission tables
CREATE TABLE submission(
  submission_id SERIAL PRIMARY KEY,
  account_id INTEGER NOT NULL REFERENCES account,
  problem_id INTEGER NOT NULL REFERENCES problem,
  language TEXT NOT NULL,
  date_created TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT current_timestamp,
  status status NOT NULL DEFAULT 'new',
  verdict verdict DEFAULT 'UNJUDGED',
  compile_error TEXT
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
  submission_id INTEGER NOT NULL REFERENCES submission ON DELETE CASCADE,
  file_path TEXT NOT NULL,
  file_contents TEXT NOT NULL
);

GRANT ALL ON submission_file TO omogenjudge;

CREATE INDEX submission_file_submission ON submission_file(submission_id);

-- Course tables
CREATE TABLE course(
  course_id SERIAL PRIMARY KEY,
  course_short_name TEXT NOT NULL,
  UNIQUE(course_short_name)
);

GRANT ALL ON course TO omogenjudge;
GRANT ALL ON course_course_id_seq TO omogenjudge;

CREATE TABLE course_localization(
  course_id INTEGER NOT NULL REFERENCES course ON DELETE CASCADE,
  course_language TEXT NOT NULL,
  course_name TEXT NOT NULL,
  course_summary TEXT NOT NULL,
  course_description TEXT NOT NULL,
  UNIQUE(course_id, course_language)
);

GRANT ALL ON course_localization TO omogenjudge;

CREATE TABLE course_chapter(
  course_id INTEGER NOT NULL REFERENCES course ON DELETE CASCADE,
  course_chapter_id SERIAL PRIMARY KEY,
  chapter_short_name TEXT NOT NULL,
  UNIQUE(course_id, course_chapter_id)
);

GRANT ALL ON course_chapter TO omogenjudge;
GRANT ALL ON course_chapter_course_chapter_id_seq TO omogenjudge;

CREATE TABLE course_chapter_localization(
  course_chapter_id INTEGER NOT NULL REFERENCES course_chapter ON DELETE CASCADE,
  chapter_language TEXT NOT NULL,
  chapter_name TEXT NOT NULL,
  chapter_summary TEXT NOT NULL,
  chapter_description TEXT NOT NULL,
  UNIQUE(course_chapter_id, chapter_language)
);

GRANT ALL ON course_chapter_localization TO omogenjudge;

CREATE TABLE course_section(
  course_chapter_id INTEGER NOT NULL REFERENCES course_chapter ON DELETE CASCADE,
  course_section_id SERIAL PRIMARY KEY,
  section_short_name TEXT NOT NULL,
  UNIQUE(course_chapter_id, section_short_name)
);

GRANT ALL ON course_section TO omogenjudge;
GRANT ALL ON course_section_course_section_id_seq TO omogenjudge;

CREATE TABLE course_section_localization(
  course_section_id INTEGER NOT NULL REFERENCES course_section ON DELETE CASCADE,
  section_language TEXT NOT NULL,
  section_name TEXT NOT NULL,
  section_summary TEXT NOT NULL,
  section_contents TEXT NOT NULL,
  UNIQUE(course_section_id, section_language)
);

GRANT ALL ON course_section_localization TO omogenjudge;

CREATE TABLE editor_file(
  editor_file_id SERIAL PRIMARY KEY,
  account_id INTEGER NOT NULL REFERENCES account_id,
  file_name TEXT NOT NULL,
  file_content TEXT NOT NULL
);

GRANT ALL ON editor_file TO omogenjudge;

CREATE UNIQUE INDEX editor_file_account_name ON editor_file(account_id, file_name);
