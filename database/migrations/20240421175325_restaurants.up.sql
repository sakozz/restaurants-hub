BEGIN;

CREATE TABLE
    IF NOT EXISTS restaurants (
        id serial PRIMARY KEY,
        manager_id int NOT NULL,
        name VARCHAR(50) NOT NULL,
        description VARCHAR(5000) NOT NULL,
        address JSONB DEFAULT '{}'::jsonb,
        email VARCHAR(50) UNIQUE NOT NULL,
        phone VARCHAR(20) UNIQUE NOT NULL,
        mobile VARCHAR(20) DEFAULT '',
        website VARCHAR(500) DEFAULT '',
        facebook_link VARCHAR(500) DEFAULT '',
        instagram_link VARCHAR(500) DEFAULT '',
        created_at timestamp NOT NULL DEFAULT now (),
        updated_at timestamp NOT NULL DEFAULT now (),
        deleted_at timestamp,
        CONSTRAINT fk_user FOREIGN KEY (manager_id) REFERENCES users (id)
    );

CREATE TRIGGER update_restaurant_updated_at BEFORE
UPDATE ON restaurants FOR EACH ROW EXECUTE PROCEDURE update_modified_column ();

COMMIT;