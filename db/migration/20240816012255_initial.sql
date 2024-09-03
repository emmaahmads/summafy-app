-- +goose Up
-- +goose StatementBegin
DROP TYPE IF EXISTS "user_activity";

CREATE TYPE "user_activity" AS ENUM (
  'uploaded',
  'generated_summary',
  'deleted',
  'change_summary',
  'downloaded'
);

CREATE TABLE "users" (
  "username" varchar UNIQUE PRIMARY KEY,
  "hashed_password" varchar NOT NULL,
  "full_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "activities" (
  "id" BIGSERIAL PRIMARY KEY,
  "username" varchar NOT NULL,
  "activity" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "document_id" bigserial NOT NULL
);

CREATE TABLE "document" (
  "id" BIGSERIAL PRIMARY KEY,
  "username" varchar NOT NULL,
  "is_private" bool NOT NULL DEFAULT true,
  "has_summary" bool NOT NULL DEFAULT false,
  "file_name" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "summary" (
  "id" BIGSERIAL PRIMARY KEY,
  "doc_id" bigserial NOT NULL,
  "param1" bool NOT NULL DEFAULT false,
  "param2" bool NOT NULL DEFAULT false,
  "summary" BYTEA
);
CREATE TABLE"sessions"("id"uuid PRIMARY KEY,"username"varchar NOT NULL,"refresh_token"varchar NOT NULL,"user_agent"varchar NOT NULL,"client_ip"varchar NOT NULL,"is_blocked"boolean NOT NULL DEFAULT false,"expires_at"timestamptz NOT NULL,"created_at"timestamptz NOT NULL DEFAULT (now()) );

CREATE INDEX ON "activities" ("username");

CREATE INDEX ON "document" ("file_name");

CREATE INDEX ON "summary" ("doc_id");

CREATE UNIQUE INDEX "doc_summary" ON "summary" ("doc_id", "summary");

ALTER TABLE "activities" ADD FOREIGN KEY ("username") REFERENCES "users" ("username") ON DELETE CASCADE;

ALTER TABLE "document" ADD FOREIGN KEY ("username") REFERENCES "users" ("username") ON DELETE CASCADE;

ALTER TABLE "activities" ADD FOREIGN KEY ("document_id") REFERENCES "document" ("id");

ALTER TABLE "summary" ADD FOREIGN KEY ("doc_id") REFERENCES "document" ("id") ON DELETE CASCADE;

ALTER TABLE "sessions" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TYPE IF EXISTS "user_activity";
DROP TABLE "sessions";
DROP TABLE "summary";
DROP TABLE "activities";
DROP TABLE "document";
DROP TABLE "users";
-- +goose StatementEnd