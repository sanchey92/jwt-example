package pg

const (
	createUser = `INSERT INTO users (id, email, password, role, created_at, updated_at)
                  VALUES ($1, $2, $3, $4, $5, $6)`

	findByEmail = `SELECT id, email, password, role, created_at, updated_at
                   FROM users
                   WHERE email = $1`

	findById = `SELECT id, email, password, role, created_at, updated_at
                FROM users
                WHERE id = $1`
)

const (
	saveToken = `INSERT INTO refresh_tokens (id, user_id, token, expires_at)
                 VALUES ($1, $2, $3, $4)`

	getToken = `SELECT id, user_id, token, expires_at
                 FROM refresh_tokens
                 WHERE token = $1`

	deleteToken = `DELETE FROM refresh_tokens
                   WHERE token = $1`
)
