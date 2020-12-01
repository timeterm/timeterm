BEGIN;

DROP TABLE student_card;

CREATE TABLE student_token
(
    id                uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
    student_id        uuid        NOT NULL,
    expiration        timestamptz NOT NULL,
    card_id_hash      bytea       NOT NULL,
    card_id_hash_salt bytea       NOT NULL,
    token             bytea       NOT NULL,

    FOREIGN KEY (student_id) REFERENCES student (id) ON DELETE CASCADE
);

COMMIT;
