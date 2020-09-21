BEGIN;

ALTER TABLE "organization"
    ADD COLUMN "zermelo_institution" text NOT NULL;

COMMIT;