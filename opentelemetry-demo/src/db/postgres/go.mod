module github.com/MichaelRobotics/Kubernetes/opentelemetry-demo/src/db/postgres

go 1.19

require (
	github.com/DATA-DOG/go-sqlmock v1.5.2
	github.com/lib/pq v1.10.9
	github.com/stretchr/testify v1.10.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/MichaelRobotics/Kubernetes/opentelemetry-demo/src/db => ../
