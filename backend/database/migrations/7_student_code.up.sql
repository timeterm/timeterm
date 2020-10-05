BEGIN;

ALTER TABLE "student"
    ADD COLUMN "zermelo_code" text;

COMMIT;
