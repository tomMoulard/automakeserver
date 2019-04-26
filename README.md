# automakeserver

Go HTTP server to catch POST request from GitHub to launch a bash script when
the conditions are good

```bash
$ make
$ ./server
```

Now you can send POST request to `localhost:8080`

### Fix

If you get an error like:
```bash
$ make
server.go:6:5: cannot find package ...
```

Try:
```
$ go get
$ make
```

## TODO

 - [ ] use [webhook repository](https://github.com/go-playground/webhooks) for
   a proper / bigger implementation of the webhook management.
