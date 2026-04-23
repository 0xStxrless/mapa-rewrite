-- queries/patrol_plans.sql

-- name: GetPatrolPlan :one
SELECT * FROM patrol_plans WHERE id = $1;

-- name: ListPatrolPlans :many
SELECT * FROM patrol_plans ORDER BY created_at DESC;

-- name: CreatePatrolPlan :one
INSERT INTO patrol_plans (name, date)
VALUES ($1, $2)
RETURNING *;

-- name: DeletePatrolPlan :exec
DELETE FROM patrol_plans WHERE id = $1;

-- name: GetPatrolPlanWithPins :many
SELECT
  ppp.sort_order,
  p.*
FROM patrol_plan_pins ppp
JOIN pins p ON p.id = ppp.pin_id
WHERE ppp.patrol_plan_id = $1
ORDER BY ppp.sort_order ASC;

-- name: AddPinToPatrolPlan :one
INSERT INTO patrol_plan_pins (patrol_plan_id, pin_id, sort_order)
VALUES ($1, $2, $3)
RETURNING *;

-- name: RemovePinFromPatrolPlan :exec
DELETE FROM patrol_plan_pins WHERE patrol_plan_id = $1 AND pin_id = $2;


-- queries/visits.sql

-- name: GetVisitsByPin :many
SELECT * FROM visits WHERE pin_id = $1 ORDER BY visited_at DESC;

-- name: CreateVisit :one
INSERT INTO visits (pin_id, name, note, image_url)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: DeleteVisit :exec
DELETE FROM visits WHERE id = $1;


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

-- name: MarkUpdateViewed :one
INSERT INTO user_updates_viewed (user_id, update_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING
RETURNING *;

-- name: GetUnviewedUpdates :many
SELECT au.*
FROM app_updates au
WHERE au.id NOT IN (
  SELECT update_id FROM user_updates_viewed WHERE user_id = $1
)
ORDER BY au.released_at DESC;


-- queries/streetwork_stats.sql

-- name: GetStatsByWorker :many
SELECT * FROM streetwork_stats WHERE worker_name = $1 ORDER BY month DESC;

-- name: GetStatsByMonth :many
SELECT * FROM streetwork_stats WHERE month = $1 ORDER BY worker_name;

-- name: UpsertStreetworkStat :one
INSERT INTO streetwork_stats (worker_name, month, interactions, new_contacts, interventions)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (worker_name, month) DO UPDATE
SET interactions = $3, new_contacts = $4, interventions = $5, updated_at = now()
RETURNING *;


-- queries/categories.sql

-- name: ListCategories :many
SELECT * FROM categories ORDER BY name;

-- name: CreateCategory :one
INSERT INTO categories (name, color) VALUES ($1, $2) RETURNING *;

-- name: DeleteCategory :exec
DELETE FROM categories WHERE name = $1;