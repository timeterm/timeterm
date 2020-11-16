BEGIN;

ALTER TABLE device_token
    ALTER COLUMN created_at SET DEFAULT now();

ALTER TABLE device_registration_token
    ALTER COLUMN created_at SET DEFAULT now();

COMMIT;
