package postgres

import (
	"database/sql"
	"fmt"

	"github.com/Car-Wash/Car-Wash-Auth-Service/config"
	"github.com/Car-Wash/Car-Wash-Auth-Service/storage"
	"github.com/go-redis/redis/v8"

	_ "github.com/lib/pq"
)

type Storage struct {
	Db    *sql.DB
	Users storage.User
	Auths storage.Auth
	rdb   *redis.Client
}

func NewPostgresStorage(redisDb *redis.Client) (storage.StorageI, error) {
	config := config.Load()
	con := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		config.PostgresUser, config.PostgresPassword,
		config.PostgresHost, config.PostgresPort,
		config.PostgresDatabase)
	db, err := sql.Open("postgres", con)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {

		return nil, err
	}

	return &Storage{Db: db, rdb: redisDb}, err
}

func (s *Storage) User() storage.User {
	if s.Users == nil {
		s.Users = &UserRepo{s.Db}
	}
	return s.Users
}

func (s *Storage) Auth() storage.Auth {
	if s.Auths == nil {
		s.Auths = &AuthRepo{db: s.Db, redis: s.rdb}
	}
	return s.Auths
}
