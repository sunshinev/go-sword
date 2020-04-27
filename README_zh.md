# Go-sword 🗡️

Go-sword(利刃)是一款基于Go语言的可视化web管理后台生成工具

![GitHub last commit](https://img.shields.io/github/last-commit/sunshinev/go-sword)
![GitHub](https://img.shields.io/github/license/sunshinev/go-sword)
![GitHub repo size](https://img.shields.io/github/repo-size/sunshinev/go-sword)
![GitHub stars](https://img.shields.io/github/stars/sunshinev/go-sword?style=social)
![GitHub forks](https://img.shields.io/github/forks/sunshinev/go-sword?style=social)

目标就是快速的创建CRUD可视化的后台

根据MySQL的表结构，创建完整的管理后台界面，开发者无需再重复手动的创建具有CRUD能力的页面
只需要点击按钮即可生成完整的管理后台

![136e8b44d5d4acf00d5a63125928bd731587996269.jpg](https://github.com/sunshinev/remote_pics/raw/master/136e8b44d5d4acf00d5a63125928bd731587996269.jpg)

## 特点
1. 一键生成，无需写一行代码
2. 支持增加、删除、编辑、列表、批量删除、分页、检索
3. 页面基于Vue.js + iView 
4. 针对每个数据表都生成了单独的逻辑文件，开发者可以求使用Vue或者iView来实现功能更加丰富的页面

![1626ee1d3300ac6db6669d63721d96381587996351.jpg](https://github.com/sunshinev/remote_pics/raw/master/1626ee1d3300ac6db6669d63721d96381587996351.jpg)

## 开始

### 安装
```
go get -u  github.com/sunshinev/go-sword
```
安装完成后，确保`go-sword`命令在`GOPATH/bin`目录下，可执行


### 启动服务
```
go-sword -db {db_database} -password {db_password} -user {db_uesr} -module {module_name}
```

例如：`go-sword -db blog -password 123456 -user root -module  go-sword-app`

以上命令，就是连接数据库`blog`，用户名`root`,密码`12345`,在go-sword命令的当前目录下创建项目`go-sword-app`

启动成功的提示
```
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


#### 参数说明
```
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

#### 参数：  -module
`-module` 参数是代表要创建的项目名称，同时也是新项目`go.mod`文件中的`module`字段的值，这点请务必保持一致。

#### 注意
新项目会在运行`go-sword`命令的当前目录下，直接创建`module`目录，作为新项目

### 开始使用服务

```
Start successful, server is running ...
Please request: http://localhost:8080
```

根据服务启动的提示，直接点击`http://localhost:8080`即可进入web的可视化工具页面

![59384a43cbc382dec53dd76d169a5d001587995174.jpg](https://github.com/sunshinev/remote_pics/raw/master/59384a43cbc382dec53dd76d169a5d001587995174.jpg)

#### 重要：页面功能介绍
1. 首先下拉选择MySQL 的表格，然后点击`Preview`按钮，即可渲染出需要创建的文件
2. 首次创建新项目文件需要点击`select all`全部选择，首次创建包含了项目启动必需的核心文件
3. 点击`Generate`按钮，既可以看到提示文件创建成功
4. 到目前为止，我们的后台已经创建成功了

注意：
1. 首次创建，文件需要全部选择
2. 如果创建第二个管理页面，那么可以只选择 `select diff & new`按钮，然后点击`Generate`按钮
3. 每次生成新的管理界面后，请重启新创建的项目

## 开始使用新项目
进入到我们新创建的项目目录
```
➜  test tree -L 2
.
└── go-sword-app
    ├── controller
    ├── core
    ├── go.mod
    ├── go.sum
    ├── main.go
    ├── model
    ├── resource
    ├── route
    └── view
```
比如说我们，刚刚是在`test`目录运行的`go-sword`命令，创建的项目就是`test/go-sword-app`

我们进入`test/go-sword-app`目录下按照以下命令启动项目

### 初始化新项目 go mod init

利用`go mod`初始化项目，这里的`module`就是我们前面讲到的要与项目名称保持一致！！

```
go mod init {module}
```

### 启动项目

```
go run main.go
```

然后会看到下面的提示，点击`http://localhost:8082`既可以进入后台管理界面

```
Enjoy your system ^ ^
Generated by Go-sword
https://github.com/sunshinev/go-sword

[Server info]
Server port : 8082

[db info]
MySQL host : localhost
MySQL port : 3306
MySQL user : root
MySQL password : 123456

Start successful, server is running ...
Please request: http://localhost:8082
```

### 管理后台效果

1. 后端报错提醒
2. 增加、删除、编辑、列表、批量删除、分页、检索

![1626ee1d3300ac6db6669d63721d96381587996351.jpg](https://github.com/sunshinev/remote_pics/raw/master/1626ee1d3300ac6db6669d63721d96381587996351.jpg)

## 一些问题
1. 因为golang的map结构遍历乱序的问题，部分页面输出的字段顺序不能保证和数据库字段顺序一致
2. 关于`module`的参数，可能还会有更好的解决方案
3. 没有提供用户注册、登录的能力，这也不符合初衷，最开始就是想做的更加基础，快速创建页面
4. 生成的项目代码，还有很大的优化空间 

## 页面功能展示

### 列表
![ea1f86ebc1b5c88aaf6484fa078584951587997286.jpg](https://github.com/sunshinev/remote_pics/raw/master/ea1f86ebc1b5c88aaf6484fa078584951587997286.jpg)

### 删除
![70279af696d9a230001f821cdf3a1ac21587997368.jpg](https://github.com/sunshinev/remote_pics/raw/master/70279af696d9a230001f821cdf3a1ac21587997368.jpg)

### 预览
![2d1871a645acc3d3544ad7f77a0d6fca1587997398.jpg](https://github.com/sunshinev/remote_pics/raw/master/2d1871a645acc3d3544ad7f77a0d6fca1587997398.jpg)

### 编辑
![a9255db26b2af0365655840f6afd27851587997440.jpg](https://github.com/sunshinev/remote_pics/raw/master/a9255db26b2af0365655840f6afd27851587997440.jpg)



## Go-sword fork
如果想要自定义的话，那么需要注意，Go-sword 项目可以打包成一个那单独的命令来执行，因为将所有的静态文件也进行了打包

静态文件压缩命令如下：
```
go-bindata -o assets/resource/dist.go -pkg resource resource/dist/...
```

```
go-bindata -o assets/stub/stub.go -pkg stub stub/...
```

```
go-bindata -o assets/view/view.go -pkg view view/...
```
