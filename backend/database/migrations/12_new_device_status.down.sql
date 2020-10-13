BEGIN;

ALTER TYPE "device_status" RENAME TO "device_status_";

CREATE TYPE "device_status" AS ENUM ('online', 'offline');

ALTER TABLE "device"
    ALTER COLUMN "status" TYPE "device_status" USING 'offline';

ALTER TABLE "device"
    ALTER COLUMN "status" SET NOT NULL;

DROP TYPE "device_status_";

COMMIT;
