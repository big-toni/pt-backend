# pt-backend

Parcel tracking app.
Project created for practicing programming skills with Go.

## NOTE:

- using gorilla/mux HTTP request multiplexer request router and dispatcher for matching incoming requests to their respective handler.
- using chromedp to drive browser and scrape web for parcel number info
- using DAO pattern to separate the application/business layer from the persistence layer
- using go.mongodb.org/mongo-driver, MongoDB supported driver for Go

### Alternative Debug with VSCode and Delve

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
