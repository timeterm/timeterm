BEGIN;

ALTER TABLE "user"
    ADD UNIQUE ("email");

COMMIT;
