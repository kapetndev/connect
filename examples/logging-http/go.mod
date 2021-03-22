module github.com/crumbandbase/logging-http

go 1.16

require (
	github.com/crumbandbase/service-core-go v0.0.0
	github.com/crumbandbase/service-core-go/logging/logrus v0.0.0
	github.com/sirupsen/logrus v1.8.1
	go.opencensus.io v0.23.0
)

replace github.com/crumbandbase/service-core-go => ../../

replace github.com/crumbandbase/service-core-go/logging/logrus => ../../logging/logrus
