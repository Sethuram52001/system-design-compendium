BEGIN;

SELECT *
FROM orders
WHERE id = 1
FOR UPDATE;

SELECT pg_sleep(15);

UPDATE orders
SET status = 'paid'
WHERE id = 1;

COMMIT;
