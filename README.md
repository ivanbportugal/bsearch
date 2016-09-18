# BSearch
## A simple search app (server side only) on the greatest text ever.
An Angular 2 app is in [BSearchClient](https://github.com/ivanbportugal/bsearchclient) that makes queries to this one. It is expected to be deployed to the `public` folder.

## Instructions
To run in development, make sure you have `go` and `boltdb` installed.
This will seed your database in the `translations` folder:
```
go run migrate.go
```

This will run the server. Feel free to test using postman or the like:
```
go run main.go
```

To run in production, compile the binaries using the `GOOS` and `GOARCH` system/command line variables:
```
GOOS="linux" GOARCH="amd64" migrate.go
GOOS="linux" GOARCH="amd64" main.go
```

Possible values for the above:
```
const goosList = "android darwin dragonfly freebsd linux nacl \ 
  netbsd openbsd plan9 solaris windows "
const goarchList = "386 amd64 amd64p32 arm arm64 ppc64 ppc64le \
   mips mipsle mips64 mips64le mips64p32 mips64p32le \ # (new)
   ppc s390 s390x sparc sparc64 "
```

After moving to your server, simply run as a process. For example,
```
chmod +x ./migrate
# Seed your BoltDB instance
./migrate

chmod +x ./main
# Should return your console immediately and run in the background.
nohup ./main >/dev/null 2>&1 &
```

To kill the background process, you have to find its PID:
```
ps -ef | grep main
# Remember the PID from this list
kill (PID)