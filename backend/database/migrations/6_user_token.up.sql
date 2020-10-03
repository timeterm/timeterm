BEGIN;

CREATE TABLE "user_token"
(
    "token_hash" bytea PRIMARY KEY,
    "user_id"    uuid        NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    "expires_at" timestamptz NOT NULL DEFAULT now() + '24 hours',

    FOREIGN KEY ("user_id") REFERENCES "user" ("id") ON DELETE CASCADE
);

COMMIT;