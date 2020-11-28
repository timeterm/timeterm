BEGIN;

CREATE TABLE jwt
(
    subject text PRIMARY KEY
);

CREATE TABLE operator
(
    subject text PRIMARY KEY,
    name    text NOT NULL UNIQUE,

    FOREIGN KEY (subject) REFERENCES jwt (subject)
);

CREATE TABLE account
(
    subject          text PRIMARY KEY,
    name             text NOT NULL,
    operator_subject text NOT NULL,

    -- name can only be used once for this operator
    UNIQUE (name, operator_subject),
    FOREIGN KEY (subject) REFERENCES jwt (subject),
    FOREIGN KEY (operator_subject) REFERENCES operator (subject) ON DELETE RESTRICT
);

CREATE TABLE "user"
(
    subject         text PRIMARY KEY,
    name            text NOT NULL,
    account_subject text NOT NULL,

    -- name can only be used once for this account
    UNIQUE (name, account_subject),
    FOREIGN KEY (subject) REFERENCES jwt (subject),
    FOREIGN KEY (account_subject) REFERENCES account (subject)
);

CREATE TABLE jwt_migration
(
    version int PRIMARY KEY UNIQUE
);

COMMIT;
