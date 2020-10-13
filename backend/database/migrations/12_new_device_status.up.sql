BEGIN;

ALTER TYPE "device_status" RENAME TO "device_status_";

CREATE TYPE "device_status" AS ENUM ('not_activated', 'ok');

ALTER TABLE "device"
    ALTER COLUMN "status" TYPE "device_status" USING 'not_activated';

ALTER TABLE "device"
    ALTER COLUMN "status" DROP NOT NULL;

DROP TYPE "device_status_";

COMMIT;
