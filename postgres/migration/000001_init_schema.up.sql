CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "email" varchar UNIQUE NOT NULL,
  "hashed_password" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "followers" (
  "follower_id" bigint NOT NULL,
  "following_id" bigint NOT NULL,

  unique ("follower_id", "following_id")
);

ALTER TABLE "followers" ADD FOREIGN KEY ("follower_id") REFERENCES "users" ("id");

ALTER TABLE "followers" ADD FOREIGN KEY ("following_id") REFERENCES "users" ("id");