-- Create "users" table
CREATE TABLE "public"."users" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "email" text NOT NULL,
  "name" text NOT NULL,
  "password_hash" text NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_users_deleted_at" to table: "users"
CREATE INDEX "idx_users_deleted_at" ON "public"."users" ("deleted_at");
-- Create index "idx_users_email" to table: "users"
CREATE UNIQUE INDEX "idx_users_email" ON "public"."users" ("email");
-- Create "workouts" table
CREATE TABLE "public"."workouts" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "user_id" bigint NOT NULL,
  "name" text NOT NULL,
  "description" text NULL,
  "duration_minutes" bigint NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_workouts_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_workouts_deleted_at" to table: "workouts"
CREATE INDEX "idx_workouts_deleted_at" ON "public"."workouts" ("deleted_at");
-- Create index "idx_workouts_user_id" to table: "workouts"
CREATE INDEX "idx_workouts_user_id" ON "public"."workouts" ("user_id");
