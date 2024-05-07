# Postgres notes

- [Руководство по PostgreSQL](https://metanit.com/sql/postgresql/)
- [Транзакции](https://www.postgresql.org/docs/current/tutorial-transactions.html)
- [Уровни изоляции транзакций с примерами на PostgreSQL](https://habr.com/ru/articles/317884/)
- [MVCC-1. Изоляция](https://habr.com/ru/companies/postgrespro/articles/442804/)
- [PostgreSQL Antipatterns: уникальные идентификаторы](https://habr.com/ru/companies/tensor/articles/515786/)
- [Функции для работы с последовательностями](https://postgrespro.ru/docs/postgrespro/9.5/functions-sequence)

## Запуск docker контейнера

```bash
mkdir postgres-initdb
mv onlinestore.sql postgres-initdb/
docker compose down; docker rm postgresql; docker volume rm -f otus_pgadmin_data otus_postgresql_data otus_postgresql_log; docker compose up postgresql
docker exec -it postgresql sh -c 'exec psql -h "$POSTGRES_PORT_5432_TCP_ADDR" -p "$POSTGRES_PORT_5432_TCP_PORT" -U postgres'
\c onlinestore manager
```

## SQL запросы

```sql
--- запросы на вставку, редактирование и удаление пользователей
INSERT INTO users (name, email, password) VALUES ('Alexey Sidorov', LOWER('AS@gmail.com'), '42145623');
UPDATE users SET email='pp@yandex.ru' WHERE id=2;
DELETE FROM users WHERE id=3;

--- запросы на вставку, редактирование и удаление продуктов
INSERT INTO products (id, name, price) VALUES (4, 'Monitor', 750.50);
---
BEGIN;
SELECT * FROM products WHERE id=4 FOR UPDATE;
UPDATE products SET name='Television' WHERE id=4;
COMMIT;
---
ROLLBACK;

DELETE FROM products WHERE id=4;

--- запрос на сохранение и удаление заказов
BEGIN;
INSERT INTO orders (user_id, order_date, total_amount) VALUES (1, now(), 99);
INSERT INTO orderproducts (order_id, product_id) VALUES (currval(pg_get_serial_sequence('orders', 'id')), 3);
COMMIT;
---

DELETE FROM orders WHERE id IN (1,2);
--- запрос на выборку пользователей и выборку товаров
SELECT * FROM users WHERE email LIKE LOWER('%VI%');
SELECT * FROM users WHERE id=1;
SELECT * FROM products WHERE name ILIKE '%Laptop%';
--- запрос на выборку заказов по пользователю
SELECT o.id, o.order_date,o.total_amount FROM users u JOIN orders o ON u.id = o.user_id WHERE u.id=1;
--- запрос на выборку статистики по пользователю (общая сумма 'стоимость?' заказов/средняя цена товара)
SELECT u.name, SUM(p.price * o.total_amount) AS total_price_orders FROM orderproducts op LEFT JOIN products p ON op.product_id=p.id RIGHT JOIN orders o ON op.order_id=o.id, users u WHERE o.user_id = u.id AND u.id = 1 GROUP BY u.name;

SELECT u.name, ROUND(AVG(p.price), 2) AS avg_price_per_orders FROM orderproducts op LEFT JOIN products p ON op.product_id=p.id RIGHT JOIN orders o ON op.order_id=o.id, users u WHERE o.user_id = u.id AND u.id = 1 GROUP BY u.name;

SELECT u.name, COUNT(*) AS total_count, SUM(total_amount) AS total_amount FROM orders o INNER JOIN users u ON u.id = o.user_id GROUP BY user_id, u.name;

--- другие запросы
SELECT u.name, SUM(o.total_amount) FROM users u JOIN orders o ON u.id = o.user_id WHERE u.id=1 GROUP BY u.name;

SELECT z.product_id, p.price FROM orderproducts AS z LEFT JOIN products AS p ON z.product_id = p.id GROUP BY z.product_id, p.price;

SELECT order_id, price FROM products, orderproducts WHERE orderproducts.product_id=products.id;
SELECT order_id, total_amount FROM orders, orderproducts WHERE orderproducts.order_id=orders.id;
```

## Транзакции

```sql
CREATE OR REPLACE PROCEDURE make_order(IN _user_id int, IN _total_amount int, IN _product_id_list int[])
LANGUAGE plpgsql
AS $$
DECLARE
    _seq INT;
    _product_id INT;
BEGIN
    SELECT nextval(pg_get_serial_sequence('orders', 'id')) INTO _seq;
    INSERT INTO orders (id, user_id, order_date, total_amount) VALUES (_seq, _user_id, now(), _total_amount);
    FOREACH _product_id IN ARRAY _product_id_list
        LOOP
            INSERT INTO orderproducts (order_id, product_id) VALUES (_seq, _product_id);
        END LOOP;
    COMMIT;
END;
$$;

SELECT proname, prosrc
FROM pg_proc
WHERE proname = 'make_order';

SELECT EXISTS (
    SELECT 1
    FROM pg_proc
    WHERE proname = 'make_order'
);

CALL make_order(1, 99, ARRAY[1, 2, 3]);
CALL make_order(2, 500, ARRAY[1, 2, 3]);
```