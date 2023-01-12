BEGIN;

-- User
CREATE TYPE user_role AS ENUM ('ADMIN', 'MERCHANT');
CREATE TYPE user_status AS ENUM ('ACTIVE', 'INACTIVE');

CREATE TABLE payment_system_user(
                                    id BIGINT NOT NULL GENERATED ALWAYS AS IDENTITY,
                                    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
                                    updated_at TIMESTAMP WITH TIME ZONE NULL,

                                    _role user_role NOT NULL,
                                    status user_status NOT NULL
);
ALTER TABLE payment_system_user ADD PRIMARY KEY(id);

-- Merchant
CREATE TABLE merchant(
                         user_id BIGINT NOT NULL,
                         email VARCHAR(255) NOT NULL,
                         name VARCHAR(255) NOT NULL,
                         description VARCHAR(255) NULL
);

ALTER TABLE merchant ADD PRIMARY KEY(user_id);
CREATE UNIQUE INDEX merchant_email_unique ON merchant USING btree(email);

ALTER TABLE merchant ADD CONSTRAINT merchant_user_id_foreign FOREIGN KEY(user_id)
    REFERENCES payment_system_user(id) ON DELETE CASCADE ON UPDATE CASCADE;


-- Transaction
CREATE TYPE transaction_status AS ENUM ('APPROVED', 'REVERSED', 'REFUNDED', 'ERROR');
CREATE TYPE transaction_type AS ENUM ('AUTHORIZE', 'CHARGE', 'REFUND', 'REVERSAL');

CREATE TABLE transaction(
                            id BIGINT NOT NULL GENERATED ALWAYS AS IDENTITY,
                            created_at TIMESTAMP WITH TIME ZONE NOT NULL,
                            updated_at TIMESTAMP WITH TIME ZONE NULL,

                            ext_uuid UUID NOT NULL,
                            merchant_id BIGINT NOT NULL,
                            belongs_to BIGINT NULL,

                            customer_email VARCHAR(255) NOT NULL,
                            customer_phone VARCHAR(255) NOT NULL,
                            amount BIGINT NOT NULL,

                            status transaction_status NOT NULL,
                            _type transaction_type NOT NULL
);
ALTER TABLE transaction ADD PRIMARY KEY(id);
CREATE UNIQUE INDEX transaction_ext_uuid_unique ON transaction USING btree(ext_uuid);
CREATE INDEX transaction_merchant_id_index ON transaction USING btree(merchant_id);

ALTER TABLE transaction ADD CONSTRAINT transaction_merchant_id_foreign FOREIGN KEY(merchant_id) REFERENCES merchant(user_id) ON DELETE RESTRICT;
ALTER TABLE transaction ADD CONSTRAINT transaction_belongs_to_foreign FOREIGN KEY(belongs_to) REFERENCES transaction(id) ON DELETE CASCADE ;

COMMIT;