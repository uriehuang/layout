package data

import (
	"layout/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/extra/redisotel"
	redis "github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"github.com/uriehuang/pkg/database"
	"github.com/uriehuang/pkg/metrics"
	"gorm.io/gorm"
	gLogger "gorm.io/gorm/logger"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo)

// Data .
type Data struct {
	XsDB *gorm.DB
	Rds  *redis.Client
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	xsDB, err := database.NewMysqlDB(&database.MysqlConf{
		MasterDsn:       c.GetDatabase().GetXs().GetMaster(),
		SlaveDsn:        c.GetDatabase().GetXs().GetSlaves(),
		MaxIdleConns:    c.GetDatabase().GetXs().GetMaxIdleConns(),
		MaxOpenConns:    c.GetDatabase().GetXs().GetMaxOpenConns(),
		ConnMaxIdleTime: c.GetDatabase().GetXs().GetConnMaxIdleTime().AsDuration(),
		ConnMaxLifetime: c.GetDatabase().GetXs().GetConnMaxLifetime().AsDuration(),
	}, "layout", gLogger.Info)
	if err != nil {
		return nil, nil, err
	}

	rds := redis.NewClient(&redis.Options{
		Network:      c.GetRedis().GetNetwork(),
		Addr:         c.GetRedis().GetAddr(),
		ReadTimeout:  c.GetRedis().GetReadTimeout().AsDuration(),
		WriteTimeout: c.GetRedis().GetWriteTimeout().AsDuration(),
		Password:     c.GetRedis().GetPassword(),
	})
	rds.AddHook(redisotel.TracingHook{})          // 链路追踪
	rds.AddHook(metrics.NewMetricsHook("layout")) // metrics

	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
		_ = database.Close(xsDB)
		_ = rds.Close()
	}

	return &Data{
		XsDB: xsDB,
		Rds:  rds,
	}, cleanup, nil
}
