BEGIN;

ALTER TABLE "device"
    DROP COLUMN "last_heartbeat";

COMMIT;
