
CREATE TABLE IF NOT EXISTS Expressions (
    id SERIAL,
    baseID BIGINT,
    Expression text NOT NULL,
    User text NOT NULL,
    Value text,
    Err boolean NOT NULL DEFAULT false,
    CurrentResult text
);

CREATE TABLE IF NOT EXISTS Configs (
    id SERIAL,
    Plus BIGINT,
    Minus BIGINT,
    Multiply BIGINT,
    Divide BIGINT,
    MaxID BIGINT
);

CREATE TABLE IF NOT EXISTS Users (
    id SERIAL,
    Login text NOT NULL,
    Password text NOT NULL
);
