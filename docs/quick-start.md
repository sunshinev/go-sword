# 快速开始

## 安装

```text
go get -u  github.com/sunshinev/go-sword
```

## 编译

```text
go build
```

安装完成后，确保`go-sword`命令在`GOPATH/bin`目录下，可执行

## 启动服务

```text
go-sword -db {db_database} -password {db_password} -user {db_uesr} -module {module_name}
```

例如：`go-sword -db blog -password 123456 -user root -module go-sword-app`

以上命令，就是连接数据库`blog`，用户名`root`,密码`12345`,在go-sword命令的当前目录下创建项目`go-sword-app`

启动成功的提示

```text
Go-Sword will create new project named go-sword-app in current directory

[Server info]
Server port : 8080
Project module : go-sword-app

[db info]
MySQL host : localhost
MySQL port : 3306
MySQL user : root
MySQL password : 123456

Start successful, server is running ...
Please request: http://localhost:8080
```

### 参数说明

```text
+---------------------------------------------------+
|                                                   |
|            Welcome to use Go-Sword                |
|                                                   |
|                Visualized tool                    |
|        Fastest to create CRUD background          |
|      https://github.com/sunshinev/go-sword        |
|                                                   |
+---------------------------------------------------+
Usage of go-sword:
  // 要连接的数据库信息
  -db string
      MySQL database
  -host string
      MySQL Host (default "localhost")
  // 重要：module参数单独作解释
  -module string
      New project module, the same as  'module' in go.mod file.   (default "go-sword-app/")
  // go-sword 服务启动的默认端口
  -p string
      Go-sword Server port (default "8080")
  -password string
      MySQL password
  -port int
      MySQL port (default 3306)
  -user string
      MySQL user
```

### 参数：  -module

`-module` 参数是代表要创建的项目名称，同时也是新项目`go.mod`文件中的`module`字段的值，这点请务必保持一致。

### 注意

新项目会在运行`go-sword`命令的当前目录下，直接创建`module`目录，作为新项目

## 开始使用服务

```text
Start successful, server is running ...
Please request: http://localhost:8080
```

根据服务启动的提示，直接点击`http://localhost:8080`即可进入web的可视化工具页面