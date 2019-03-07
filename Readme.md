# grpc_proxy

支持通过`postman` 调用grpc接口

## 使用方式

1. 在postman上设置代理 localhost:7001
2. 发送JSON请求

使用参考： https://github.com/jnewmano/grpc-json-proxy

```bash
curl -X POST \
  http://{{host}}/{{service}}/{{method}} \
  -H 'Content-Type: application/json' \
  -H 'Postman-Token: 135f874d-d5c2-49c9-94e0-5300e6a424dc' \
  -H 'cache-control: no-cache' \
  -H 'proxy-grpc: true' \
  -d '{
	  json payload
}'

```

## 实现

对`grpcurl`进行了一层封装, server端调用了`grpcurl`


