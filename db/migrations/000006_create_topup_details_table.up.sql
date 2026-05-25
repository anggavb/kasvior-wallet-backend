CREATE TABLE "topup_details" (
  "id" INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  "transaction_id" int NOT NULL,
  "payment_method_id" int NOT NULL,
  "discount" decimal NOT NULL,
  "tax" decimal NOT NULL,
  "sub_total" decimal NOT NULL,
  "created_at" timestamp DEFAULT (now()),
  FOREIGN KEY ("transaction_id") REFERENCES "transactions" ("id") DEFERRABLE INITIALLY IMMEDIATE,
  FOREIGN KEY ("payment_method_id") REFERENCES "payment_methods" ("id") DEFERRABLE INITIALLY IMMEDIATE
);
