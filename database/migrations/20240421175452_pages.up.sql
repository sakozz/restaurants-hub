BEGIN;

CREATE TABLE
    IF NOT EXISTS pages (
        id serial PRIMARY KEY,
        title VARCHAR(50) NOT NULL,
        slug VARCHAR(50) NOT NULL,
        excerpt VARCHAR(2000) NOT NULL,
        body VARCHAR(5000) NOT NULL,
        visibility VARCHAR(50) NOT NULL,
        created_at timestamp NOT NULL DEFAULT now (),
        updated_at timestamp NOT NULL DEFAULT now (),
        deleted_at timestamp,
        profile_id int NOT NULL,
        restaurant_id int NOT NULL,
        CONSTRAINT fk_restaurant FOREIGN KEY (restaurant_id) REFERENCES restaurants (id),
        CONSTRAINT fk_author FOREIGN KEY (profile_id) REFERENCES profiles (id)
    );

CREATE TRIGGER update_pages_updated_at BEFORE
UPDATE ON pages FOR EACH ROW EXECUTE PROCEDURE update_modified_column ();

COMMIT;