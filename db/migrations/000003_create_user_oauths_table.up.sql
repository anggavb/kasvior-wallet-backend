CREATE TABLE "user_oauths" (
  "id" INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  "user_id" int NOT NULL,
  "provider" provider_name NOT NULL,
  "access_token" VARCHAR(255) NOT NULL,
  "refresh_token" VARCHAR(255),
  "expires_at" date,
  FOREIGN KEY ("user_id") REFERENCES "users" ("id") DEFERRABLE INITIALLY IMMEDIATE
);
