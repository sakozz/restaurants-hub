BEGIN;

CREATE TABLE
    IF NOT EXISTS users (
        id serial PRIMARY KEY,
        username VARCHAR(50) UNIQUE NOT NULL,
        password VARCHAR(50) NOT NULL,
        email VARCHAR(300) UNIQUE NOT NULL,
        reset_password_token VARCHAR(50),
        reset_password_sent_at timestamp default NULL,
        created_at timestamp NOT NULL  DEFAULT CURRENT_TIMESTAMP,
        updated_at timestamp NOT NULL  DEFAULT now(),
        /* Trackable*/
        sign_in_count integer,
        current_sign_in_at timestamp default NULL,
        last_sign_in_at timestamp default NULL,
        current_sign_in_ip inet default NULL,
        last_sign_in_ip inet default NULL,
        /* confirmable */
        confirmation_token VARCHAR(50) default NULL,
        confirmed_at timestamp default NULL,
        confirmation_sent_at timestamp default NULL
        
    );

CREATE TRIGGER update_user_updated_at BEFORE UPDATE ON users FOR EACH ROW EXECUTE PROCEDURE  update_modified_column();

COMMIT;