BEGIN;

CREATE TYPE admin_message_severity AS ENUM('error', 'info');

CREATE TABLE admin_message (
    organization_id uuid                   NOT NULL,
    logged_at       timestamptz            NOT NULL,
    severity        admin_message_severity NOT NULL,
    verbosity       int                    NOT NULL,
    nonce           bytea                  NOT NULL,
    data            bytea                  NOT NULL,

    PRIMARY KEY (organization_id, logged_at),
    FOREIGN KEY (organization_id) REFERENCES organization (id) ON DELETE CASCADE
);

COMMIT;

