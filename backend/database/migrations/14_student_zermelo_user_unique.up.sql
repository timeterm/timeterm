BEGIN;

ALTER TABLE "student" ADD UNIQUE ("organization_id", "zermelo_user");

COMMIT;
