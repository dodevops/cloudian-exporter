package internal

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"time"
)

type metric struct {
	size *prometheus.GaugeVec
}

type CloudianExporter struct {
	api         CloudianAPI
	refresh     time.Duration
	exitChannel chan bool
	metric      *metric
}

func NewCloudianExporter(refresh time.Duration, api CloudianAPI, exitChannel chan bool, registry *prometheus.Registry) CloudianExporter {
	m := &metric{
		size: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "cloudian_bucket_size",
			Help: "Size of bucket in bytes",
		}, []string{"group_id", "user_id", "bucket"}),
	}
	registry.MustRegister(m.size)
	return CloudianExporter{
		api:         api,
		refresh:     refresh,
		exitChannel: exitChannel,
		metric:      m,
	}
}

func (e CloudianExporter) Run() {
	go func() {
		for {
			select {
			case <-e.exitChannel:
				return
			default:
				logrus.Debug("Fetching metrics")
				if r, err := e.api.GetGroups(); err != nil {
					logrus.Errorf("Error fetching groups: %v", err)
				} else {
					for _, groupID := range r {
						if b, err := e.api.GetBuckets(groupID); err != nil {
							logrus.Errorf("Error fetching buckets for group %s: %v", groupID, err)
						} else {
							for _, userBucket := range b {
								if s, err := e.api.GetBucketSize(groupID, userBucket.UserID, userBucket.Bucket); err != nil {
									logrus.Errorf("Error fetching bucket size for group %s, user %s, bucket %s: %v", groupID, userBucket.UserID, userBucket.Bucket, err)
								} else {
									e.metric.size.With(prometheus.Labels{
										"group_id": groupID,
										"user_id":  userBucket.UserID,
										"bucket":   userBucket.Bucket,
									}).Set(float64(s))
								}
							}
						}
					}
				}
				time.Sleep(e.refresh)
			}
		}
	}()
}
