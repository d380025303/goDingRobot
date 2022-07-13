# 使用钉钉机器人，实现简单的备忘录

![image](https://user-images.githubusercontent.com/31104430/178475816-6a41146c-f66c-4ab1-b6ea-e42b73f3a8ba.png)

# 配置
```
{
  "ding": {
    "appKey": "",
    "appSecret": "",
    "agentId": ""
  },
  "datasource": {
    "url": ""
  }
}
```
ding.*: 钉钉的配置
url: 数据库地址(Sqlite3)

# Docker
```docker-compose
version: "3"
services:
  go_ding_robot:
    build: .
    restart: always
    ports:
      - "8000:8000"
    environment:
      APP_KEY: "YOUR_APP_KEY"
      APP_SECRET: "YOUR_APP_SECRET"
    volumes:
      - /yourvolume:/usr/data

```