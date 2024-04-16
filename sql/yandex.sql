
CREATE TABLE Expressions (
    id SERIAL,
    baseID BIGINT,
    Expression text NOT NULL,
    Login text NOT NULL,
    Value double precision,
    Err boolean NOT NULL DEFAULT false,
    CurrentResult text
);

CREATE TABLE Configs (
    id SERIAL,
    Plus BIGINT,
    Minus BIGINT,
    Multiply BIGINT,
    Divide BIGINT,
    MaxID BIGINT
);

CREATE TABLE Users (
    id SERIAL,
    Login text NOT NULL,
    Password bytea NOT NULL
);
