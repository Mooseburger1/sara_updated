WORKING_DIR="/home/scott/repos/sara_updated"

protoc -I ${WORKING_DIR}/backend/grpc/proto ${WORKING_DIR}/backend/grpc/proto/*.proto --go_out=plugins=grpc:../