BEGIN;

CREATE TABLE
    IF NOT EXISTS dishes (
        id serial PRIMARY KEY,
        restaurant_id int NOT NULL,
        name VARCHAR(50) NOT NULL,
        description VARCHAR(5000) NOT NULL,
        price INT NOT NULL,
        category VARCHAR(50),
        tags VARCHAR(50),
        website VARCHAR(500),
        enable_reviews BOOLEAN NOT NULL DEFAULT FALSE,
        published timestamp NOT NULL DEFAULT now (),
        created_at timestamp NOT NULL DEFAULT now (),
        updated_at timestamp NOT NULL DEFAULT now (),
        deleted_at timestamp,
        CONSTRAINT fk_restaurant FOREIGN KEY (restaurant_id) REFERENCES restaurants (id)
    );

CREATE TRIGGER update_dish_updated_at BEFORE
UPDATE ON dishes FOR EACH ROW EXECUTE PROCEDURE update_modified_column ();

COMMIT;