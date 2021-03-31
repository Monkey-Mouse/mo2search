# Mo2 Search Service

[Mo2](https://github.com/Monkey-Mouse/Mo2)的微服务之一，提供搜索相关服务  
> [ATTENTION] 此服务**没有任何安全性**，任何情况下不应该将该服务的接口直接暴露给外界
## API
- 创建索引
  - POST: `/api/index`
  - QUERY: `name`索引名
- 列出索引
  - GET: `/api/index`
- 更新/插入索引的文档
  - PUT: `/api/{index}`
  - QUERY: `id`文档的id
  - BODY: 文档，json格式
- 删除索引的文档
  - DELETE: `/api/{index}`
  - QUERY: `id`文档的id
- 搜索文档
  - POST: `/api/search`
  - QUERY: `index`搜索的索引
  - BODY: `bleve.SearchRequest`
  - Response: `bleve.SearchResult`

## Dokcer
用一下命令拉取我们的docker image
```bash
sudo docker pull registry.cn-hangzhou.aliyuncs.com/mo2/mo2search:latest
```


