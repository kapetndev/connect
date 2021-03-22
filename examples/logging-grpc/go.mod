module github.com/crumbandbase/logging-grpc

go 1.16

require (
	github.com/crumbandbase/service-core-go v0.0.0
	github.com/crumbandbase/service-core-go/logging/logrus v0.0.0
	github.com/sirupsen/logrus v1.8.0
	go.opencensus.io v0.22.6
	google.golang.org/grpc v1.36.0
)

replace github.com/crumbandbase/service-core-go => ../../

replace github.com/crumbandbase/service-core-go/logging/logrus => ../../logging/logrus
