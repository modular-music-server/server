CGO_CFLAGS=`pkg-config luajit --cflags`
CGO_LDFLAGS=`pkg-config luajit --libs-only-L`
go run -tags luajit server.go
