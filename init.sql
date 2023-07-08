--name: create-schema
CREATE SCHEMA IF NOT EXISTS cinema_management;

-- name: create-sales-table
CREATE TABLE IF NOT EXISTS cinema_management.article_sales
(
    id    bigserial PRIMARY KEY,
    name  text    NOT NULL,
    count integer NOT NULL,
    time  timestamp DEFAULT NOW()
);

-- name: create-register-table
CREATE TABLE IF NOT EXISTS cinema_management.cash_registers
(
    id          uuid DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    name        text                           NOT NULL,
    description text
);

-- name: create-article-table
CREATE TABLE IF NOT EXISTS cinema_management.articles
(
    id    uuid             DEFAULT gen_random_uuid() NOT NULL
        PRIMARY KEY,
    name  text                                       NOT NULL,
    price double precision DEFAULT 0.00              NOT NULL,
    icon  text
);

-- name: create-transaction-table
CREATE TABLE IF NOT EXISTS cinema_management.transactions
(
    id          uuid DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    title       text                           NOT NULL,
    description text,
    amount      numeric,
    by          text                           NOT NULL,
    register    uuid
        REFERENCES cinema_management.cash_registers
            ON UPDATE RESTRICT ON DELETE RESTRICT
);
