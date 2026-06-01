EXPLAIN ANALYZE 
SELECT *
FROM orders
WHERE user_id = 1
AND status = 'pending';

CREATE INDEX IF NOT EXISTS orders_user_id_idx ON orders(user_id);
CREATE INDEX IF NOT EXISTS orders_status_idx ON orders(user_id, status);
CREATE INDEX IF NOT EXISTS orders_pending_user_idx ON orders(user_id) WHERE status = 'pending';
CREATE INDEX IF NOT EXISTS users_lower_email_idx ON users(lower(email));
