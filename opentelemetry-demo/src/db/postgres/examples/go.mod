module example

go 1.22.1

replace github.com/MichaelRobotics/Kubernetes/opentelemetry-demo/src/db/postgres => ../

require github.com/MichaelRobotics/Kubernetes/opentelemetry-demo/src/db/postgres v0.0.0-00010101000000-000000000000

require github.com/lib/pq v1.10.9 // indirect
