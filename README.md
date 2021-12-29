# pt-backend

Parcel tracking service.


# TODO: ...

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
