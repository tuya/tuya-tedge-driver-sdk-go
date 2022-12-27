### pb文件生成

proto文件所在的目录下执行，对应的pb.go文件会生成在上级目录下

```shell
cd internal/proto/pb
protoc -I=. --go_out=plugins=grpc,paths=source_relative:../ *.proto 
```

### protoc 工具安装

```shell
go get -u github.com/golang/protobuf/protoc-gen-go
```