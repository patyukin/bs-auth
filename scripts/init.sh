#!/bin/bash

include ./../.env

set -e

psql -v ON_ERROR_STOP=1 --username "$PG_USER" <<-EOSQL
  CREATE DATABASE "$PG_BN_NAME";
EOSQL

exec "$@"
