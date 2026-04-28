CREATE TABLE app_updates (
  id          SERIAL PRIMARY KEY,
  version     TEXT NOT NULL UNIQUE,
  title       TEXT NOT NULL,
  description TEXT NOT NULL,
  features    TEXT NOT NULL,
  released_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE categories (
  name  TEXT PRIMARY KEY,
  color TEXT NOT NULL
);

CREATE TABLE users (
  id                  SERIAL PRIMARY KEY,
  email               TEXT NOT NULL UNIQUE,
  password_hash       TEXT NOT NULL,
  must_change_password BOOLEAN DEFAULT true,
  created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
  last_login          TIMESTAMPTZ
);

CREATE TABLE pins (
  id          SERIAL PRIMARY KEY,
  title       TEXT NOT NULL,
  description TEXT,
  lat         DOUBLE PRECISION NOT NULL,
  lng         DOUBLE PRECISION NOT NULL,
  category    TEXT NOT NULL,
  image_url   TEXT,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  version     INTEGER NOT NULL DEFAULT 1,
  visits_count INTEGER DEFAULT 0
);

CREATE TABLE visits (
  id         SERIAL PRIMARY KEY,
  pin_id     INTEGER NOT NULL REFERENCES pins(id),
  name       TEXT NOT NULL,
  note       TEXT,
  image_url  TEXT,
  visited_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE patrol_plans (
  id         SERIAL PRIMARY KEY,
  name       TEXT NOT NULL,
  date       TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE patrol_plan_pins (
  id            SERIAL PRIMARY KEY,
  patrol_plan_id INTEGER NOT NULL REFERENCES patrol_plans(id),
  pin_id         INTEGER NOT NULL REFERENCES pins(id),
  sort_order     INTEGER NOT NULL DEFAULT 0,
  created_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE streetwork_stats (
  id            SERIAL PRIMARY KEY,
  worker_name   TEXT NOT NULL,
  month         TEXT NOT NULL,
  interactions  INTEGER NOT NULL DEFAULT 0,
  new_contacts  INTEGER NOT NULL DEFAULT 0,
  interventions INTEGER NOT NULL DEFAULT 0,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
  avatar        TEXT,
  bg_color      TEXT,
  UNIQUE (worker_name, month)
);

CREATE TABLE user_updates_viewed (
  id        SERIAL PRIMARY KEY,
  user_id   INTEGER NOT NULL REFERENCES users(id),
  update_id INTEGER NOT NULL REFERENCES app_updates(id),
  viewed_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
