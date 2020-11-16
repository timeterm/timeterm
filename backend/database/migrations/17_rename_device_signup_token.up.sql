BEGIN;

ALTER TABLE device_signup_token
    RENAME TO device_registration_token;

COMMIT;
