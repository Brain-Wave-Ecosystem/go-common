package models

import "github.com/jackc/pgx/v5/pgxpool"

type Assemblable interface {
	Assemble(db *pgxpool.Pool) error
}
