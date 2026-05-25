CREATE TABLE "password_reset_tokens" (
  "id" INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  "user_id" INT NOT NULL,
  "token_hash" CHAR(64) UNIQUE NOT NULL,
  "expires_at" timestamp NOT NULL,
  "used_at" timestamp,
  "created_at" timestamp DEFAULT (now()),
  FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE DEFERRABLE INITIALLY IMMEDIATE
);

CREATE INDEX "idx_password_reset_tokens_user_id" ON "password_reset_tokens" ("user_id");
