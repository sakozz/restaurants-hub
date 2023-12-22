BEGIN;

CREATE TABLE
    IF NOT EXISTS tokens (
        id serial PRIMARY KEY,
        user_id int NOT NULL,
        token VARCHAR(300) NOT NULL,
        client VARCHAR(50) NOT NULL,
        created_at timestamp NOT NULL  DEFAULT now(),
        updated_at timestamp NOT NULL  DEFAULT now(),
        expires_at timestamp NOT NULL  DEFAULT now(),
        CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id)
    );
  
CREATE TRIGGER update_token_updated_at BEFORE UPDATE ON tokens FOR EACH ROW EXECUTE PROCEDURE  update_modified_column();

COMMIT;