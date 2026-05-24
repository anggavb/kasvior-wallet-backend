CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE "wallets" (
  "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  "user_id" int UNIQUE NOT NULL,
  "balance" decimal DEFAULT 0,
  FOREIGN KEY ("user_id") REFERENCES "users" ("id") DEFERRABLE INITIALLY IMMEDIATE
);
