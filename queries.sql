-- name: get-register-items
SELECT
    *
FROM
    cinema_management.articles;

-- name: get-registers
SELECT
    *
FROM
    cinema_management.cash_registers;

-- name: insert-transaction
INSERT INTO
    cinema_management.transactions(title, description, amount, by, register)
VALUES
    ($1, $2, $3, $4, $5::uuid);

-- name: insert-article-sale
INSERT INTO
    cinema_management.article_sales(name, count)
VALUES
    ($1, $2);

-- name: get-article-statistics
SELECT
    name, sum(count) as count
FROM cinema_management.article_sales
WHERE time BETWEEN $1 AND $2
GROUP BY name
ORDER BY name;