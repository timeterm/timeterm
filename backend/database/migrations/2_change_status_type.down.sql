BEGIN;

ALTER TABLE "device"
    ALTER COLUMN "status" TYPE text USING "status"::text;

DROP TYPE "device_status";

COMMIT;