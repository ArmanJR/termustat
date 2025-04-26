# Misc docs

To make a `digest.txt` of Termustat:

```shell
gitingest ./ -e "frontend/node_modules/*,node_modules,go.sum,LICENSE,.env,*.sample,*.sql,engine/legacy/,/frontend/package-lock.json,.gitignore,/frontend/.gitignore,/engine/temp/,/engine/.gitignore,api/docs/,/docs/swagger.yaml,api/digest.txt,*/digest.txt,digest.txt"
```