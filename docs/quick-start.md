## 安装须知
在整个文档中，将使用Gin框架作为演示，Go-sword工具本身可以适用于任意Golang框架

## 安装依赖包
```
go get -u github.com/sunshinev/go-sword/v2
```

## 启动Go-sword 工具

使用`gosword.Init`加载配置文件，并且在项目中开启工具；
该命令会解析yaml配置文件中的数据库配置，以及启动端口，并且在对应端口启动服务
```
gosword.Init("config/go-sword.yaml").Run()
```

在项目中的main函数中，import依赖包，开启工具，参考代码如下：
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

	// 这里我们可以调整日志的详细程度
	log.SetFlags(log.Llongfile | log.Ldate)

	// 核心注意：这里通过该命令，加载配置文件，并且启动服务
	gosword.Init("config/go-sword.yaml").Run()

	// 原始gin项目
	_ = r.Run()
}
```


## 控制台输出
如果配置成功，那么控制台会打印如下信息

```
START----------------------------------------
Go-Sword will create new project named admin22 in current directory
[Server info]
Server port : 8081
Project module : admin22
[db info]
MySQL host : localhost
MySQL port : 3306
MySQL user : root
MySQL password : 123456

Start successful, server is running ...
Please request: http://localhost:8081
END----------------------------------------
```


## 工具页面

点击链接，即可以在本地启动该工具

![136e8b44d5d4acf00d5a63125928bd731587996269.jpg](https://cdn.jsdelivr.net/gh/sunshinev/remote_pics/136e8b44d5d4acf00d5a63125928bd731587996269.jpg)
