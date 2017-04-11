#!/usr/bin/env bash
#
# Generate all etcd protobuf bindings.
# Run from repository root.
set -e
set -u

if ! [[ "$0" =~ "scripts/genproto.sh" ]]; then
	echo "must be run from repository root"
	exit 255
fi

if ! [[ $(protoc --version) =~ "3.2.0" ]]; then
	echo "could not find protoc 3.2.0, is it installed + in PATH?"
	exit 255
fi

GOGOPROTO_ROOT="${GOPATH}/src/github.com/gogo/protobuf"
GOGOPROTO_PATH="${GOGOPROTO_ROOT}:${GOGOPROTO_ROOT}/protobuf"
GRPC_GATEWAY_ROOT="${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway"
PROMETHEUS_ROOT="${GOPATH}/src/github.com/prometheus"

DIRS="nflog/nflogpb silence/silencepb/ api/grpcapi/apipb/"

for dir in ${DIRS}; do
	pushd ${dir}
		protoc --go_out=plugins=grpc:. -I=. \
            -I="${PROMETHEUS_ROOT}" \
            -I="${GRPC_GATEWAY_ROOT}/third_party/googleapis" \
            *.proto
	popd
done

protoc -I. \
    -I="${GOGOPROTO_PATH}" \
    -I="${PROMETHEUS_ROOT}" \
    -I="${GRPC_GATEWAY_ROOT}/third_party/googleapis" \
    --grpc-gateway_out=logtostderr=true:. \
    --swagger_out=logtostderr=true:./Documentation/dev/apispec/swagger/. \
    api/grpcapi/apipb/api.proto

mv Documentation/dev/apispec/swagger/api/grpcapi/apipb/api.swagger.json Documentation/dev/apispec/swagger
rm -rf Documentation/dev/apispec/swagger/api