BEGIN;

-- DML
INSERT INTO payment_system_user(created_at, updated_at, _role, status)
VALUES
    (now(), NULL, 'ADMIN'::user_role,'ACTIVE'::user_status),
    (now(), NULL, 'MERCHANT'::user_role,'ACTIVE'::user_status),
    (now(), NULL, 'MERCHANT'::user_role,'ACTIVE'::user_status),
    (now(), NULL, 'MERCHANT'::user_role,'INACTIVE'::user_status);

INSERT INTO merchant(user_id, email, name, description)
VALUES
    (2, 'merchant1@gmail.com', 'Merchant One', 'Description'),
    (3, 'merchant2@gmail.com', 'Merchant Two', 'Description'),
    (4, 'merchant3@gmail.com', 'Merchant Three', 'Description');

INSERT INTO transaction(created_at, updated_at, ext_uuid, merchant_id, belongs_to, customer_email, customer_phone, amount, status, _type)
VALUES
    (now(), NULL, gen_random_uuid(), 2, NULL, 'customer1@gmail.com', '1111-22-33', 200, 'APPROVED'::transaction_status, 'AUTHORIZE'::transaction_type),
    (now(), NULL, gen_random_uuid(), 2, NULL, 'customer2@gmail.com', '1111-22-33', 500, 'APPROVED'::transaction_status, 'AUTHORIZE'::transaction_type),
    (now(), NULL, gen_random_uuid(), 2, 2, 'customer2@gmail.com', '1111-22-33', 500, 'APPROVED'::transaction_status, 'CHARGE'::transaction_type),
    (now(), NULL, gen_random_uuid(), 2, NULL, 'customer2@gmail.com', '1111-22-33', 600, 'APPROVED'::transaction_status, 'AUTHORIZE'::transaction_type),
    (now(), NULL, gen_random_uuid(), 2, 4, 'customer2@gmail.com', '1111-22-33', 600, 'APPROVED'::transaction_status, 'CHARGE'::transaction_type),
    (now(), NULL, gen_random_uuid(), 2, NULL, 'customer3@gmail.com', '1111-22-33', 100, 'REVERSED'::transaction_status, 'AUTHORIZE'::transaction_type),
    (now(), NULL, gen_random_uuid(), 2, 6, 'customer3@gmail.com', '1111-22-33', 100, 'REVERSED'::transaction_status, 'CHARGE'::transaction_type),
    (now(), NULL, gen_random_uuid(), 2, 7, 'customer3@gmail.com', '1111-22-33', 100, 'REVERSED'::transaction_status, 'REFUND'::transaction_type);

COMMIT;