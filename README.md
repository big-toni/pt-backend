
# pt-server

Parcel tracking server deployed to Heroku.

## Running Locally

Make sure you have [Go](http://golang.org/doc/install) version 1.12 or newer and the [Heroku Toolbelt](https://toolbelt.heroku.com/) installed.

```sh
$ git clone https://github.com/heroku/go-getting-started.git
$ cd go-getting-started
$ go build -o bin/go-getting-started -v . # or `go build -o bin/go-getting-started.exe -v .` in git bash
github.com/mattn/go-colorable
gopkg.in/bluesuncorp/validator.v5
golang.org/x/net/context
github.com/heroku/x/hmetrics
github.com/gin-gonic/gin/render
github.com/manucorporat/sse
github.com/heroku/x/hmetrics/onload
github.com/gin-gonic/gin/binding
github.com/gin-gonic/gin
github.com/heroku/go-getting-started
$ heroku local
```

Your app should now be running on [localhost:5000](http://localhost:5000/).

## Deploying to Heroku

```sh
$ heroku create
$ git push heroku master
$ heroku open
```

or

[![Deploy](https://www.herokucdn.com/deploy/button.png)](https://heroku.com/deploy)


## Documentation

For more information about using Go on Heroku, see these Dev Center articles:

- [Go on Heroku](https://devcenter.heroku.com/categories/go)


## Alternative Debug with VSCode and Delve

1. Build code with: go build -o $GOPATH/bin/devmarks -i
2. Run server with: devmarks serve
3. In another terminal run delve in headless mode: dlv attach --headless --listen=:2345 --api-version=2 $(pgrep devmarks)
4. Use configuration:
        {
            "name": "Connect to Delve",
            "type": "go",
            "request": "attach",
            "remotePath": "${workspaceRoot}",
            "mode": "remote",
            "port": 2345,
            "host": "127.0.0.1",
            "apiVersion": 2,
            "showLog": true,
            "trace": "verbose"
        }