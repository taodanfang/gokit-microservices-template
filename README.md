# gokit-microservices-template
基于go kit封装的微服务应用模板：http-iris-consul-gateway-reverse-websocket-zaplog-keystone-cors-device-grpc(client-consul)
支持：
1) 协议：http(iris), grpc, websocket(gosf)
2) 服务注册于发现（http, grpc) -> consul
3) http网关：端口代理网关、反向代理、双网关（http, ws)
4) 统一了各类协议的request/response结构， 模拟http的json方式
5）支持zap logger
6) 支持跨域 cors
7) 支持grpc客户端调用封装，便于实现服务之间的rpc调用
8）内置了认证微服务(keystone），支持oauth/jwt认证
9）内置了微服务实例（http_server, grpc_server, ws_server）

详细文档，请参考博客：https://taodanfang.github.io/post/go-kit-microservice/
