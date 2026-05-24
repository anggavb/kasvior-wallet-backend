CREATE TYPE "provider_name" AS ENUM (
  'google',
  'facebook'
);

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

CREATE TYPE "method_type" AS ENUM (
  'online',
  'bank'
);
