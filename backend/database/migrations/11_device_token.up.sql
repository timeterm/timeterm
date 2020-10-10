BEGIN;

CREATE TABLE device_token
(
    token_hash bytea PRIMARY KEY,
    device_id  uuid,
    created_at timestamptz NOT NULL,
    expires_at timestamptz,

    FOREIGN KEY (device_id) REFERENCES device (id) ON DELETE CASCADE
);

COMMIT;
