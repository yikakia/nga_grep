```shell
# sync 
nga_grep sync --cid="" --uid="" --db="./nga.db" 
# api
nga_grep api-server --cors=localhost,dashidai.yikakia.com --port=";11648" --db="./nga.db"
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

```shell
git submodule update --init --recursive
```