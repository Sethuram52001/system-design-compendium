BEGIN;

SELECT pg_sleep(25);

UPDATE orders
SET status = 'cancelled'
WHERE id = 1;

COMMIT;