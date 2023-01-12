module github.com/krasish/payment-system

go 1.19

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/docker/distribution v2.8.1+incompatible
	github.com/docker/go-connections v0.4.0
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/jackc/pgx/v5 v5.2.0
	github.com/lib/pq v1.10.7
	github.com/onsi/ginkgo/v2 v2.6.1
	github.com/onsi/gomega v1.24.2
	github.com/sirupsen/logrus v1.9.0
	github.com/testcontainers/testcontainers-go v0.17.0
	github.com/vrischmann/envconfig v1.3.0
	gorm.io/driver/postgres v1.4.6
	gorm.io/gorm v1.24.3
)

require (
	github.com/Azure/go-ansiterm v0.0.0-20210617225240-d185dfc1b5a1 // indirect
	github.com/Microsoft/go-winio v0.5.2 // indirect
	github.com/cenkalti/backoff/v4 v4.2.0 // indirect
	github.com/containerd/containerd v1.6.12 // indirect
	github.com/docker/docker v20.10.20+incompatible // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/felixge/httpsnoop v1.0.1 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/klauspost/compress v1.11.13 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/moby/patternmatcher v0.5.0 // indirect
	github.com/moby/sys/sequential v0.5.0 // indirect
	github.com/moby/term v0.0.0-20221128092401-c43b287e0e0f // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.0-rc2 // indirect
	github.com/opencontainers/runc v1.1.3 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/crypto v0.4.0 // indirect
	golang.org/x/net v0.5.0 // indirect
	golang.org/x/sys v0.4.0 // indirect
	golang.org/x/text v0.6.0 // indirect
	google.golang.org/genproto v0.0.0-20220617124728-180714bec0ad // indirect
	google.golang.org/grpc v1.47.0 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

//Needed for the usage of testcontainers-go. See point 2 here: https://golang.testcontainers.org/quickstart/
replace github.com/docker/docker => github.com/docker/docker v20.10.3-0.20221013203545-33ab36d6b304+incompatible // 22.06 branch
