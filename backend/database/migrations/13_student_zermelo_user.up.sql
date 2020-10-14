BEGIN;

ALTER TABLE "student"
    RENAME COLUMN "zermelo_code" TO "zermelo_user";

COMMIT;
