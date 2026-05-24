CREATE TABLE "payment_methods" (
  "id" INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  "name" VARCHAR(255) NOT NULL,
  "logo" VARCHAR(255),
  "method" method_type NOT NULL,
  "tax" decimal DEFAULT 0 NOT NULL,
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp
);
