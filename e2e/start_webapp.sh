#!/usr/bin/env bash

set -e

if [ -f "db-test.sqlite3" ]; then
    rm db-test.sqlite3
fi

if [ -f "logs-test.log" ]; then
    rm logs-test.log
fi

cd ..
make
./bin/webapp ./e2e/db-test.sqlite3

