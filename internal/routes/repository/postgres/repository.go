package routes_postgres_repository

import core_postgres_pool "github.com/SALutHere/onelia/internal/core/repository/postgres/pool"

type RoutesRepository struct {
	pool core_postgres_pool.Pool
}

func NewRoutesRepository(pool core_postgres_pool.Pool) *RoutesRepository {
	return &RoutesRepository{
		pool: pool,
	}
}
