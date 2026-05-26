CREATE TABLE "transfer_details" (
  "id" INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  "transaction_id" int NOT NULL,
  "recipient_wallet_id" uuid NOT NULL,
  "notes" TEXT,
  "created_at" timestamp DEFAULT (now()),
  FOREIGN KEY ("transaction_id") REFERENCES "transactions" ("id") DEFERRABLE INITIALLY IMMEDIATE,
  FOREIGN KEY ("recipient_wallet_id") REFERENCES "wallets" ("id") DEFERRABLE INITIALLY IMMEDIATE
);
