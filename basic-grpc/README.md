# BASIC RPC DEMO
*Tutorial from `gRPC` office website*  

Tool in use:  
* GRPC

GRPC Workflow
1. Define `protobuf`
2. Generate `server` and `client` by using`protobuf` compiler
3. Use the API to write our `server` and `client`

---
> This simple demo is about route mapping app that let client gets information about feature on their route.

---
Compile:  
> protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/basic.proto

* Generate `xxxx.pb.go`
  * Contains all the protocol buffer code to populate, serialize and retrieve request and response message type
* Generate `xxx_grpc.pb.go`
  * Including:
    * An interface type(stub) for clients to call with the methods defined in the `service`
    * Ab interface type for server to implement, also with the methods defined in the `service`

