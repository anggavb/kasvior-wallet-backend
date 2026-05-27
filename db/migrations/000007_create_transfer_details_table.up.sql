CREATE TABLE "transfer_details" (
  "id" INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  "transaction_id" int NOT NULL,
  "recipient_wallet_id" uuid NOT NULL,
  "notes" TEXT,
  "created_at" timestamp DEFAULT (now()),
  FOREIGN KEY ("transaction_id") REFERENCES "transactions" ("id") DEFERRABLE INITIALLY IMMEDIATE,
  FOREIGN KEY ("recipient_wallet_id") REFERENCES "wallets" ("id") DEFERRABLE INITIALLY IMMEDIATE
);

CREATE INDEX "idx_transfer_details_recipient_wallet_transaction"
  ON "transfer_details" ("recipient_wallet_id", "transaction_id");

CREATE INDEX "idx_transfer_details_transaction_id"
  ON "transfer_details" ("transaction_id");
