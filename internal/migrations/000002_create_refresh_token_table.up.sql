CREATE TABLE "refreshToken" (
    "createdAt" timestamptz NOT NULL,
    "updatedAt" timestamptz NOT NULL,
    "id" SERIAL PRIMARY KEY,
    "userId" UUID NOT NULL,
    "value" VARCHAR(255) NOT NULL,
    "clientInfo" VARCHAR(255)
);