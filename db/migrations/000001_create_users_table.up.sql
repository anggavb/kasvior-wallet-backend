CREATE TABLE "users" (
  "id" INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  "fullname" VARCHAR(255),
  "email" VARCHAR(255) UNIQUE NOT NULL,
  "password" VARCHAR(255) NOT NULL,
  "pin" char(6),
  "phone_number" VARCHAR(255),
  "photo" VARCHAR(255),
  "is_verified" boolean DEFAULT false,
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp
);
