BEGIN;

CREATE TABLE jwt
(
    subject text PRIMARY KEY,
    jwt     text NOT NULL
);

CREATE TABLE operator
(
    name    text PRIMARY KEY,
    subject text NOT NULL,

    FOREIGN KEY (subject) REFERENCES jwt (subject) ON DELETE RESTRICT
);

CREATE TABLE account
(
    name          text NOT NULL,
    subject       text NOT NULL,
    operator_name text NOT NULL,

    PRIMARY KEY (name, operator_name),
    FOREIGN KEY (subject) REFERENCES jwt (subject) ON DELETE RESTRICT,
    FOREIGN KEY (operator_name) REFERENCES operator (name) ON DELETE RESTRICT
);

CREATE TABLE "user"
(
    name          text NOT NULL,
    subject       text NOT NULL,
    account_name  text NOT NULL,
    operator_name text NOT NULL,

    PRIMARY KEY (name, account_name, operator_name),
    FOREIGN KEY (subject) REFERENCES jwt (subject) ON DELETE RESTRICT,
    FOREIGN KEY (account_name, operator_name) REFERENCES account (name, operator_name) ON DELETE RESTRICT
);

COMMIT;
