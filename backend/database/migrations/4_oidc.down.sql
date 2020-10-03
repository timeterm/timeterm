BEGIN;

DROP TABLE "oauth2_state";

DROP TABLE "oidc_federation";

ALTER TABLE "user"
    ADD COLUMN "keycloak_subject" text NOT NULL DEFAULT '';

COMMIT;