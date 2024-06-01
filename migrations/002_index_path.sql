SET search_path TO post_service;

CREATE INDEX IF NOT EXISTS comment_path_index
ON comment USING GIST (path);
