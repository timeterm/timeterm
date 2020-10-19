#!/usr/bin/env sh

set -ex

mc alias set dest $S3_ENDPOINT $S3_ACCESS_KEY_ID $S3_SECRET_ACCESS_KEY --api ${S3_API_SIGNATURE:-S3v4}

current_date=`date -u +"%Y-%m-%dT%H:%M:%SZ"`

pg_dumpall --globals-only | zstd | mc pipe dest/timeterm-postgres-backup/$current_date-globals.sql.zst

psql \
	-X \
	-c "SELECT datname FROM pg_database WHERE datistemplate = false AND datname != 'postgres'" \
	--single-transaction \
	--set AUTOCOMMIT=off \
	--set ON_ERROR_STOP=on \
	--no-align \
	-t \
	--field-separator ' ' \
	--quiet \
| while read datname; do
	pg_dump -Fc $datname | zstd | mc pipe dest/timeterm-postgres-backup/$current_date-$datname.dump
done

