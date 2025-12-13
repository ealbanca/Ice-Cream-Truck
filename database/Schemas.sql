-- SQL to create the events table
CREATE TABLE IF NOT EXISTS events (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(20) NOT NULL,
    date TIMESTAMP NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- SQL to create the contact table
CREATE TABLE IF NOT EXISTS contacts (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    reason VARCHAR(255),
    message TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- SQL to create the users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    lastname VARCHAR(100) NOT NULL,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    phone VARCHAR(20),
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- SQL to create the products tables
-- SQL to create the sizes table
CREATE TABLE sizes (
    id SERIAL PRIMARY KEY,
    size_label VARCHAR(50) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    size_img VARCHAR(255) NOT NULL
);
-- Add sizes data
INSERT INTO sizes (size_label, price, size_img)
VALUES ('Small', 10.00, 'small.jpg'),
    ('Medium', 15.00, 'medium.jpeg'),
    ('Large', 20.00, 'large.jpg');
-- SQL to create the flavors table
CREATE TABLE flavors (
    id SERIAL PRIMARY KEY,
    flavor_name VARCHAR(100) NOT NULL,
    falvor_img VARCHAR(255) NOT NULL
);
-- Add flavors data
INSERT INTO flavors (flavor_name, falvor_img)
VALUES ('Vanilla', 'Vanilla.webp'),
    ('Chocolate', 'Chocolate.webp'),
    ('Strawberry', 'Strawberry.webp'),
    ('Mint Chocolate Chip', 'Mint-Chip.webp'),
    ('Cookies and Cream', 'Cookies-N-Cream.webp'),
    ('Pistachio', 'Pistachio.webp'),
    ('Mango', 'Mango.webp'),
    ('Black Raspberry', 'Black-Raspberry.webp');
-- SQL to create the toppings table
CREATE TABLE toppings (
    id SERIAL PRIMARY KEY,
    topping_name VARCHAR(100) NOT NULL,
    additional_price DECIMAL(10, 2) NOT NULL,
    topping_img VARCHAR(255) NOT NULL
);
-- Add toppings data
INSERT INTO toppings (topping_name, additional_price, topping_img)
VALUES ('Sprinkles', 0.50, '/images/build/Sprinkles.jpg'),
    (
        'Chocolate Chips',
        0.75,
        '/images/build/Chocolate-Chips.webp'
    ),
    (
        'Nerds',
        0.60,
        '/images/build/Nerds.jpg'
    ),
    (
        'Nuts',
        0.80,
        '/images/build/Nuts.webp'
    ),
    (
        'Gummi Bears',
        0.70,
        '/images/build/Gummi-Bears.jpg'
    ),
    (
        'Marshmallows',
        0.65,
        '/images/build/Marshmallow.jpg'
    );
-- SQL to create the products table
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    product_name VARCHAR(255) NOT NULL,
    size_id INT REFERENCES sizes(id) Not NULL,
    flavor_id1 INT REFERENCES flavors(id) Not NULL,
    flavor_id2 INT REFERENCES flavors(id),
    flavor_id3 INT REFERENCES flavors(id),
    topping_id1 INT REFERENCES toppings(id) Not NULL,
    topping_id2 INT REFERENCES toppings(id),
    topping_id3 INT REFERENCES toppings(id),
    total_price DECIMAL(10, 2) NOT NULL
);
-- SQL to create the orders table
-- Orders table: one row per order
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    guest_cart_id VARCHAR(64),
    total_price DECIMAL(10, 2) NOT NULL,
    order_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(20) NOT NULL DEFAULT 'cart'
);
-- Order items table: one row per product in an order
CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,
    order_id INT REFERENCES orders(id) ON DELETE CASCADE,
    product_id INT REFERENCES products(id),
    quantity INT NOT NULL DEFAULT 1
);
-- Each guest cart is identified by a guest_cart_id (string), and contains products and quantities
CREATE TABLE guest_carts (
    id SERIAL PRIMARY KEY,
    guest_cart_id VARCHAR(64) NOT NULL,
    product_id INT REFERENCES products(id),
    quantity INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- Create table to check out information
CREATE TABLE checkout_info (
    id SERIAL PRIMARY KEY,
    order_id INT REFERENCES orders(id) ON DELETE CASCADE,
    full_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    credit_card VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);