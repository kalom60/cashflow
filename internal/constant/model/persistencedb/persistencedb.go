package persistencedb

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/kalom60/cashflow/internal/constant/model/db"
	"github.com/kalom60/cashflow/platform/logger"
)

type PersistenceDB struct {
	*db.Queries
	Pool *pgxpool.Pool
	log  logger.Logger
}

type Sibling string

func New(pool *pgxpool.Pool, log logger.Logger) PersistenceDB {
	return PersistenceDB{
		Queries: db.New(pool),
		Pool:    pool,
		log:     log,
	}
}
