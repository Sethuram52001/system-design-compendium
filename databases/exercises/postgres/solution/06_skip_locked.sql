BEGIN;

SELECT pg_sleep(10);

SELECT id
FROM orders
WHERE status = 'pending'
ORDER BY id
FOR UPDATE SKIP LOCKED
LIMIT 1;

SELECT pg_sleep(20);

COMMIT;