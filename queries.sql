-- name: is-schema-available
SELECT
    COUNT(*)
FROM
    information_schema.schemata
WHERE
    schema_name = ?;