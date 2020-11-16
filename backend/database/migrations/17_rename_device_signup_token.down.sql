BEGIN;

ALTER TABLE device_registration_token
    RENAME TO device_signup_token;

COMMIT;
