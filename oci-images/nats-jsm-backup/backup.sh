#!/usr/bin/env sh

set -ex

mc alias set dest $S3_ENDPOINT $S3_ACCESS_KEY_ID $S3_SECRET_ACCESS_KEY --api ${S3_API_SIGNATURE:-S3v4}

current_date=`date -u +"%Y-%m-%dT%H:%M:%SZ"`

nats backup --data nats-jsm-backup
tar cf - nats-jsm-backup | zstd | mc pipe dest/timeterm-nats-jsm-backup/$current_date-all.tar.zst


