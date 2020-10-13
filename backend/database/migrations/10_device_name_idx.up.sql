BEGIN;

-- See https://www.postgresql.org/docs/current/collation.html#COLLATION-NONDETERMINISTIC
CREATE COLLATION ndcoll (provider = icu, locale = 'und', deterministic = false);

ALTER TABLE "device" ALTER COLUMN "name" TYPE text COLLATE ndcoll;

CREATE UNIQUE INDEX ON "device" ("name", "id");

COMMIT;
