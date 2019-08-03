#!/usr/bin/env bash

psql omogenjudge -c "DELETE FROM course WHERE course_short_name = '$1'";
