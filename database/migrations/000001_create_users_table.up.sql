BEGIN;

CREATE TABLE
    IF NOT EXISTS users (
        id serial PRIMARY KEY,
        first_name VARCHAR(50) NOT NULL,
        last_name VARCHAR(50) NOT NULL,
        email VARCHAR(300) UNIQUE NOT NULL,
        avatar_url VARCHAR(500),
        created_at timestamp NOT NULL  DEFAULT CURRENT_TIMESTAMP,
        updated_at timestamp NOT NULL  DEFAULT now() ,
        deleted_at timestamp 
    );

CREATE TRIGGER update_user_updated_at BEFORE UPDATE ON users FOR EACH ROW EXECUTE PROCEDURE  update_modified_column();

COMMIT;