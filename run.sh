export CGO_CFLAGS=`pkg-config luajit --cflags`
export CGO_LDFLAGS=`pkg-config luajit --libs-only-L`
go run -tags luajit server.go
