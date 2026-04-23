CREATE TABLE public.app_updates (
  id integer NOT NULL DEFAULT nextval('app_updates_id_seq'::regclass),
  version text NOT NULL UNIQUE,
  title text NOT NULL,
  description text NOT NULL,
  features text NOT NULL,
  released_at timestamp with time zone NOT NULL DEFAULT now(),
  CONSTRAINT app_updates_pkey PRIMARY KEY (id)
);
CREATE TABLE public.categories (
  name text NOT NULL,
  color text NOT NULL,
  CONSTRAINT categories_pkey PRIMARY KEY (name)
);
CREATE TABLE public.patrol_plan_pins (
  id integer NOT NULL DEFAULT nextval('patrol_plan_pins_id_seq'::regclass),
  patrol_plan_id integer NOT NULL,
  pin_id integer NOT NULL,
  sort_order integer NOT NULL DEFAULT 0,
  created_at timestamp with time zone NOT NULL DEFAULT now(),
  CONSTRAINT patrol_plan_pins_pkey PRIMARY KEY (id),
  CONSTRAINT patrol_plan_pins_patrol_plan_id_fkey FOREIGN KEY (patrol_plan_id) REFERENCES public.patrol_plans(id),
  CONSTRAINT patrol_plan_pins_pin_id_fkey FOREIGN KEY (pin_id) REFERENCES public.pins(id)
);
CREATE TABLE public.patrol_plans (
  id integer NOT NULL DEFAULT nextval('patrol_plans_id_seq'::regclass),
  name text NOT NULL,
  date text NOT NULL,
  created_at timestamp with time zone NOT NULL DEFAULT now(),
  updated_at timestamp with time zone NOT NULL DEFAULT now(),
  CONSTRAINT patrol_plans_pkey PRIMARY KEY (id)
);
CREATE TABLE public.pins (
  id integer NOT NULL DEFAULT nextval('pins_id_seq'::regclass),
  title text NOT NULL,
  description text,
  lat double precision NOT NULL,
  lng double precision NOT NULL,
  category text NOT NULL,
  image_url text,
  created_at timestamp with time zone NOT NULL DEFAULT now(),
  updated_at timestamp with time zone NOT NULL DEFAULT now(),
  version integer NOT NULL DEFAULT 1,
  visits_count integer DEFAULT 0,
  CONSTRAINT pins_pkey PRIMARY KEY (id)
);
CREATE TABLE public.streetwork_stats (
  id integer NOT NULL DEFAULT nextval('streetwork_stats_id_seq'::regclass),
  worker_name text NOT NULL,
  month text NOT NULL,
  interactions integer NOT NULL DEFAULT 0,
  new_contacts integer NOT NULL DEFAULT 0,
  interventions integer NOT NULL DEFAULT 0,
  created_at timestamp with time zone NOT NULL DEFAULT now(),
  updated_at timestamp with time zone NOT NULL DEFAULT now(),
  avatar text,
  bg_color text,
  CONSTRAINT streetwork_stats_pkey PRIMARY KEY (id)
);
CREATE TABLE public.user_updates_viewed (
  id integer NOT NULL DEFAULT nextval('user_updates_viewed_id_seq'::regclass),
  user_id integer NOT NULL,
  update_id integer NOT NULL,
  viewed_at timestamp with time zone NOT NULL DEFAULT now(),
  CONSTRAINT user_updates_viewed_pkey PRIMARY KEY (id),
  CONSTRAINT user_updates_viewed_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id),
  CONSTRAINT user_updates_viewed_update_id_fkey FOREIGN KEY (update_id) REFERENCES public.app_updates(id)
);
CREATE TABLE public.users (
  id integer NOT NULL DEFAULT nextval('users_id_seq'::regclass),
  email text NOT NULL UNIQUE,
  password_hash text NOT NULL,
  must_change_password boolean DEFAULT true,
  created_at timestamp with time zone NOT NULL DEFAULT now(),
  last_login timestamp with time zone,
  CONSTRAINT users_pkey PRIMARY KEY (id)
);
CREATE TABLE public.visits (
  id integer NOT NULL DEFAULT nextval('visits_id_seq'::regclass),
  pin_id integer NOT NULL,
  name text NOT NULL,
  note text,
  image_url text,
  visited_at timestamp with time zone NOT NULL DEFAULT now(),
  CONSTRAINT visits_pkey PRIMARY KEY (id),
  CONSTRAINT visits_pin_id_fkey FOREIGN KEY (pin_id) REFERENCES public.pins(id)
);