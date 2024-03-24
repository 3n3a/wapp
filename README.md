# wapp - web application part pot

> i forgot about the exact naming lol

but why does this exist.

well i wanted to build a front-/backend web framework.

## goals

* component architecture
* server side rendered components
* frontend html injection with htmx
* option to choose how output is rendered
	* html - for displaying with htmx
	* json - for building apis
	* xml - idk for lols
* each module should be able to:
	* add a parent menu
	* have submodules
	* have a page
	* have data endpoint
	* have validation options for endpoint
	* have data transformation options

## build your binary

build as static binary, stripped (-s & -w)

```bash
go build -ldflags "-s -w -extldflags=-static" -buildvcs=false .
```

(optional) but reduces binary size dramatically

then compress with upx (will be detected as a virus in some cases)
and then uncompressed at runtime

```bash
upx --brute <name>
```

### benchmark for binary size

script to get sizes in table below

```bash
go build -o test-app-normal .
go build -ldflags "-extldflags=-static" -buildvcs=false -o test-app-static .
go build -ldflags "-s -w -extldflags=-static" -buildvcs=false -o test-app-static-stripped .
go build -ldflags "-s -w -extldflags=-static" -buildvcs=false -o test-app-static-stripped-upx . && upx --brute test-app-static-stripped-upx
```

with test-app

| methods | size |
| --- | --- |
| normal go build | 15M |
| normal go build, static | 18M |
| normal go build, static, stripped | 13M |
| normal go build, static, stripped, packed with upx | 3.7M |
