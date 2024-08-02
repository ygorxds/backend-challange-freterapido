CREATE TABLE quotes (
    id SERIAL PRIMARY KEY,
    carrier VARCHAR(255),
    service VARCHAR(255),
    deadline VARCHAR(10),
    price FLOAT8,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
