-- name: InsertBadgeQuery :exec
INSERT INTO badge_dispatch (ghUsername, badge_name) 
VALUES ($1, $2);
