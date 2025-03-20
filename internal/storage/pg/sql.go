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
