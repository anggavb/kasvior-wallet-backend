CREATE TYPE "type_transaction" AS ENUM (
  'topup',
  'transfer',
  'receiver'
);

CREATE TYPE "status_transaction" AS ENUM (
  'pending',
  'success',
  'failed'
);

CREATE TABLE "transactions" (
  "id" INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  "wallet_id" uuid NOT NULL,
  "amount" decimal NOT NULL,
  "type" type_transaction NOT NULL,
  "status" status_transaction NOT NULL,
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp,
  FOREIGN KEY ("wallet_id") REFERENCES "wallets" ("id") DEFERRABLE INITIALLY IMMEDIATE
);
