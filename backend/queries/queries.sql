-- queries/patrol_plans.sql

-- name: GetPatrolPlan :one
SELECT * FROM patrol_plans WHERE id = $1;

-- name: ListPatrolPlans :many
SELECT * FROM patrol_plans ORDER BY created_at DESC;

-- name: ListPatrolPlansWithPinCount :many
SELECT
  pp.*,
  COUNT(ppp.pin_id) AS pin_count
FROM patrol_plans pp
LEFT JOIN patrol_plan_pins ppp ON ppp.patrol_plan_id = pp.id
GROUP BY pp.id
ORDER BY pp.created_at DESC;

-- name: ListPatrolPlansByDate :many
SELECT * FROM patrol_plans WHERE date = $1 ORDER BY created_at DESC;

-- name: GetPatrolPlansByDateRange :many
SELECT * FROM patrol_plans
WHERE date BETWEEN $1 AND $2
ORDER BY date DESC;

-- name: GetPatrolPlanWithPins :many
SELECT
  ppp.sort_order,
  p.*
FROM patrol_plan_pins ppp
JOIN pins p ON p.id = ppp.pin_id
WHERE ppp.patrol_plan_id = $1
ORDER BY ppp.sort_order ASC;

-- name: CreatePatrolPlan :one
INSERT INTO patrol_plans (name, date)
VALUES ($1, $2)
RETURNING *;

-- name: AddPinToPatrolPlan :one
INSERT INTO patrol_plan_pins (patrol_plan_id, pin_id, sort_order)
VALUES ($1, $2, $3)
RETURNING *;

-- name: RemovePinFromPatrolPlan :exec
DELETE FROM patrol_plan_pins WHERE patrol_plan_id = $1 AND pin_id = $2;

-- name: DeletePatrolPlan :exec
DELETE FROM patrol_plans WHERE id = $1;


-- queries/pins.sql

-- name: GetPin :one
SELECT * FROM pins WHERE id = $1;

-- name: GetPinWithVisits :many
SELECT
  p.*,
  v.id AS visit_id,
  v.name AS visitor_name,
  v.note,
  v.image_url AS visit_image_url,
  v.visited_at
FROM pins p
LEFT JOIN visits v ON v.pin_id = p.id
WHERE p.id = $1
ORDER BY v.visited_at DESC;

-- name: ListPins :many
SELECT * FROM pins ORDER BY created_at DESC;

-- name: ListPinsWithVisitCount :many
SELECT
  p.*,
  COUNT(v.id) AS visit_count,
  MAX(v.visited_at) AS last_visited_at
FROM pins p
LEFT JOIN visits v ON v.pin_id = p.id
GROUP BY p.id
ORDER BY p.created_at DESC;

-- name: ListPinsWithLastVisit :many
SELECT
  p.*,
  v.name AS last_visitor,
  v.visited_at AS last_visited_at,
  v.note AS last_note
FROM pins p
LEFT JOIN LATERAL (
  SELECT * FROM visits WHERE pin_id = p.id ORDER BY visited_at DESC LIMIT 1
) v ON true
ORDER BY p.created_at DESC;

-- name: ListPinsByCategory :many
SELECT * FROM pins WHERE category = $1 ORDER BY created_at DESC;

-- name: ListPinsByCategoryWithVisitCount :many
SELECT
  p.*,
  COUNT(v.id) AS visit_count,
  MAX(v.visited_at) AS last_visited_at
FROM pins p
LEFT JOIN visits v ON v.pin_id = p.id
WHERE p.category = $1
GROUP BY p.id
ORDER BY p.created_at DESC;

-- name: SearchPins :many
SELECT * FROM pins
WHERE
  title ILIKE '%' || $1 || '%'
  OR description ILIKE '%' || $1 || '%'
ORDER BY created_at DESC;

-- name: GetPinsNearLocation :many
SELECT *,
  round(cast(point(lng, lat) <-> point($1, $2) as numeric), 6) AS distance
FROM pins
ORDER BY point(lng, lat) <-> point($1, $2)
LIMIT $3;

-- name: CreatePin :one
INSERT INTO pins (title, description, lat, lng, category, image_url)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: UpdatePin :one
UPDATE pins
SET title = $2, description = $3, lat = $4, lng = $5,
    category = $6, image_url = $7, updated_at = now(), version = version + 1
WHERE id = $1
RETURNING *;

-- name: IncrementPinVisits :exec
UPDATE pins SET visits_count = COALESCE(visits_count, 0) + 1 WHERE id = $1;

-- name: DeletePin :exec
DELETE FROM pins WHERE id = $1;


-- queries/visits.sql

-- name: GetVisitByID :one
SELECT * FROM visits WHERE id = $1;

-- name: GetVisitsByPin :many
SELECT * FROM visits WHERE pin_id = $1 ORDER BY visited_at DESC;

