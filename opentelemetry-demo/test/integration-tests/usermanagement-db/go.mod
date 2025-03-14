module github.com/MichaelRobotics/Kubernetes/opentelemetry-demo/test/integration-tests/usermanagement-db

go 1.22.1

require google.golang.org/grpc v1.71.0

require (
	github.com/google/go-cmp v0.7.0 // indirect
	go.opentelemetry.io/otel v1.35.0 // indirect
)

require (
	github.com/MichaelRobotics/Kubernetes/opentelemetry-demo/src/usermanagementservice/genproto/oteldemo v0.0.0-00010101000000-000000000000
	golang.org/x/net v0.35.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250218202821-56aae31c358a // indirect
	google.golang.org/protobuf v1.36.5 // indirect
)

replace github.com/MichaelRobotics/Kubernetes/opentelemetry-demo/src/usermanagementservice/genproto/oteldemo => ../../../src/usermanagementservice/genproto/oteldemo

// Add additional replacements for needed modules
replace github.com/MichaelRobotics/Kubernetes/opentelemetry-demo/src/usermanagementservice => ../../../src/usermanagementservice
replace github.com/MichaelRobotics/Kubernetes/opentelemetry-demo/src/usermanagementservice/handlers => ../../../src/usermanagementservice/handlers
replace github.com/MichaelRobotics/Kubernetes/opentelemetry-demo/src/usermanagementservice/tests/mocks => ../../../src/usermanagementservice/tests/mocks
