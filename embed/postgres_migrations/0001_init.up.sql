CREATE TABLE IF NOT EXISTS thing (
  id TEXT PRIMARY KEY NOT NULL,
  created timestamp with time zone default NOW(),
  updated timestamp with time zone default NOW(),
  name TEXT
);

