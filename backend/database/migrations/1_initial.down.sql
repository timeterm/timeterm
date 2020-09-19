BEGIN;

-- We're not dropping the uuid-ossp extension because we can't be sure that it was created by us.

DROP TABLE "device";

DROP TABLE "user";

DROP TABLE "student_token";

DROP TABLE "student";

DROP TABLE "organization";

COMMIT;