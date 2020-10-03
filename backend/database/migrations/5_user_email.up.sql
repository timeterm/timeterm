BEGIN;

ALTER TABLE "user"
    ADD COLUMN "email" text;

UPDATE "user"
SET "email" = ''
WHERE "email" IS NULL;

ALTER TABLE "user"
    ALTER COLUMN "email" SET NOT NULL;

COMMIT;