-- name: GetRecentVisits :many
SELECT
  v.*,
  p.title AS pin_title,
  p.category AS pin_category
FROM visits v
JOIN pins p ON p.id = v.pin_id
ORDER BY v.visited_at DESC
LIMIT $1;

-- name: GetVisitsByDateRange :many
SELECT
  v.*,
  p.title AS pin_title,
  p.category AS pin_category
FROM visits v
JOIN pins p ON p.id = v.pin_id
WHERE v.visited_at BETWEEN $1 AND $2
ORDER BY v.visited_at DESC;

-- name: CreateVisit :one
INSERT INTO visits (pin_id, name, note, image_url)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateVisit :one
UPDATE visits
SET name = $2,
    note = $3,
    image_url = $4,
    pin_id = $5
WHERE id = $1
RETURNING *;

-- name: DeleteVisit :exec
DELETE FROM visits WHERE id = $1;


-- queries/categories.sql

-- name: GetCategory :one
SELECT * FROM categories WHERE name = $1;

-- name: ListCategories :many
SELECT * FROM categories ORDER BY name;

-- name: ListCategoriesWithPinCount :many
SELECT
  c.*,
  COUNT(p.id) AS pin_count
FROM categories c
LEFT JOIN pins p ON p.category = c.name
GROUP BY c.name
ORDER BY pin_count DESC;

-- name: CreateCategory :one
INSERT INTO categories (name, color) VALUES ($1, $2) RETURNING *;

-- name: UpdateCategory :one
UPDATE categories
SET color = $2
WHERE name = $1
RETURNING *;

-- name: DeleteCategory :exec
DELETE FROM categories WHERE name = $1;


-- queries/streetwork_stats.sql

-- name: GetAllStats :many
SELECT * FROM streetwork_stats ORDER BY month DESC, worker_name;

-- name: GetStatsByWorker :many
SELECT * FROM streetwork_stats WHERE worker_name = $1 ORDER BY month DESC;

-- name: GetStatsByMonth :many
SELECT * FROM streetwork_stats WHERE month = $1 ORDER BY worker_name;

-- name: GetStatsByDateRange :many
SELECT * FROM streetwork_stats
WHERE month BETWEEN $1 AND $2
ORDER BY month DESC, worker_name;

-- name: GetStatsSummary :many
SELECT
  worker_name,
  COUNT(*) AS months_recorded,
  SUM(interactions) AS total_interactions,
  SUM(new_contacts) AS total_new_contacts,
  SUM(interventions) AS total_interventions,
  AVG(interactions) AS avg_interactions_per_month,
  MIN(month) AS first_recorded,
  MAX(month) AS last_recorded
FROM streetwork_stats
GROUP BY worker_name
ORDER BY total_interactions DESC;

-- name: GetMonthlyTotals :many
SELECT
  month,
  SUM(interactions) AS total_interactions,
  SUM(new_contacts) AS total_new_contacts,
  SUM(interventions) AS total_interventions,
  COUNT(DISTINCT worker_name) AS active_workers
FROM streetwork_stats
GROUP BY month
ORDER BY month DESC;

-- name: GetTopWorkersByMonth :many
SELECT
  month,
  worker_name,
  interactions,
  RANK() OVER (PARTITION BY month ORDER BY interactions DESC) AS rank
FROM streetwork_stats
ORDER BY month DESC, rank;

-- name: UpsertStreetworkStat :one
INSERT INTO streetwork_stats (worker_name, month, interactions, new_contacts, interventions)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (worker_name, month) DO UPDATE
SET interactions = $3, new_contacts = $4, interventions = $5, updated_at = now()
RETURNING *;


-- queries/users.sql

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: CreateUser :one
INSERT INTO users (email, password_hash)
VALUES ($1, $2)
RETURNING *;

-- name: UpdateLastLogin :exec
UPDATE users SET last_login = now() WHERE id = $1;

-- name: SetPasswordHash :exec
UPDATE users
SET password_hash = $2, must_change_password = false
WHERE id = $1;


-- queries/app_updates.sql

-- name: ListAppUpdates :many
SELECT * FROM app_updates ORDER BY released_at DESC;

-- name: GetLatestAppUpdate :one
SELECT * FROM app_updates ORDER BY released_at DESC LIMIT 1;

-- name: GetUnviewedUpdates :many
SELECT au.*
FROM app_updates au
WHERE au.id NOT IN (
  SELECT update_id FROM user_updates_viewed WHERE user_id = $1
)
ORDER BY au.released_at DESC;

-- name: MarkUpdateViewed :one
INSERT INTO user_updates_viewed (user_id, update_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING
RETURNING *;

-- name: SelectWorkers :many
SELECT DISTINCT worker_name FROM streetwork_stats;