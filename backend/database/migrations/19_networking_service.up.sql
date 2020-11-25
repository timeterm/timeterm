BEGIN;

CREATE TABLE networking_service
(
    "id"              uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    "organization_id" uuid NOT NULL,
    "name"            text NOT NULL,

    FOREIGN KEY (organization_id) REFERENCES organization (id) ON DELETE CASCADE
);

COMMIT;