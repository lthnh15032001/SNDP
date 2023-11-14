<h1 align="center">
  Streaming Network Data Platform
</h1>

Streaming Network Data Platform (**SNDP**) is a project which based on chisel - TCP/UDP tunnel, transported over HTTP, secured via SSH, reference [here](https://github.com/jpillora/chisel). Integrate with API, store authen user with mysql or sqlite & interactive user website interface,... IP address restriction to specific destination. 

*Some of project use case*
- Host a website on Emmbed without opening ports on router.
- Fast delivery to public network from your private network (Testing without deploying)
- Running personal services from your home


## 🚀 Quick start


1. **Start coding!**

        make build
        make build-prod

## 🧐 What's inside?

A quick look at the top-level files and directories you'll see in this project.

    .
    ├── build
    ├── cmd
    ├── docs
    ├── internal
    ├── utils
    ├── seeds
    ├── share
    ├── Dockerfile
    ├── docker-compose.yml
    ├── Makefile
    ├── go.mod
    ├── go.sum
    ├── .gitignore
    └── README.md
    └── users.json

1. **`/build`**: Binary after a build
 
2. **`/internal/api`**:  Viper configration and the property files, controllers that the router/handlers are configured to call..

3. **`/internal/client`**:  Internal Client tunnel 
4. **`/internal/server`**:  Internal Server tunnel 
5. **`/share`**:  Share library
6. **`/cmd/stream/client`**:  Root command (api, client, server)
7. **`/cmd/gobind`**:  Export for android aar file


### Deploy with version 
env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath ${LDFLAGS} ${GCFLAGS} ${ASMFLAGS} -o ${DIR}/github.com/lthnh15032001/ngrok-impl-v1 .