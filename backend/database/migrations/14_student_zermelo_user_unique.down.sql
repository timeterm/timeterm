BEGIN;

ALTER TABLE "student" DROP CONSTRAINT "student_organization_id_zermelo_user_key";

COMMIT;
