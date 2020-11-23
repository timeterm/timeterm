BEGIN;

CREATE TABLE operator
(
    subject text PRIMARY KEY,
    name    text NOT NULL UNIQUE
);

CREATE TABLE account
(
    subject          text PRIMARY KEY,
    name             text NOT NULL UNIQUE,
    operator_subject text NOT NULL,

    FOREIGN KEY (operator_subject) REFERENCES operator (subject) ON DELETE RESTRICT
);

CREATE TABLE "user"
(
    subject         text PRIMARY KEY,
    name            text NOT NULL UNIQUE,
    account_subject text NOT NULL,

    FOREIGN KEY (account_subject) REFERENCES account (subject)
);

COMMIT;
