BEGIN;

CREATE TABLE
    IF NOT EXISTS pages (
        id serial PRIMARY KEY,
        title VARCHAR(50) NOT NULL,
        slug VARCHAR(50) UNIQUE NOT NULL,
        excerpt VARCHAR(2000) NOT NULL,
        body VARCHAR(5000) NOT NULL,
        visibility VARCHAR(50) NOT NULL DEFAULT 'draft',
        created_at timestamp NOT NULL DEFAULT now (),
        updated_at timestamp NOT NULL DEFAULT now (),
        deleted_at timestamp,
        author_id int NOT NULL,
        restaurant_id int NOT NULL,
        parent_page_id int ,
        CONSTRAINT fk_restaurant FOREIGN KEY (restaurant_id) REFERENCES restaurants (id),
        CONSTRAINT fk_author FOREIGN KEY (author_id) REFERENCES users (id),
        CONSTRAINT fk_parent_page FOREIGN KEY (parent_page_id) REFERENCES pages (id)
    );

CREATE TRIGGER update_pages_updated_at BEFORE
UPDATE ON pages FOR EACH ROW EXECUTE PROCEDURE update_modified_column ();

COMMIT;