namespace php rpc_thrift.services

exception RpcException {
  1: i32  code,
  2: string msg
}

service RpcServiceBase {
    void ping();
}