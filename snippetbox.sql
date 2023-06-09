REVOKE CREATE ON SCHEMA public FROM PUBLIC;
REVOKE ALL ON DATABASE snippetbox FROM PUBLIC;

CREATE SCHEMA app;
SET search_path TO app;
CREATE TABLE snippets (
    id serial NOT NULL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    created TIMESTAMP NOT NULL,
    expires TIMESTAMP NOT NULL
);

CREATE INDEX idx_snippets_created ON snippets(created);

INSERT INTO snippets (title, content, created, expires) VALUES (
    'An old silent pond',
    'An old silent pond...\nA frog jumps into the pond,\nsplash! Silence again.\n\n– Matsuo Bashō',
    NOW(),
    NOW() + INTERVAL '1 YEAR'
);

INSERT INTO snippets (title, content, created, expires) VALUES (
    'Over the wintry forest',
    'Over the wintry\nforest, winds howl in rage\nwith no leaves to blow.\n\n– Natsume Soseki',
    NOW(),
    NOW() + INTERVAL '1 YEAR'
);

INSERT INTO snippets (title, content, created, expires) VALUES (
    'First autumn morning',
    'First autumn morning\nthe mirror I stare into\nshows my father''s face.\n\n– Murakami Kijo',
    NOW(),
    NOW() + INTERVAL '7 DAY'
);

INSERT INTO snippets (title, content, created, expires) VALUES (
    'From time to time',
    'From time to time
The clouds give rest
To the moon-beholders.

- Matsu Basho',
    NOW(),
    NOW() + INTERVAL '30 DAY'
);

CREATE TABLE sessions (
    token TEXT PRIMARY KEY,
    data BYTEA NOT NULL,
    expiry TIMESTAMPTZ NOT NULL
);

CREATE INDEX sessions_expiry_idx ON sessions(expiry);

CREATE TABLE users (
    id serial NOT NULL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    hashed_password CHAR(60) NOT NULL,
    created TIMESTAMP NOT NULL
);

ALTER TABLE users ADD CONSTRAINT users_uc_email UNIQUE (email);

CREATE ROLE readonly;
GRANT CONNECT ON DATABASE snippetbox TO readonly;
GRANT USAGE ON SCHEMA app TO readonly;
GRANT SELECT ON ALL TABLES IN SCHEMA app TO readonly;
ALTER DEFAULT PRIVILEGES IN SCHEMA app GRANT SELECT ON TABLES TO readonly;

CREATE ROLE readwrite;
GRANT CONNECT ON DATABASE snippetbox TO readwrite;
GRANT USAGE, CREATE ON SCHEMA app TO readwrite;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA app TO readwrite;
ALTER DEFAULT PRIVILEGES IN SCHEMA app GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO readwrite;
GRANT USAGE ON ALL SEQUENCES IN SCHEMA app TO readwrite;
ALTER DEFAULT PRIVILEGES IN SCHEMA app GRANT USAGE ON SEQUENCES TO readwrite;

CREATE USER app_user WITH PASSWORD 'huy2000';
GRANT readwrite TO app_user;
