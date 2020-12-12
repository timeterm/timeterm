BEGIN;

ALTER TABLE device
    DROP COLUMN status;

DROP TYPE device_status;

COMMIT;
