BEGIN;

CREATE TABLE student_card
(
    id_hash         bytea NOT NULL,
    organization_id uuid  NOT NULL,
    student_id      uuid  NOT NULL,

    PRIMARY KEY (id_hash, organization_id),
    FOREIGN KEY (organization_id) REFERENCES organization (id) ON DELETE CASCADE,
    FOREIGN KEY (student_id) REFERENCES student (id) ON DELETE CASCADE
);

DROP TABLE student_token;

COMMIT;
