#!/usr/bin/env bash

psql omogenjudge -c "DELETE FROM problem WHERE short_name = '$1'";
