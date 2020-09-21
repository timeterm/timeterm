BEGIN;

ALTER TABLE "organization"
    DROP COLUMN "zermelo_institution";

COMMIT;