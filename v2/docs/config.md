## 配置文件说明
> 配置文件是服务启动的关键，提供必备信息

Go-sword服务启动，依赖yaml格式的配置文件，读取数据库等相关字段，完整的配置文件参考如下：
```yaml
db: # 数据库相关配置
  user: root
  password: '123456'
  database: test
  host: localhost
  port: 3306
root_path: admin22   # Go-sword会在项目根目录创建该目录，并且释放generate按钮生成的所有文件到该目录
tool_port: '8081'    # Go-sword服务启动的端口
server_port: '8082'  # Go-sword后端CRUD服务的端口
```


## 注意

Go-sword工具服务 与 Go-sword CRUD服务启动加载的配置文件格式基本一致，可以在一个文件中进行配置