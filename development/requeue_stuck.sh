#!/usr/bin/env bash
set -e

sudo service omogenjudge-queue stop || true

sudo -u postgres psql omogenjudge -c "DELETE FROM submission_group_run WHERE submission_run_id IN (SELECT submission_run_id FROM submission_run WHERE status != 'done'); UPDATE submission_run SET status = 'queued' WHERE status = 'running' OR status = 'compiling'"

sudo service omogenjudge-queue start || true
