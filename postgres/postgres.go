package postgres

import (
	"context"
	"fmt"
	"net/url"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"go.uber.org/fx"

	"pu/config"
)

type (
	DB interface {
		DB() *pgxpool.Pool
	}
	postgres struct {
		lock *sync.RWMutex
		pool *pgxpool.Pool
	}
	pConf struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		Database string `mapstructure:"database"`
	}
)

func New(lc fx.Lifecycle, config config.Config) (db DB, err error) {
	var pCfg *pConf
	if err = config.UnmarshalKey("postgres", &pCfg); err != nil {
		return nil, err
	}

	p := &postgres{
		lock: &sync.RWMutex{},
		pool: nil,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if p.pool, err = p.connect(ctx, getDsn(pCfg)); err != nil {
				return err
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			p.DB().Close()
			return nil
		},
	})

	return p, nil
}

func (p *postgres) DB() *pgxpool.Pool {
	p.lock.RLock()
	defer p.lock.RUnlock()
	return p.pool
}

func (p *postgres) connect(ctx context.Context, connString string) (pool *pgxpool.Pool, err error) {
	if pool, err = pgxpool.New(ctx, connString); err != nil {
		return nil, errors.Wrap(err, "can't create conn")
	}

	if err = pool.Ping(ctx); err != nil {
		return nil, errors.Wrap(err, "can't ping")
	}

	return pool, nil
}

func getDsn(pCfg *pConf) string {
	endpoint := fmt.Sprintf("%s:%d", pCfg.Host, pCfg.Port)
	dsn := &url.URL{
		Scheme: "postgresql",
		Host:   endpoint,
		Path:   pCfg.Database,
	}

	q := dsn.Query()
	q.Add("sslmode", "disable")
	dsn.RawQuery = q.Encode()

	if pCfg.Username == "" {
		return dsn.String()
	}

	if pCfg.Password == "" {
		dsn.User = url.User(pCfg.Username)
		return dsn.String()
	}

	dsn.User = url.UserPassword(pCfg.Username, pCfg.Password)
	return dsn.String()
}
