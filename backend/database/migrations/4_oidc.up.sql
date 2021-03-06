BEGIN;

CREATE TABLE "oauth2_state"
(
    "state"        uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    "issuer"       text        NOT NULL,
    "redirect_url" text        NOT NULL,
    "created_at"   timestamptz NOT NULL DEFAULT now(),
    "expires_at"   timestamptz NOT NULL DEFAULT now() + '30 minutes'
);

CREATE TABLE "oidc_federation"
(
    "oidc_subject"  text NOT NULL,
    "oidc_issuer"   text NOT NULL,
    "oidc_audience" text NOT NULL,
    "user_id"       uuid NOT NULL,

    PRIMARY KEY ("oidc_subject", "oidc_issuer"),
    FOREIGN KEY ("user_id") REFERENCES "user" ("id") ON DELETE CASCADE
);

ALTER TABLE "user"
    DROP COLUMN "keycloak_subject";

COMMIT;