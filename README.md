# microsvc-base

## Dependencies
* Go 1.11
* Go Modules turned on

## Developement
**If you haven't turned them on yet, turn on go modules**
```
GO111MODULE=on
```
**Fetch repo**
```
$ go get github.com/hathbanger/microsvc-base
```

**Run it**
```
$ cd $GOPATH/src/github.com/hathbanger/microsvc-base/cmd
$ go run microsvc-base
```

_______________________________________________
## Dependency management - Go Modules

**Standard commands like go build or go test will automatically add new dependencies as needed to satisfy imports (updating go.mod and downloading the new dependencies).**
```
# After adding package to a source file
$ go build
$ go test
```

**When needed, more specific versions of dependencies can be chosen with commands such as:**
```
$ go get foo@v1.2.3 ||  go get foo@master || go get foo@e3702bed2
```

**You can even update all direct and indirect dependencies to latest minor or patch upgrades**
```
$ go get -u || go get -u=patch
```
