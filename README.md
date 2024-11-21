# ezcat 🐱 easy reverse shell handler
https://github.com/ngn13/ezcat/assets/78868991/75c3c7c5-6768-47e4-9ef1-0a9e66710773

---

### 📋 Features
- Easy to install
- Simple web interface
- Agent communication over TCP
- Receive TCP reverse shells
- Linux & Windows support
- Self deletion because why not

### 🚀 Install
You can easily install ezcat with docker:
```
docker run --rm --network host       \
    -e EZCAT_PASSWORD=securepassword \
    ghcr.io/ngn13/ezcat
```

### ⚙️ Configuration
Configuration is handled with environment variables, here are all the options:

- **`EZCAT_PASSWORD`**: Used to change the login password, by default it's `ezcat`, and for security, you should
definitely change it
- **`EZCAT_SHELLIP`**: By default ezcat will try to detect your interface IP address (giving priority to tunnel interfaces).
If you want set this IP address to something else by default, you can use the `SHELLIP` environment variable
- **`EZCAT_DISABLE_MEGAMIND`**: When set to `1`, it disables the "no shells?" megamind meme that's displayed on the dashboard if you don't have
any active shells
- **`EZCAT_HTTP_PORT`**: Used to change the port that the API server will listen on, default is 5566
- **`EZCAT_AGENT_PORT`**: Used to change the agent communication port, default is 1053
- **`EZCAT_API_URL`**: Used to change the API URL for the front-end application
- **`EZCAT_DIST_DIR`**: Directory that the server will use to store payload/stage builds, default is `./data`
- **`EZCAT_STATIC_DIR`**: Used to change the front-end application (static) directory, it's pre-set in the Dockerfile,
you don't need to worry about it unless you are working on something
- **`EZCAT_PAYLOAD_DIR`**: Specifies the directory that contains the payloads, it's pre-set in the Dockerfile, just like
the `EZCAT_STATIC_DIR` option, don't worry about it
- **`EZCAT_DEBUG`**: When set to `1`, it enables debug output for the server and the stage builds

### ⚒️ Build
To build ezcat, install a recent version go. Then download and [extract the latest release](http://github.com/ngn13/ezcat/releases/latest).
- To build the server, install a recent version of go, change directory into the `server/` directory and run:
```bash
go build
```

- To build the front-end application, install a recent version of node and npm change directory into the `app/` directory and run:
```bash
npm i
npm run build
```

To build different payloads during runtime, you will need GNU `coreutils` and `bash`, `build-essential` tools and optionally `mingw`
for windows builds. After installing these tools, you can run the `server/` binary with the desired configuration.
