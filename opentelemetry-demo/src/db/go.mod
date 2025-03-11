module github.com/MichaelRobotics/Kubernetes/opentelemetry-demo/src/db

go 1.19

require (
	github.com/DATA-DOG/go-sqlmock v1.5.2
	github.com/MichaelRobotics/Kubernetes/opentelemetry-demo/src/db/postgres v0.0.0-00010101000000-000000000000
	github.com/lib/pq v1.10.9
)

replace github.com/MichaelRobotics/Kubernetes/opentelemetry-demo/src/db/postgres => ./postgres
