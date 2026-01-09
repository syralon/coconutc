grpc:
  name: example/grpc
  listen_on: 0.0.0.0:9000
  timeout: 5s

gateway:
  name: example/gateway
  listen_on: 0.0.0.0:8000
  timeout: 5s

etcd:
  endpoints:
    - 127.0.0.1:32379

database:
    driver: mysql
    dsn: root:root@tcp(127.0.0.1:3306)/example?parseTime=true&loc=Local