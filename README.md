## Introduction

dynamic create router and make it works, you can choose different proxy types, `mesh`, `relay` or `static`

## install

```bash
git clone git@github.com:crackeer/go-gateway.git
go build main.go
```
a `.env` file looks like this

```env
ENV="develop"
PORT=9000
CONFIG_DRIVER="file"
ROUTER_DIR="./config/router"
API_DIR="./config/service"
SYNC_INTERVAL=3
LOG_DIR="./log"
DEBUG=false
```

## Config Your APIs

#### API config sample

> ./config/service/tenapi.json
```json
{
    "service_map" : {
        "default" : {
            "host" : "https://tenapi.cn"
        }
    },
    "data_key" : "",
    "api_list" : [
        {
            "name" : "v2_yiyan",
            "path" : "v2/yiyan",
            "method" : "GET",
            "content_type": ""
        }
    ]
}
```

## Add your routers

router config directory is `./config/router`, if you want add a router `your/path/mesh.json`,you can add the config.json file to the path like this:
#### mesh
> ./config/router/your/path/mesh.json
```json
{
    "mode": "mesh",
    "mesh_config": [
        [
            {
                "api": "uomg/api_qq_info",
                "params": {
                    "qq": "@_input_.qq1"
                },
                "as": "qq1",
                "error_exit": true
            },
            {
                "api": "uomg/api_qq_info",
                "params": {
                    "qq": "@_input_.qq2"
                },
                "as": "qq2",
                "error_exit": true
            }
        ]
    ],
    "response": {
        "qq1" : "@_input_.qq1",
        "data" : "@qq1"
    }
}
```

#### relay

```json
{
    "mode" : "relay",
    "relay_api" : "tenapi/v2_yiyan",
    "response" : {
        "simple" : "@v2_yiyan"
    }
}
```

#### static

```json
{
    "mode" : "static",
    "response" : {
        "simple" : "simple"
    }
}
```