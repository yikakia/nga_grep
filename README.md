```shell
# sync 
nga_grep sync --cid="" --uid="" --db="./nga.db" 
# api
nga_grep api-server --cors=localhost,dashidai.yikakia.com --port=";11648" --db="./nga.db"

# 初始化 db
nga_grep migrate --db="./nga.db"
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
2. **构建镜像**
```sh
# db 采用挂载的形式 由变量 NGA_DATA_DIR 决定
# export NGA_DATA_DIR=./data/
docker-compose build
```

3. **启动全部服务**

```sh
docker-compose up -d
```

- `sync` 容器会执行数据爬取命令并将结果写入 `/data/nga.db`。
- `api` 容器在 `11648` 端口提供 HTTP API。

4. **查看日志/停止**
```sh
docker-compose logs -f      # 查看两个服务的输出
docker-compose down         # 停止并移除容器
```

> 注意：如果需要调整启动参数，可直接修改 `docker-compose.yml` 中的 `command` 字段，
> 或使用环境变量/扩展的 Compose 配置覆盖。