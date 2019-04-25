# automakeserver

Go HTTP server to catch POST request from GitHub to launch a bash script when
the conditions are good

```bash
$ make
$ ./server
```

Now you can send POST request to `localhost:8080`

### Fix

if you get an error like:
```bash
$ make
server.go:6:5: cannot find package ...
```

Try:
```
$ go get
$ make
```
