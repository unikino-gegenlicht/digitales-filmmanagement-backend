-- name: is-schema-available
SELECT
    COUNT(*)
FROM
    information_schema.schemata
WHERE
    schema_name = $1;

-- name: is-table-available
SELECT
    COUNT(*)
FROM
    information_schema.TABLES
WHERE
    table_schema  = $1
AND
    table_name = $2;

-- name: get-register-items
SELECT
    *
FROM
    gegenlicht.register_items;

-- name: get-registers
SELECT
    *
FROM
    gegenlicht.cash_registers;

-- name: insert-transaction
INSERT INTO
    gegenlicht.transactions(title, description, amount, by, register)
VALUES
    ($1, $2, $3, $4, $5::uuid);

-- name: insert-article-sale
INSERT INTO
    gegenlicht.article_sales(name, count)
VALUES
    ($1, $2);

-- name: get-article-statistics
SELECT
    name, sum(count) as count
FROM gegenlicht.article_sales
WHERE time BETWEEN $1 AND $2
GROUP BY name
ORDER BY name;