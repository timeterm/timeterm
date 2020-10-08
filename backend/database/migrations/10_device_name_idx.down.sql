BEGIN;

DROP COLLATION ndcoll;

ALTER TABLE "device" ALTER COLUMN "name" TYPE text;

DROP INDEX "device_name_id_idx";

COMMIT;
