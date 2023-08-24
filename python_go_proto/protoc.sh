/usr/local/bin/python3.11 -m grpc_tools.protoc -I=./proto/ --python_out=./python_plugin/ --grpc_python_out=./python_plugin/ ./proto/*.proto

protoc -I="proto" --go_out=plugins=grpc:./go_plugin/ ./proto/*.proto