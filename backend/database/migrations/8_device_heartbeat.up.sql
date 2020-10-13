BEGIN;

ALTER TABLE "device"
    ADD COLUMN "last_heartbeat" timestamptz;

COMMIT;
