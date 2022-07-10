# update swagger ui

https://github.com/elazarl/go-bindata-assetfs

``` bash
go get github.com/go-bindata/go-bindata/...
go get github.com/elazarl/go-bindata-assetfs/...
```

git clone https://github.com/swagger-api/swagger-ui

update url in swagger-initializer.js from https://petstore.swagger.io/v2/swagger.json to ./api.swagger.json

go-bindata-assetfs dist/...
