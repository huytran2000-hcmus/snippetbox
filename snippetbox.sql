REVOKE CREATE ON SCHEMA public FROM PUBLIC;
REVOKE ALL ON DATABASE snippetbox FROM PUBLIC;

CREATE SCHEMA app;
CREATE TABLE app.snippets (
    id serial NOT NULL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    created TIMESTAMP NOT NULL,
    expired TIMESTAMP NOT NULL
);

-- Add an index on the created column.
CREATE INDEX idx_snippets_created ON app.snippets(created);

INSERT INTO app.snippets (title, content, created, expires) VALUES (
    'An old silent pond',
    'An old silent pond...\nA frog jumps into the pond,\nsplash! Silence again.\n\n– Matsuo Bashō',
    NOW(),
    NOW() + INTERVAL '1 YEAR'
);

INSERT INTO app.snippets (title, content, created, expires) VALUES (
    'Over the wintry forest',
    'Over the wintry\nforest, winds howl in rage\nwith no leaves to blow.\n\n– Natsume Soseki',
    NOW(),
    NOW() + INTERVAL '1 YEAR'
);

INSERT INTO app.snippets (title, content, created, expires) VALUES (
    'First autumn morning',
    'First autumn morning\nthe mirror I stare into\nshows my father''s face.\n\n– Murakami Kijo',
    NOW(),
    NOW() + INTERVAL '7 DAY'
);


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
