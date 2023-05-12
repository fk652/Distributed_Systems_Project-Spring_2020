module project/version1

require commonpb v1.0.0
replace commonpb v1.0.0 => ./commonpb

go 1.13

require (
	github.com/coreos/etcd v3.3.20+incompatible
	github.com/gin-gonic/gin v1.6.2
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/golang/protobuf v1.4.0-rc.4.0.20200313231945-b860323f09d0
	golang.org/x/net v0.0.0-20190311183353-d8887717615a
	google.golang.org/grpc v1.28.1
	google.golang.org/protobuf v1.21.0
)
