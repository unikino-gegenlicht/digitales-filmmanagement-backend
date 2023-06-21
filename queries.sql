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
