CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    deposit_amount DECIMAL(10, 2) DEFAULT 0
);

CREATE TABLE mesin_bor (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    stock_availability INT NOT NULL,
    rental_costs DECIMAL(10, 2) NOT NULL,
    category VARCHAR(100)
);

CREATE TABLE rental_history (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    drill_id INT NOT NULL,
    rental_date DATE NOT NULL,
    return_date DATE,
    total_cost DECIMAL(10, 2) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (drill_id) REFERENCES mesin_bor(id)
);

CREATE TABLE reviews (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    drill_id INT NOT NULL,
    review_text TEXT,
    rating INT CHECK (rating BETWEEN 1 AND 5),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (drill_id) REFERENCES mesin_bor(id)
);

CREATE TABLE maintenance (
    id INT AUTO_INCREMENT PRIMARY KEY,
    drill_id INT NOT NULL,
    maintenance_date DATE NOT NULL,
    cost DECIMAL(10, 2) NOT NULL,
    description TEXT,
    FOREIGN KEY (drill_id) REFERENCES mesin_bor(id)
);
