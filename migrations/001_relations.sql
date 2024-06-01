CREATE SCHEMA IF NOT EXISTS post_service;

SET search_path TO post_service;

CREATE EXTENSION IF NOT EXISTS ltree;
CREATE EXTENSION IF NOT EXISTS moddatetime;

CREATE TABLE IF NOT EXISTS post(
	id serial PRIMARY KEY,
	author text NOT NULL,
	title text NOT NULL,
	content text NOT NULL,
	allow_comment bool NOT NULL,
	created_at timestamptz NOT NULL DEFAULT NOW(),
	updated_at timestamptz NOT NULL DEFAULT NOW(),
	deleted_at timestamptz
);

CREATE OR REPLACE TRIGGER modify_post_updated_at
	BEFORE UPDATE
	ON post
	FOR EACH ROW
EXECUTE PROCEDURE moddatetime(updated_at);

CREATE TABLE IF NOT EXISTS comment(
	id serial PRIMARY KEY,
	author text NOT NULL,
	content text NOT NULL,
	path ltree,
	created_at timestamptz NOT NULL DEFAULT NOW(),
	updated_at timestamptz NOT NULL DEFAULT NOW(),
	deleted_at timestamptz
);

CREATE OR REPLACE TRIGGER modify_comment_updated_at
	BEFORE UPDATE
	ON comment
	FOR EACH ROW
EXECUTE PROCEDURE moddatetime(updated_at);
