# Easy coding

此项目是一个golang结构的实例项目，旨在解决项目在工程方面不标准的情况

## 项目架构

此项目是根据[Uncle Bob's Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)来设计的

## 设计原则

[Single source of truth (SSOT)](https://en.wikipedia.org/wiki/Single_source_of_truth)


<img src="pics/project_architecture.jpg" style="zoom:50%;" />

## 功能

- 所有接口都是由protobuf定义
- 自动生成grpc，grpc-gateway，validate文件
- 每个接口同时提供rest和grpc访问接口
- 自动生成swagger ui文档
- 内置基础的prometheus指标
- 支持将接口导入到postman中进行调试
- 在docker中运行
- 配置管理，配置生成
- 数据库表结构的升级降级
- mock数据库进行单元测试
- golang，protobuf等文件的静态检测与自动修复
- error的分类与管理
- 使用拦截器自动输出日志
- 单元测试和测试覆盖率
- 优雅停止
- 支持启动后台进程
- 健康检查

## 运行前的依赖

- [protoc](https://github.com/protocolbuffers/protobuf#protocol-compiler-installation)

- protoc plugins, go, grpc, grpc-gateway, openapi, validate

``` bash
go install \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest \
    google.golang.org/protobuf/cmd/protoc-gen-go@latest \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest \
    github.com/envoyproxy/protoc-gen-validate@latest
```

- golang 1.18+

- [protobuf management](https://docs.buf.build/installation)

- [go swagger cli](https://github.com/go-swagger/go-swagger/releases)

- docker and docker compose

- (可选择不安装) pre-commit

``` bash
pip3 install pre-commit
pre-commit install
```

- (可选择不安装) golang lint

``` bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

## 运行程序

**如果下载go的依赖包过慢，更新dockerfile中的`GOPROXY`**

``` bash
make deps
make run
```

**注意: 第一次运行需要手动创建 `test` 数据库，否则会出现以下错误**

``` text
failed to initialize database, got error [driver: bad connection]
```

``` bash
docker exec -it easycoding-mysql-1 bash
mysql -u root -p123456
create database test;
```

执行后会生成以下文件

- api/{module_name}/{module_name}.pb.go
- api/{module_name}/{module_name}.pb.validate.go
- api/{module_name}/{module_name}.pb.swagger.json
- api/{module_name}/rpc_grpc.pb.go
- api/{module_name}/rpc.pb.go
- api/{module_name}/rpc.pb.gw.go
- api/{module_name}/rpc.pb.validate.go
- api/{module_name}/rpc.swagger.json

服务监听了三个端口

- 10000: rest api server
- 10001: grpc api server
- 10002: swagger api and prometheus server

检查rest接口

``` bash
curl http://localhost:10000/ping
```

检查grpc接口

``` bash
go run cmd/client/main.go
```

使用浏览器打开以下链接

- http://localhost:10002/swagger/
- http://localhost:10002/metrics

<img src="pics/swagger.png" style="zoom:50%;" />
<img src="pics/metrics.png" style="zoom:50%;" />

### 专题1 数据库升级与降级

#### 动机

在代码中写sql来对数据库进行操作是很难维护的，通常是使用ORM作为中间层来进行交互，本项目使用[Gorm](https://github.com/go-gorm/gorm)，另一个问题是随着业务的迭代，我们经常会更新数据库结构，这是必然的，但是在很多公司，写代码和部署代码不是同一个人，部署代码的人只能僵硬的去执行写代码的人给的sql升级文件，如果一旦出现错误，无法判断如何解决，而且经常业务会有升级和回滚，一方面代码本身要使用【双写】等等的机制来保证数据表的向前兼容，但是表结构本身也要支持版本管理，那就必然需要回滚。另一方面，我们如果有了升级的文件就可以比较容易的结合CI、CD流程集成测试，自动升级。最后一个问题是我们手动的写升级和降级的sql文件是非常困难的，也比较难以维护，所以尽量可以做到自动生成。

现在 `test` 数据库完全是空的，使用以下命令来创建数据库初始化sql文件

``` bash
make migrate-create
```

执行成功后会生成以下文件, 如果想了解为什么会是这种结构，可以查看[migrate](https://github.com/golang-migrate/migrate)

``` text
migrations/pet/{timestamp}_pet.up.sql
migrations/pet/{timestamp}_pet.down.sql
```

通常在云原生的场景下，升级数据库结构，通常是要启动一个[kubernetes job](https://kubernetes.io/docs/concepts/workloads/controllers/job/)，所以这个命令没有和makefile结合

``` bash
go run cmd/migrate/main.go step --latest
```

``` text
INFO[0000] Start buffering 20220723144816/u pet         
INFO[0000] Read and execute 20220723144816/u pet        
INFO[0000] Finished 20220723144816/u pet (read 5.465976ms, ran 57.983119ms)
```

升级成功，使用 `describe pet` 查看表结构

``` text
+------------+----------+------+-----+-------------------+-------------------+
| Field      | Type     | Null | Key | Default           | Extra             |
+------------+----------+------+-----+-------------------+-------------------+
| id         | int      | YES  |     | NULL              |                   |
| name       | text     | YES  |     | NULL              |                   |
| type       | int      | YES  |     | NULL              |                   |
| created_at | datetime | YES  |     | CURRENT_TIMESTAMP | DEFAULT_GENERATED |
+------------+----------+------+-----+-------------------+-------------------+
4 rows in set (0.01 sec)
```

更新 pkg/orm/pet.go

``` text
--- a/pkg/orm/pet.go
+++ b/pkg/orm/pet.go
@@ -12,6 +12,7 @@ type Pet struct {
        Name string
        // TODO(qujiabao): replace int32 to pet_pb.PetType, because of `sqlize`
        Type      int32
+       Age       int32
        CreatedAt time.Time `gorm:"default:now()"`
 }
```

再次生成数据库升级和降级文件，又有两个文件生成了，这时候 migrations/pet下面会有四个文件

``` bash
make migrate-create
```

升级

``` bash
go run cmd/migrate/main.go step --latest
```

查看数据库的版本

``` bash
go run cmd/migrate/main.go version
```

``` text
Version: 20220723150428, Dirty: false
```

``` text
+------------+----------+------+-----+-------------------+-------------------+
| Field      | Type     | Null | Key | Default           | Extra             |
+------------+----------+------+-----+-------------------+-------------------+
| id         | int      | YES  |     | NULL              |                   |
| name       | text     | YES  |     | NULL              |                   |
| type       | int      | YES  |     | NULL              |                   |
| age        | int      | YES  |     | NULL              |                   |
| created_at | datetime | YES  |     | CURRENT_TIMESTAMP | DEFAULT_GENERATED |
+------------+----------+------+-----+-------------------+-------------------+
5 rows in set (0.00 sec)
```

数据库降级

``` bash
go run cmd/migrate/main.go step 1 --reverse
```

``` text
Version: 20220723144816, Dirty: false

+------------+----------+------+-----+-------------------+-------------------+
| Field      | Type     | Null | Key | Default           | Extra             |
+------------+----------+------+-----+-------------------+-------------------+
| id         | int      | YES  |     | NULL              |                   |
| name       | text     | YES  |     | NULL              |                   |
| type       | int      | YES  |     | NULL              |                   |
| created_at | datetime | YES  |     | CURRENT_TIMESTAMP | DEFAULT_GENERATED |
+------------+----------+------+-----+-------------------+-------------------+
4 rows in set (0.00 sec)
```

### 专题2 接口管理

#### 动机

接口是一个抽象的概念，是用来描述两个服务沟通的方式，在真实世界中，有很多文件来描述同一个接口

- golang、java中的struct或class
- typescript、javascript中的class
- 人可以读的文档
- 等等

这个情况违背了[`Single source of
truth`](https://en.wikipedia.org/wiki/Single_source_of_truth)原则，所以接口必须要在一个地方定义，多个地方生成。

## TODO

- Use reflect in configration
- Benchmark
- Fix linting
- Intergration test
- Auth
- More options in configuration
- Property based test
- GraphQL server

## 灵感

- https://github.com/OFFLINE-GmbH/go-webapp-example
- https://github.com/golang-standards/project-layout