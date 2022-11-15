package postgres

import (
	"context"
	"fmt"
	"net/url"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"pu/config"
	"pu/logger"
)

type (
	DB interface {
		DB() *pgx.Conn
	}
	postgres struct {
		lock *sync.RWMutex
		conn *pgx.Conn
	}
	pConf struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		Database string `mapstructure:"database"`
	}
)

func New(lc fx.Lifecycle, config config.Config, l logger.Logger) (db DB, err error) {
	var pCfg *pConf
	if err = config.UnmarshalKey("postgres", &pCfg); err != nil {
		return nil, err
	}

	p := &postgres{
		lock: &sync.RWMutex{},
		conn: nil,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if p.conn, err = p.connect(ctx, getDsn(pCfg)); err != nil {
				return err
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if errD := p.DB().Close(ctx); errD != nil {
				l.Error("can't close postgres", zap.Error(errD))
			}
			return nil
		},
	})

	return p, nil
}

func (p *postgres) DB() *pgx.Conn {
	p.lock.RLock()
	defer p.lock.RUnlock()
	return p.conn
}

func (p *postgres) connect(ctx context.Context, connString string) (conn *pgx.Conn, err error) {
	if conn, err = pgx.Connect(ctx, connString); err != nil {
		return nil, errors.Wrap(err, "can't create conn")
	}

	if err = conn.Ping(ctx); err != nil {
		return nil, errors.Wrap(err, "can't ping")
	}

	return conn, nil
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
