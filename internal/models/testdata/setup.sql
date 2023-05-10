-- CREATE SCHEMA app;

SET search_path TO app;

CREATE TABLE snippets (
    id serial NOT NULL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    created TIMESTAMP NOT NULL,
    expires TIMESTAMP NOT NULL
);

CREATE INDEX idx_snippets_created ON snippets(created);

CREATE TABLE users (
    id serial NOT NULL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    hashed_password CHAR(60) NOT NULL,
    created TIMESTAMP NOT NULL
);

ALTER TABLE users ADD CONSTRAINT users_uc_email UNIQUE (email);

INSERT INTO users (name, email, hashed_password, created) VALUES (
    'Alice Jones',
    'alice@example.com',
    '$2a$12$NuTjWXm3KKntReFwyBVHyuf/to.HEwTy.eS206TNfkGfr6HzGJSWG',
    '2023-05-09 10:00:00'
);

-- CREATE ROLE test_readwrite;
-- GRANT CONNECT ON DATABASE test_snippetbox TO test_readwrite;
-- GRANT USAGE, CREATE ON SCHEMA app TO test_readwrite;
-- GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA app TO test_readwrite;
-- ALTER DEFAULT PRIVILEGES IN SCHEMA app GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO test_readwrite;
-- GRANT USAGE ON ALL SEQUENCES IN SCHEMA app TO test_readwrite;
-- ALTER DEFAULT PRIVILEGES IN SCHEMA app GRANT USAGE ON SEQUENCES TO test_readwrite;
