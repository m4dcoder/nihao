# 你好, Metaverse!
Sample REST service written in Go

### Quick Test
Run `make build` and the binary in a terminal to bring up the server.

```
% make build
Building binaries...
Completed building binaries.

% ./builds/nihao 
INFO[0000] Starting the API server on port 6688.        
INFO[0000] HTTP server will listen at :6688.            
INFO[0000] Launch the http server without socket activation. 
INFO[0000] Successfully started the API server on port 6688.
```

In another terminal, get the endpoint http://localhost:6688/hello.

```
% curl http://localhost:6688/hello
{"message":"你好, Metaverse!"}
```

### Build binary and docker image
Run `make image` to build the binary and docker image.

### Run as container
Run `docker run -p 6688:6688 --platform linux/amd64 m4dcoder/nihao:latest -it /bin/ash`
