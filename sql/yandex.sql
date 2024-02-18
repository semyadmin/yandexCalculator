
CREATE TABLE Expressions (
    id SERIAL,
    baseID BIGINT,
    Expression text NOT NULL,
    Value text,
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
