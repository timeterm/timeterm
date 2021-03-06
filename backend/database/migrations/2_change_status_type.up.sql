BEGIN;

CREATE TYPE "device_status" AS ENUM ('online', 'offline');
ALTER TABLE "device"
    ALTER COLUMN "status" TYPE "device_status" USING "status"::"device_status";

COMMIT;