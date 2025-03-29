#!/usr/bin/env bash

set -e

if [ -z "db-test.sqlite3" ]; then
    rm db-test.sqlite3
fi

cd ..
make
./bin/webapp ./e2e/db-test.sqlite3

