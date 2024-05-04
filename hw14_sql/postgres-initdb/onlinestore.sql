DO $$ BEGIN
    IF EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'manager') THEN
        RAISE NOTICE 'Роль существует.';	
    ELSE
        CREATE ROLE manager WITH LOGIN;
		ALTER USER manager WITH PASSWORD 'password';
        RAISE NOTICE 'Роль manager не существует.';
    END IF;
END $$;
\dg
DROP DATABASE IF EXISTS onlinestore;
CREATE DATABASE onlinestore OWNER manager LOCALE 'en_US.utf8' ENCODING 'SQL_ASCII' TEMPLATE template0;

\l
GRANT ALL PRIVILEGES ON DATABASE onlinestore TO manager;
\c onlinestore manager

CREATE TABLE users
(
    id SERIAL NOT NULL PRIMARY KEY,
    name VARCHAR(60),
    email VARCHAR(30) NOT NULL,
    password VARCHAR(30) NOT NULL,
    UNIQUE(email)
);

CREATE TABLE orders
(
    id SERIAL NOT NULL PRIMARY KEY,
    user_id INT NOT NULL,
    order_date TIMESTAMP WITHOUT TIME ZONE NULL,
    total_amount REAL NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE products
(
    id SERIAL NOT NULL PRIMARY KEY,
    name VARCHAR(200),
    price BIGINT
);

CREATE TABLE orderproducts
(
    id SERIAL NOT NULL PRIMARY KEY,
    order_id INT NOT NULL,
    product_id INT NOT NULL,
    FOREIGN KEY (order_id) REFERENCES orders (id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products (id) ON DELETE CASCADE
);

CREATE INDEX idx_users_email ON users (email);
CREATE INDEX idx_orders_user_id ON orders (user_id);
CREATE INDEX idx_products_name ON products (name);
CREATE INDEX idx_orderproducts_order_id ON orderproducts (order_id);
CREATE INDEX idx_orderproducts_product_id ON orderproducts (product_id);
\dt
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO manager;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO manager;
-- GRANT ALL ON ALL FUNCTIONS IN SCHEMA public to manager;
-- ALTER DATABASE shop OWNER TO manager;
\z
\c onlinestore manager
INSERT INTO users (id, name, email, password) VALUES 
(1, 'Vasiliy Ivanov', 'vi@mail.ru', '12345678'),
(2, 'Pavel Petrov', 'pp@mail.ru', '87654321');
INSERT INTO orders (id, user_id, order_date, total_amount) VALUES 
(1, 1, NOW(), 100),
(2, 2, NOW(), 200),
(3, 1, '2024-05-03', 100);
INSERT INTO products (id, name, price) VALUES 
(1, 'Phone', 500.0),
(2, 'Laptop', 1200.0),
(3, 'Headphones', 100.0);
INSERT INTO orderproducts (order_id, product_id) VALUES
(1, 1),
(1, 2),
(2, 2),
(3, 3); 
SELECT * FROM orderproducts;
SELECT u.name, COUNT(o.id) 
FROM users u
JOIN orders o ON u.id = o.user_id
GROUP BY u.name;
