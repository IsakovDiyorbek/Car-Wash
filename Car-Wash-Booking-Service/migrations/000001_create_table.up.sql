CREATE TYPE payment_status AS ENUM ('pending', 'completed', 'failed', 'refunded');

CREATE TABLE payments (
    id uuid PRIMARY KEY, 
    booking_id uuid NOT NULL, 
    amount NUMERIC(10, 2) NOT NULL, 
    status payment_status NOT NULL, -
    payment_method VARCHAR(50) NOT NULL, 
    transaction_id VARCHAR(100) NOT NULL, 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP 
    deleted_at TIMESTAMP DEFAULT NULL
    
);
