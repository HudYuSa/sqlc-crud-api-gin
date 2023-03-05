CREATE TABLE
    "users" (
        "id" uuid NOT NULL DEFAULT (uuid_generate_v4 ()),
        "name" varchar NOT NULL,
        "email" varchar NOT NULL,
        "photo" varchar NOT NULL,
        "verified" boolean NOT NULL,
        "password" varchar NOT NULL,
        "role" varchar NOT NULL,
        "created_at" timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
        "updated_at" timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
        CONSTRAINT "users_pkey" PRIMARY KEY ("id")
    );

CREATE UNIQUE INDEX "users_email_key" ON "users" ("email");