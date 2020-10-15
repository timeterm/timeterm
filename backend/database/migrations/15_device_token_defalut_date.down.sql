BEGIN;

ALTER TABLE device_token
    ALTER created_at DROP DEFAULT;

COMMIT;
