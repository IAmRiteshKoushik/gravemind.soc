-- name: InsertBadgeQuery :exec
INSERT INTO badge_dispatch (ghUsername, badge_name) 
VALUES ($1, $2);

-- name: FindUserStackQuery :one 
SELECT stack from user_account
WHERE ghUsername = $1;

-- name: UpdateUserStackQuery :exec
UPDATE user_account
SET stack = array_cat(COALESCE(stack, ARRAY[]::TEXT[]), $1)
WHERE ghUsername = $2;
