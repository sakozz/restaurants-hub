BEGIN;

CREATE TABLE
    IF NOT EXISTS resturants (
        id serial PRIMARY KEY,
        profile_id int NOT NULL,
        name VARCHAR(50) NOT NULL,
        description VARCHAR(50) NOT NULL,
        address VARCHAR(500),
        email VARCHAR(50) UNIQUE NOT NULL,
        phone VARCHAR(20) UNIQUE NOT NULL,
        mobile VARCHAR(20),
        website VARCHAR(500),
        facebook_link VARCHAR(500),
        instagram_link VARCHAR(500),
        created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at timestamp NOT NULL DEFAULT now (),
        deleted_at timestamp,
        CONSTRAINT fk_profile FOREIGN KEY (profile_id) REFERENCES profiles (id)
    );

CREATE TRIGGER update_resturant_updated_at BEFORE
UPDATE ON resturants FOR EACH ROW EXECUTE PROCEDURE update_modified_column ();

COMMIT;