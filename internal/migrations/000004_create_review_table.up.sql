CREATE TABLE "review" (
    "createdAt"   timestamptz        NOT NULL,
    "updatedAt"   timestamptz        NOT NULL,
    "id"          UUID             PRIMARY KEY,
    "userId"      UUID             NOT NULL,
    "hospitalId"   UUID             NOT NULL,
    "content"     TEXT             NOT NULL,
    "rating"      INT              NOT NULL
);