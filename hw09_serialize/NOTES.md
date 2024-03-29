# Protocol Buffers Documentation

- [Defining Your Protocol Format](https://protobuf.dev/getting-started/gotutorial/#protocol-format)
- [Scalar Value Types](https://protobuf.dev/programming-guides/proto3/#scalar)
- [Go и Protocol Buffers толика практики](https://habr.com/ru/articles/252455/)

```go
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
protoc --go_out=. --go_opt=Mbook.proto=./book *.proto
```
