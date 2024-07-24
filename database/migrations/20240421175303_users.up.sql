BEGIN;

CREATE TABLE
    IF NOT EXISTS users (
        id serial PRIMARY KEY,
        role VARCHAR(50) NOT NULL DEFAULT 'public',
        first_name VARCHAR(50) NOT NULL,
        last_name VARCHAR(50) NOT NULL,
        email VARCHAR(300) UNIQUE NOT NULL,
        avatar_url VARCHAR(500),
        restaurant_id int DEFAULT NULL,
        created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at timestamp NOT NULL DEFAULT now (),
        deleted_at timestamp
    );

CREATE TRIGGER update_user_updated_at BEFORE
UPDATE ON users FOR EACH ROW EXECUTE PROCEDURE update_modified_column ();

INSERT INTO users (role, first_name, last_name, email, avatar_url)  VALUES ('admin', 'System', 'admin', 'mygmail@gmail.com'  , 'https://avatars.githubusercontent.com/u/2129058?v=4');    
COMMIT;