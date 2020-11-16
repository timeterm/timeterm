BEGIN;

ALTER TABLE device_token
    ALTER COLUMN created_at DROP DEFAULT;

ALTER TABLE device_registration_token
    ALTER COLUMN created_at DROP DEFAULT;

COMMIT;
