```shell
# sync 
nga_grep sync --cid="" --uid="" --db="./nga.db" 
# api
nga_grep api-server --cors=localhost,dashidai.yikakia.com --port=";11648" --db="./nga.db"
```

用 git submodule 嵌入了前端的项目，之后可以考虑一起发布