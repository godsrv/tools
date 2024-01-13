package es

import (
	"crypto/tls"
	"fmt"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

var (
	Client *elastic.Client
)

// @author: lipper
// @object: *elastic.Client
// @function: NewEsClient
// @description: 实例elastic
// @return: *elastic.Client
func NewEsClient(conf EsConf) *elastic.Client {
	var (
		err       error
		version   string
		settings  = make([]elastic.ClientOptionFunc, 0)
		endpoints = make([]string, 0)
		okNodes   = make([]string, 0)
		errNodes  = make([]string, 0)
	)

	if len(conf.Endpoints) == 0 {
		panic("must specify elastic endpoints")
	}

	if !conf.Https {
		for idx := range conf.Endpoints {
			endpoints = append(endpoints, fmt.Sprintf("http://%s", conf.Endpoints[idx]))
		}
	} else {
		for idx := range conf.Endpoints {
			endpoints = append(endpoints, fmt.Sprintf("https://%s", conf.Endpoints[idx]))
		}
		settings = append(settings, elastic.SetScheme("https"))
		settings = append(settings, elastic.SetHttpClient(
			&http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{
						InsecureSkipVerify: true,
					},
				},
			},
		))
	}

	for idx := range endpoints {
		settings = append(settings, elastic.SetURL(endpoints[idx]))
	}

	settings = append(settings, elastic.SetSniff(conf.Sniff))

	settings = append(settings, elastic.SetGzip(conf.Gzip))

	if conf.Auth {
		settings = append(settings, elastic.SetBasicAuth(conf.UserPass.Username, conf.UserPass.Password))
	}
	if conf.EnableTrace {
		settings = append(settings, elastic.SetTraceLog(logrus.New().WithField("trace", "elastic_trace")))
	}

	if !conf.EnableHealthCheck {
		settings = append(settings, elastic.SetHealthcheck(conf.EnableHealthCheck))
	} else {
		settings = append(settings, elastic.SetHealthcheckInterval(60*time.Second))
		settings = append(settings, elastic.SetHealthcheckTimeout(5*time.Second))
	}

	Client, err = elastic.NewClient(settings...)
	if err != nil {
		logrus.Panicf("init es client err: %v", err)
	}

	switch len(okNodes) {
	case 0:
		logrus.Panicf("all nodes: %+v unavailable", endpoints)
	case len(conf.Endpoints):
		logrus.Infof("connect to elastic[version: %s] success, all nodes: %+v available", version, endpoints)
	default:
		logrus.Warnf("connect to elastic all nodes: %+v, err nodes: %v", endpoints, errNodes)
	}

	return Client
}
