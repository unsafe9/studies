CREATE USER admin;
CREATE DATABASE test_db OWNER admin;

\connect test_db

CREATE TABLE greeting
(
    id      BIGSERIAL PRIMARY KEY,
    content TEXT      NOT NULL
);
