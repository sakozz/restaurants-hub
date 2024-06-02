BEGIN;

CREATE TABLE
    IF NOT EXISTS invitations (
        id serial PRIMARY KEY,
        role VARCHAR(50) NOT NULL DEFAULT 'public',
        email VARCHAR(300) UNIQUE NOT NULL,
        token VARCHAR(300) NOT NULL,
        created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at timestamp NOT NULL DEFAULT now (),
        expires_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP + interval '1' day
    );

CREATE TRIGGER update_invitation_updated_at BEFORE
UPDATE ON invitations FOR EACH ROW EXECUTE PROCEDURE update_modified_column ();

COMMIT;