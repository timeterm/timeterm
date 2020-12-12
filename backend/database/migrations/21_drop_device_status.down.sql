BEGIN;

CREATE TYPE device_status AS ENUM ('not_activated', 'ok');

ALTER TABLE device
    ADD COLUMN status device_status;

COMMIT;
