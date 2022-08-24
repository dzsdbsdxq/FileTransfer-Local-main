# 文件快传
只需简单1步，极速文件互传，从未这样爽过！
### 在线预览

[示例地址](https://mapi.net.cn/)

### 前期准备

除过 mysql 以外，其他服务都开箱即可使用，无需进行初始化之类的操作

mysql 需要创建数据库`db_file_transfer`，使用 data/tb_files.sql 来初始化表结构和表数据


### 项目结构

```
-FileTransfer-Local-main
    |-bootstrap 插件目录
    |-config 配置文件目录
    |-controller 控制器目录
    |-data  数据库sql存放目录
    |-middlewares 中间件目录
    |-models 数据库模型目录
    |-pkg 第三方包目录
    |-routers 路由目录
    |-runtime 默认下载目录
    |-web web资源目录
        |-css css文件目录
        |-fonts 字体目录
        |-js js文件目录
    |-main.go 程序执行入口
```

### 如何运行？

+ 把代码clone到自己电脑或者服务器上

  >  git clone https://github.com/dzsdbsdxq/FileTransfer-Local-main.git

+ 修改配置`config/app.ini`
  ```golang
    #debug or release
    RUN_MODE = release

    [app]
    # 每秒请求限制次数
    LIMIT_COUNT_PER_SECOND = 7
    # 此配置目前使用不到
    JWT_SECRET = 23347$040412
    #上传路径
    UPLOAD_DIR = /FileTransfer/runtime/upload/
    # 下载url
    DOWN_DOMAIN = /api/d/
    # 网站标题
    WEB_NAME = 点点笔记
    PROXY = false

    [server]
    ADDR_IP = 0.0.0.0
    HTTP_PORT = 8999
    READ_TIMEOUT = 60
    WRITE_TIMEOUT = 60

    [database]
    TYPE = mysql
    USER = root
    PASSWORD = root
    HOST = localhost:3306
    NAME = db_file_transfer
    TABLE_PREFIX = tb_
    ```
    
 + 运行
   ```bash
    go run main.go
    ```
    后台运行
    
    ```bash
    nohup go run main.go &
    ```
    
### nginx反向代理

```bash
server{
  listen 443 ssl;
  server_name you-domain;

  location / {
      proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
      proxy_set_header Host $http_host;
      proxy_set_header X-Real-IP $remote_addr;
      proxy_redirect off;
      proxy_pass http://127.0.0.1:8999;
  }
}
```


    
    


