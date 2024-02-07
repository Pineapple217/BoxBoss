CREATE TABLE IF NOT EXISTS repositories (
  id            INTEGER PRIMARY KEY AUTOINCREMENT,
  name          text    NOT NULL,
  url           text NOT NULL,
  container_repo text,
  container_tag text,
  container_id  text,
  compose_file  text,
  compose_service text
);