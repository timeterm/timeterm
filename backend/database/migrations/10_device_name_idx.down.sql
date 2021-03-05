BEGIN;

ALTER TABLE "device" ALTER COLUMN "name" TYPE text;

DROP COLLATION ndcoll;

DROP INDEX "device_name_id_idx";

COMMIT;
