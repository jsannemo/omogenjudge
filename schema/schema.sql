-- Database schema for the OmogenJudge judging system.
-- NOTE: documentation of the database tables and fields can be found in the corresponding model structs in storage/models.

CREATE TABLE stored_file(
	file_hash VARCHAR(256) PRIMARY KEY,
	url bytea NOT NULL
);

GRANT ALL ON stored_file TO omogenjudge;

-- Extracts the URL of a file hash.
CREATE FUNCTION file_url(in VARCHAR(256), out bytea)
AS $$ SELECT url FROM stored_file WHERE file_hash = $1 $$
LANGUAGE SQL;

CREATE TABLE account(
	account_id SERIAL PRIMARY KEY,
	username TEXT NOT NULL,
	password_hash TEXT NOT NULL,
	full_name TEXT NOT NULL,
	email TEXT NOT NULL
);

CREATE UNIQUE INDEX account_username ON account(username);
CREATE UNIQUE INDEX account_email ON account(email);

GRANT ALL ON account TO omogenjudge;
GRANT ALL ON account_account_id_seq TO omogenjudge;

CREATE TABLE problem(
	problem_id SERIAL PRIMARY KEY,
	short_name TEXT NOT NULL,
	author TEXT NOT NULL,
	license TEXT NOT NULL,
	current_version INTEGER
);

CREATE UNIQUE INDEX problem_shortname ON problem(short_name);

GRANT ALL ON problem TO omogenjudge;
GRANT ALL ON problem_problem_id_seq TO omogenjudge;

CREATE TABLE problem_version(
	problem_version_id SERIAL PRIMARY KEY,
	problem_id INTEGER NOT NULL REFERENCES problem ON DELETE CASCADE,
	time_limit_ms INTEGER NOT NULL,
	memory_limit_kb INTEGER NOT NULL
);

GRANT ALL ON problem_version TO omogenjudge;
GRANT ALL ON problem_version_problem_version_id_seq TO omogenjudge;

-- NOTE: this is added here to resolve the cyclic dependence between the two tables.
-- The DEFERRABLE INITIALLY IMMEDIATE clause is used to resolve the corresponding cyclic dependence between the removal
-- of a problem_version row and its praent problem row, defering the check to the end of a transaction.
ALTER TABLE problem ADD FOREIGN KEY (current_version) REFERENCES problem_version(problem_version_id) DEFERRABLE INITIALLY IMMEDIATE;

CREATE TABLE problem_output_validator(
	problem_version_id INTEGER NOT NULL REFERENCES problem_version ON DELETE CASCADE,
	validator_language_id TEXT NOT NULL,
	validator_source_zip VARCHAR(256) NOT NULL,
	PRIMARY KEY(problem_version_id)
);

GRANT ALL ON problem_output_validator TO omogenjudge;

CREATE TABLE problem_statement(
	problem_id INTEGER NOT NULL REFERENCES problem ON DELETE CASCADE,
	language TEXT NOT NULL,
	title TEXT NOT NULL,
	html TEXT NOT NULL,
	PRIMARY KEY(problem_id, language)
);

GRANT ALL ON problem_statement TO omogenjudge;

CREATE TABLE problem_statement_file(
	problem_id INTEGER NOT NULL,
	language TEXT NOT NULL,
	FOREIGN KEY(problem_id, language) REFERENCES problem_statement ON DELETE CASCADE,
	file_path TEXT NOT NULL,
	file_hash VARCHAR(256) NOT NULL REFERENCES stored_file,
	PRIMARY KEY(problem_id, language, file_path)
);

CREATE TABLE problem_testgroup(
	problem_testgroup_id SERIAL PRIMARY KEY,
	problem_version_id INTEGER NOT NULL REFERENCES problem_version ON DELETE CASCADE,
	testgroup_name TEXT NOT NULL,
	public_visibility BOOLEAN NOT NULL,
	score INTEGER NOT NULL,
	output_validator_flags TEXT NOT NULL
);

CREATE INDEX problem_testgroup_problem_id ON problem_testgroup(problem_version_id);

GRANT ALL ON problem_testgroup TO omogenjudge;
GRANT ALL ON problem_testgroup_problem_testgroup_id_seq TO omogenjudge;

CREATE TABLE problem_testcase(
	problem_testcase_id SERIAL PRIMARY KEY,
	problem_testgroup_id INTEGER NOT NULL REFERENCES problem_testgroup ON DELETE CASCADE,
	testcase_name TEXT NOT NULL,
	input_file_hash VARCHAR(256) NOT NULL REFERENCES stored_file,
	output_file_hash VARCHAR(256) NOT NULL REFERENCES stored_file
);

GRANT ALL ON problem_testcase TO omogenjudge;
GRANT ALL ON problem_testcase_problem_testcase_id_seq TO omogenjudge;

CREATE INDEX problem_testcase_problem_testgroup_id ON problem_testcase(problem_testgroup_id);

