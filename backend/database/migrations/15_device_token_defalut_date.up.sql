BEGIN;

ALTER TABLE device_token
    ALTER created_at SET DEFAULT now();

COMMIT;
