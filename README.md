```shell
# 初始化 db
nga_grep migrate --db="./nga.db"

# 仅运行同步爬取（兼容旧入口，内部等价于 api-server --mode=sync）
nga_grep sync --cid="" --uid="" --db="./nga.db"

# 仅运行 HTTP API（默认模式 http）
nga_grep api-server --db="./nga.db" --cors=localhost,dashidai.yikakia.com --port=":11648"

# 同时运行 HTTP API + 同步爬取
nga_grep api-server --db="./nga.db" --mode=http,sync --cid="" --uid=""
```

用 git submodule 嵌入了前端的项目，之后可以考虑一起发布




要使用 submodule, 拉取后执行一下命令即可初始化，分为两种情况

- 之前未拉取本仓库

执行以下命令即可拉取本仓库以及子模块
```shell
git clone --recursive git@github.com:yikakia/nga_grep.git
```

- 之前已拉取本仓库

执行以下命令即可初始化子模块

## Docker 使用

项目提供了 `Dockerfile` 与 `docker-compose.yml`，方便通过容器启动两部分服务。

1. **准备数据库文件**

SQLite 数据文件由用户负责提供并挂载到容器，示例：
```sh
mkdir -p data
touch data/nga.db     # 或拷贝已有文件
```

2. **构建镜像**（可选）
如果需要测试修改的话，就在本地构建镜像
```sh
# db 采用挂载的形式 由变量 NGA_DATA_DIR 决定
# export NGA_DATA_DIR=./data/

docker build . -t ghcr.io/yikakia/nga_grep
```

3. **启动全部服务**

```sh
# 拉最新镜像 需要更新时执行
docker compose pull
# 重启容器
docker compose up -d
```

- `sync` 容器会执行数据爬取命令并将结果写入 `/data/nga.db`。
- `api` 容器在 `11648` 端口提供 HTTP API。

4. **查看日志/停止**
```sh
docker compose logs --tail=20 -f      # 查看两个服务的输出
docker compose down                   # 停止并移除容器
```

5. 配置相关
涉及配置为 .api.env .sync.env .env 三个文件，可以通过复制对应的 sample 文件后进行修改

```sh
cp .api.env.sample .api.env
cp .sync.env.sample .sync.sample
cp .env.sample .env
```

.env 中 主要就是
- 挂载 sqlite 文件的地址 
- 是否开启 gin 的 debug 模式的配置 
- 以及 otel 相关配置 servicename endpoint protocol headers 等

.api.env .sync.env 两个文件，也是 otel 相关的配置，主要就是启动的文件的服务名称和服务的环境

> 注意：如果需要调整启动参数，可直接修改 `docker-compose.yml` 中的 `command` 字段，
> 或使用环境变量/扩展的 Compose 配置覆盖。