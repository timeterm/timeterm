BEGIN;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE "organization" (
    "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    "name" text NOT NULL
);

CREATE TABLE "student" (
    "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    "organization_id" uuid NOT NULL,

    FOREIGN KEY (organization_id) REFERENCES organization (id) ON DELETE CASCADE
);

CREATE TABLE "student_token" (
    "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    "student_id" uuid NOT NULL,
    "expiration" timestamptz NOT NULL,
    "card_id_hash" bytea NOT NULL,
    "card_id_hash_salt" bytea NOT NULL,
    "token" bytea NOT NULL,

    FOREIGN KEY (student_id) REFERENCES student (id) ON DELETE CASCADE
);

CREATE TABLE "user" (
    "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    "name" text NOT NULL,
    "organization_id" uuid NOT NULL,
    "keycloak_subject" text NOT NULL,

    FOREIGN KEY (organization_id) REFERENCES organization (id) ON DELETE CASCADE
);

CREATE TABLE "device" (
    "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    "organization_id" uuid NOT NULL,
    "name" text NOT NULL,
    "status" text NOT NULL,

    FOREIGN KEY (organization_id) REFERENCES organization (id) ON DELETE CASCADE
);

COMMIT;