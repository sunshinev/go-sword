Go-sword(利刃) V2.0.0

![GitHub last commit](https://img.shields.io/github/last-commit/sunshinev/go-sword)
![GitHub](https://img.shields.io/github/license/sunshinev/go-sword)
![GitHub repo size](https://img.shields.io/github/repo-size/sunshinev/go-sword)
![GitHub stars](https://img.shields.io/github/stars/sunshinev/go-sword?style=social)
![GitHub forks](https://img.shields.io/github/forks/sunshinev/go-sword?style=social)


> 一款基于Go语言的可视化web管理后台生成工具包
> 根据MySQL的表结构，创建CRUD的管理后台界面，开发者无需再重复手动的创建具有CRUD能力的页面，只需要点击按钮即可生成完整的管理后台


[官网 https://sunshinev.github.io/go-sword-home/](https://sunshinev.github.io/go-sword-home/)

## 升级说明
v2.0.0 
1. 修改为以包引入的方式来启动工具+创建后端
2. 修改了stub的业务逻辑代码
3. 将所有的代码在生成前，使用gofmt进行格式化
4. 缺点，需要重启项目后端项目才能生效

v1.0.0
1. 单独的服务启动项目


Gosword会在项目指定目录，释放一个完整的后台代码，包括前端、后端

### 安装
```
go get -u  github.com/sunshinev/go-sword
```

### 配置文件说明
项目需要一个配置文件，采用yaml格式，除了数据库的配置，主要包括释放的目录、工具端口、后台端口
```
db:
  user: root
  password: '123456'
  database: test
  host: localhost
  port: 3306
root_path: admin22 # 后端项目释放的目录
tool_port: '8081'  # go-sword代码生成工具的端口
server_port: '8082' # 生成的后台项目的端口
```

### GIN框架中的应用
1. 在Gin项目中的main函数中，开启工具，import 代码包

```
package main

import (
	"log"

	"github.com/app/admin22/sword"
	"github.com/gin-gonic/gin"
	gosword "github.com/sunshinev/go-sword"
)

func main() {

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	// 日志
	log.SetFlags(log.Llongfile | log.Ldate)

	// 1. 开启工具->根据sql生成项目
	gosword.Init("config/go-sword.yaml").Run()

	// 原始gin项目
	_ = r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
```

使用`gosword.Init`加载配置文件，并且在项目中开启工具
```
gosword.Init("config/go-sword.yaml").Run()
```

2. 后台创建成功后，加入`sword.Run`，使用另外一个端口开启后端项目
```
func main() {

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	// 日志
	log.SetFlags(log.Llongfile | log.Ldate)

	// 1. 开启工具->根据sql生成项目
	gosword.Init("config/go-sword.yaml").Run()

	// 2. 加载生成的项目->重新启动
	sword.Run("config/go-sword.yaml")

	// 原始gin项目
	_ = r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

```


根据服务启动的提示，直接点击`http://localhost:8081`即可进入web的可视化工具页面

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

