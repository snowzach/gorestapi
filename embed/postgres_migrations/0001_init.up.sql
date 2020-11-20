CREATE TABLE IF NOT EXISTS thing (
  id TEXT PRIMARY KEY NOT NULL,
  created timestamp with time zone default NOW(),
  updated timestamp with time zone default NOW(),
  name TEXT,
  description TEXT
);

CREATE TABLE IF NOT EXISTS widget (
  id TEXT PRIMARY KEY NOT NULL,
  created timestamp with time zone default NOW(),
  updated timestamp with time zone default NOW(),
  name TEXT,
  description TEXT,
  thing_id TEXT
);

ALTER TABLE ONLY widget ADD CONSTRAINT fkey_widget_thing_id FOREIGN KEY (thing_id) REFERENCES public.thing(id) ON DELETE CASCADE;
