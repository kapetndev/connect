module github.com/crumbandbase/recovery-grpc

go 1.16

require (
	github.com/crumbandbase/service-core-go v0.0.0
	github.com/golang/protobuf v1.4.3 // indirect
	golang.org/x/net v0.7.0 // indirect
	google.golang.org/grpc v1.36.0
)

replace github.com/crumbandbase/service-core-go => ../../