CREATE TYPE status AS ENUM('new', 'compiling', 'compilation_failed', 'running', 'successful');
CREATE TYPE verdict AS ENUM('unjudged', 'accepted', 'wrong_answer', 'time_limit_exceeded', 'run_time_error');

CREATE TABLE submission(
	submission_id SERIAL PRIMARY KEY,
	account_id INTEGER NOT NULL REFERENCES account,
	problem_id INTEGER NOT NULL REFERENCES problem,
	language TEXT NOT NULL,
	current_run INTEGER,
	date_created TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT current_timestamp
);

GRANT ALL ON submission TO omogenjudge;
GRANT ALL ON submission_submission_id_seq TO omogenjudge;

CREATE TABLE submission_file(
	submission_id INTEGER NOT NULL REFERENCES submission ON DELETE CASCADE,
	file_path TEXT NOT NULL,
	file_contents TEXT NOT NULL
);

GRANT ALL ON submission_file TO omogenjudge;

CREATE INDEX submission_file_submission ON submission_file(submission_id);

CREATE TABLE submission_run(
	submission_run_id SERIAL PRIMARY KEY,
    submission_id INTEGER NOT NULL REFERENCES submission ON DELETE CASCADE,
	problem_version_id INTEGER NOT NULL REFERENCES problem_version ON DELETE SET NULL,
	date_created TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT current_timestamp,
	status status NOT NULL,
    time_usage_ms INTEGER NOT NULL,
    score INTEGER NOT NULL,
	verdict verdict NOT NULL,
	compile_error TEXT
);

ALTER TABLE submission ADD FOREIGN KEY (current_run) REFERENCES submission_run(submission_run_id) DEFERRABLE INITIALLY IMMEDIATE;

CREATE FUNCTION notify_run() RETURNS TRIGGER AS $$
BEGIN
    PERFORM pg_notify('new_run', (NEW.submission_run_id)::text);
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER "new_run"
    AFTER INSERT ON submission_run
    FOR EACH ROW EXECUTE PROCEDURE notify_run();

CREATE INDEX submission_run_status ON submission_run(status);
CREATE INDEX submission_run_submission ON submission_run(submission_id);

GRANT ALL ON submission_run TO omogenjudge;
GRANT ALL ON submission_run_submission_run_id_seq TO omogenjudge;

CREATE TABLE submission_group_run(
	submission_run_id INTEGER NOT NULL REFERENCES submission_run ON DELETE CASCADE,
	problem_testgroup_id INTEGER NOT NULL REFERENCES problem_testgroup,
	date_created TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT current_timestamp,
    time_usage_ms INTEGER NOT NULL,
    score INTEGER NOT NULL,
	verdict verdict NOT NULL
);

GRANT ALL ON submission_group_run TO omogenjudge;

CREATE TABLE submission_case_run(
	submission_run_id INTEGER NOT NULL REFERENCES submission_run ON DELETE CASCADE,
	problem_testcase_id INTEGER NOT NULL REFERENCES problem_testcase,
	date_created TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT current_timestamp,
    time_usage_ms INTEGER NOT NULL,
    score INTEGER NOT NULL,
	verdict verdict NOT NULL
);

GRANT ALL ON submission_case_run TO omogenjudge;

CREATE TABLE contest(
	contest_id SERIAL PRIMARY KEY,
	short_name TEXT NOT NULL,
	host_name TEXT,
	start_time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT current_timestamp,
	duration INTERVAL NOT NULL,
	title TEXT NOT NULL,
	hidden_scoreboard BOOLEAN NOT NULL
);

CREATE INDEX contest_host_name ON contest(host_name);
CREATE UNIQUE INDEX contest_short_name ON contest(short_name);

GRANT ALL ON contest TO omogenjudge;
GRANT ALL ON contest_contest_id_seq TO omogenjudge;

CREATE TABLE team(
	team_id SERIAL PRIMARY KEY,
	contest_id INTEGER NOT NULL REFERENCES contest ON DELETE CASCADE,
	team_name TEXT,
	virtual BOOLEAN NOT NULL,
	unofficial BOOLEAN NOT NULL,
	start_time TIMESTAMP WITH TIME ZONE,
	team_data json NOT NULL
);

CREATE INDEX team_contest_id ON team(contest_id);

GRANT ALL ON team TO omogenjudge;
GRANT ALL ON team_team_id_seq TO omogenjudge;

CREATE TABLE team_member(
	team_id INTEGER NOT NULL REFERENCES team ON DELETE CASCADE,
	account_id INTEGER NOT NULL REFERENCES account,
	PRIMARY KEY(team_id, account_id)
);

GRANT ALL ON team_member TO omogenjudge;

CREATE TABLE contest_problem(
	contest_id INTEGER NOT NULL REFERENCES contest ON DELETE CASCADE,
	problem_id INTEGER NOT NULL REFERENCES problem,
	label TEXT NOT NULL,
	PRIMARY KEY (contest_id, problem_id)
);

GRANT ALL ON contest_problem TO omogenjudge;
