BEGIN;
CREATE TABLE IF NOT EXISTS sessions (
        id serial PRIMARY KEY,
        provider VARCHAR(300),
        user_id int NOT NULL,
        email VARCHAR(300),
        access_token VARCHAR(4000),
        access_token_secret VARCHAR(4000) DEFAULT '',
        refresh_token VARCHAR(4000),
        id_token VARCHAR(4000),
        expires_at timestamp NOT NULL  DEFAULT now(),
        created_at timestamp NOT NULL  DEFAULT now(),
        updated_at timestamp NOT NULL  DEFAULT now(),
        CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id)
 );


CREATE TRIGGER update_session_updated_at BEFORE UPDATE ON sessions FOR EACH ROW EXECUTE PROCEDURE  update_modified_column();

COMMIT;
