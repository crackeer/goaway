# Introduction

dynamic create router and make it works, you can choose different proxy types, `mesh`, `relay` or `static`

# install

## 1. go build
```bash
git clone git@github.com:crackeer/go-gateway.git
go build main.go
```

## 2. .env file looks like this

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


# Database SQL

## SQLite

```sql
CREATE TABLE service(
   id INTEGER PRIMARY KEY AUTOINCREMENT,
   name TEXT NOT NULL,
   service TEXT NOT NULL,
   host TEXT NOT NULL,
   env TEXT NOT NULL,
   timeout int,
   sign TEXT NOT NULL,
   sign_config TEXT NOT NULL,
   disable_extract INTEGER NOT NULL,
   code_key TEXT NOT NULL,
   message_key TEXT NOT NULL,
   data_key TEXT NOT NULL,
   success_code TEXT NOT NULL,
   description TEXT NOT NULL,
   create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   modify_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_service_env on service (service, env);

CREATE TABLE service_api(
   id INTEGER PRIMARY KEY AUTOINCREMENT,
   service TEXT NOT NULL,
   api TEXT NOT NULL,
   path TEXT NOT NULL,
   method TEXT NOT NULL,
   content_type TEXT NOT NULL,
   timeout int,
   description TEXT NOT NULL,
   create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   modify_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) CREATE INDEX idx_service_api on service_api (service, api);

CREATE TABLE router (
   id INTEGER PRIMARY KEY AUTOINCREMENT,
   path TEXT NOT NULL,
   status INT NOT NULL,
   classify TEXT NOT NULL,
   mode TEXT NOT NULL,
   request TEXT NOT NULL,
   response TEXT NOT NULL,
   description TEXT NOT NULL,
   create_at INT,
   modify_at INT
);
```
