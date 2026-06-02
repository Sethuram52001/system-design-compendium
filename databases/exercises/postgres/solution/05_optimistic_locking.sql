UPDATE orders
SET status = 'cancelled',
    version = version + 1
WHERE id = 2
  AND version = 1;