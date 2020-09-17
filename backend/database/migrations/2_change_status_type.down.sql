BEGIN;

ALTER TABLE "device"
    ALTER COLUMN "status" TYPE text;

DROP TYPE "device_status";

COMMIT;