BEGIN;

CREATE TABLE device_signup_token
(
    token_hash bytea PRIMARY KEY,
    organization_id  uuid,
    created_at timestamptz NOT NULL,
    expires_at timestamptz,

    FOREIGN KEY (organization_id) REFERENCES organization (id) ON DELETE CASCADE
);

COMMIT;
