package arc

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/influxdata/influxdb/client/v2"
)

type InfluxAdapter struct {
	client.Client
	client.BatchPoints
}

func NewInfluxAdapter() *InfluxAdapter {
	cfg := BeanFactory.Get(IConfig(nil)).(IConfig).GetInfluxDB()
	username, passwd := cfg.GetAuth()
	cli, err := client.NewHTTPClient(client.HTTPConfig{Addr: cfg.GetUrl(), Username: username, Password: passwd})
	if err != nil {
		zap.L().Error("Init Influx DB Filed.", zap.Error(err))
		return nil
	}
	createDBSql := client.NewQuery(fmt.Sprintf("Create Database %s", cfg.GetDatabase()), "", "")
	if _, err = cli.Query(createDBSql); err != nil {
		panic(err)
	}
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Precision: "s",
		Database:  cfg.GetDatabase(),
	})
	return &InfluxAdapter{cli, bp}
}

func (this *InfluxAdapter) Name() string {
	return "InfluxAdapter"
}
