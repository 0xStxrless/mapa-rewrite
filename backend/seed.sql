-- Categories
INSERT INTO categories (name, color) VALUES
  ('Hotspot',       '#FF5733'),
  ('Shelter',       '#33A1FF'),
  ('Food Point',    '#28A745'),
  ('Risk Zone',     '#FFC300'),
  ('Meeting Point', '#8E44AD');

-- Pins
INSERT INTO pins (title, description, lat, lng, category, image_url) VALUES
  ('Dworzec Główny',        'Główny węzeł komunikacyjny, dużo bezdomnych wieczorami', 54.3571, 18.6444, 'Hotspot',       null),
  ('Noclegownia Wzgórze',   'Schronisko z 40 miejscami, czynne od 18:00',             54.3520, 18.6380, 'Shelter',       null),
  ('Kuchnia Społeczna',     'Ciepły posiłek codziennie 12:00–14:00',                  54.3610, 18.6500, 'Food Point',    null),
  ('Park Kolberga',         'Skupisko osób w kryzysie, szczególnie w weekendy',        54.3490, 18.6410, 'Risk Zone',     null),
  ('Punkt przy Biedronce',  'Regularne spotkania zespołu w każdy wtorek',             54.3555, 18.6460, 'Meeting Point', null),
  ('Poczekalnia SKM',       'Osoby nocujące w budynku stacji',                        54.3580, 18.6420, 'Hotspot',       null),
  ('Ogrzewalnia Caritas',   'Dzienna ogrzewalnia, max 20 osób',                       54.3530, 18.6390, 'Shelter',       null);

-- Visits
INSERT INTO visits (pin_id, name, note, visited_at) VALUES
  (1, 'Anna K.',   'Kontakt z 3 osobami, rozdano ulotki',          now() - interval '2 days'),
  (1, 'Marek W.',  'Mężczyzna ok. 50 lat, odmówił pomocy',         now() - interval '1 day'),
  (3, 'Anna K.',   'Wydano 12 posiłków, brak incydentów',          now() - interval '3 days'),
  (4, 'Tomasz R.', 'Interwencja — wezwano pogotowie',              now() - interval '5 hours'),
  (2, 'Marek W.',  'Przekazano info o noclegowni nowej osobie',    now() - interval '6 days'),
  (6, 'Tomasz R.', 'Patrol nocny, 5 osób w poczekalni',           now() - interval '12 hours');

-- Streetwork stats
INSERT INTO streetwork_stats (worker_name, month, interactions, new_contacts, interventions, avatar, bg_color) VALUES
  ('Anna K.',   '2025-03', 42, 8, 2, null, '#33A1FF'),
  ('Marek W.',  '2025-03', 35, 5, 1, null, '#FF5733'),
  ('Tomasz R.', '2025-03', 28, 3, 3, null, '#28A745'),
  ('Anna K.',   '2025-04', 38, 6, 0, null, '#33A1FF'),
  ('Marek W.',  '2025-04', 41, 9, 2, null, '#FF5733'),
  ('Tomasz R.', '2025-04', 33, 4, 1, null, '#28A745');

-- Patrol plans
INSERT INTO patrol_plans (name, date) VALUES
  ('Patrol Śródmieście – maj',  '2025-05-03'),
  ('Patrol Nocny Piątek',       '2025-05-09'),
  ('Obchód Stacja SKM',         '2025-05-15');

-- Pin assignments to patrol plans (sort_order = intended walking order)
INSERT INTO patrol_plan_pins (patrol_plan_id, pin_id, sort_order) VALUES
  (1, 1, 1),
  (1, 3, 2),
  (1, 5, 3),
  (2, 6, 1),
  (2, 4, 2),
  (3, 6, 1),
  (3, 2, 2);

-- App updates
INSERT INTO app_updates (version, title, description, features, released_at) VALUES
  ('1.0.0', 'Pierwsze wydanie',   'Podstawowa wersja aplikacji',       'Mapa, piny, wizyty',              now() - interval '30 days'),
  ('1.1.0', 'Plany patroli',      'Dodano moduł planowania patroli',   'Patrol plans, przypisywanie pinów', now() - interval '14 days'),
  ('1.2.0', 'Statystyki',         'Moduł statystyk streetworkerów',    'Wykresy, eksport, filtry miesięczne', now() - interval '3 days');

-- Users (password is "changeme" as SHA-256, same as your Next.js seed)
INSERT INTO users (email, password_hash, must_change_password) VALUES
  ('anna@example.com',   '057ba03d6c44104863dc7361fe4578965d1887360f90a0895882e58a6248fc86', true),
  ('marek@example.com',  '057ba03d6c44104863dc7361fe4578965d1887360f90a0895882e58a6248fc86', true),
  ('admin@example.com',  '057ba03d6c44104863dc7361fe4578965d1887360f90a0895882e58a6248fc86', false);