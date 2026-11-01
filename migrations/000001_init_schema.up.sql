CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    password TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE UNIQUE INDEX users_email_key ON users (lower(email));
/*
we do not use 
> email VARCHAR(255) NOT NULL UNIQUE
because for this foo@gmail.com and Foo@gmail.com are different, we want uniqueness for any case
that is why we use UNIQUE INDEX on lower(email) 
*/