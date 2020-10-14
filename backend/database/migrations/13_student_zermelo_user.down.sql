BEGIN;

ALTER TABLE "student"
    RENAME COLUMN "zermelo_user" TO "zermelo_code";

COMMIT;
