BEGIN;

SET client_encoding = 'LATIN1';

CREATE TABLE city (
    id SERIAL,
    baseID BIGINT,
    Expression text NOT NULL,
    Value text,
    Err boolean NOT NULL DEFAULT false,
    CurrentResult text
);

CREATE TABLE config (
    id SERIAL,
    Plus BIGINT,
    Minus BIGINT,
    Multiply BIGINT,
    Divide BIGINT
);