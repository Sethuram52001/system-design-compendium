BEGIN;

WITH new_user AS (
    INSERT INTO users(email, name)
    VALUES ('sethu@gmail.com', 'Sethu')
    RETURNING id
)

INSERT INTO orders(user_id, status, total_cents)
SELECT id, 'pending', 1000
FROM new_user;

COMMIT;

BEGIN;

INSERT INTO users(email, name)
VALUES('rollback@gmail.com', 'Rollback Demo');

ROLLBACK;

BEGIN;

INSERT INTO users (email, name)
VALUES ('sethu@gmail.com', 'Duplicate Sethu');

INSERT INTO orders (user_id, status, total_cents)
VALUES (1, 'pending', 1000);

ROLLBACK;

BEGIN;

INSERT INTO users (email, name)
VALUES ('order-fails@example.com', 'Order Fails');

INSERT INTO orders (user_id, status, total_cents)
SELECT id, 'pending', -100
FROM users
WHERE email = 'order-fails@example.com';

ROLLBACK;

SELECT *
FROM users
WHERE email = 'order-fails@example.com';
