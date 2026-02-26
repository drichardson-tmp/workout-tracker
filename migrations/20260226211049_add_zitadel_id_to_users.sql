-- Modify "users" table
ALTER TABLE "public"."users" ALTER COLUMN "password_hash" DROP NOT NULL, ADD COLUMN "zitadel_id" text NULL;
-- Create index "idx_users_zitadel_id" to table: "users"
CREATE UNIQUE INDEX "idx_users_zitadel_id" ON "public"."users" ("zitadel_id");
