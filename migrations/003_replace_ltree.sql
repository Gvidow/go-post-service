SET search_path TO post_service;

DROP INDEX IF EXISTS comment_path_index;

CREATE EXTENSION IF NOT EXISTS intarray;

ALTER TABLE comment ADD COLUMN new_path int[];

UPDATE comment SET new_path = string_to_array(ltree2text(path), '.')::int[];

ALTER TABLE comment DROP COLUMN path;

ALTER TABLE comment RENAME COLUMN new_path TO path;

ALTER TABLE comment ALTER COLUMN path SET NOT NULL;

DROP EXTENSION IF EXISTS ltree;

CREATE INDEX IF NOT EXISTS comment_path_index
ON comment USING GIN (path);
