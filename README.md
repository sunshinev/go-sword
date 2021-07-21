Go-sword 利刃 V2.0.0

![GitHub last commit](https://img.shields.io/github/last-commit/sunshinev/go-sword)
![GitHub](https://img.shields.io/github/license/sunshinev/go-sword)
![GitHub repo size](https://img.shields.io/github/repo-size/sunshinev/go-sword)
![GitHub stars](https://img.shields.io/github/stars/sunshinev/go-sword?style=social)
![GitHub forks](https://img.shields.io/github/forks/sunshinev/go-sword?style=social)


> 一款基于Go语言的可视化web管理后台生成工具包
> 根据MySQL的表结构，创建CRUD的管理后台界面，开发者无需再重复手动的创建具有CRUD能力的页面，只需要点击按钮即可生成完整的管理后台


[官网 https://sunshinev.github.io/go-sword-home/](https://sunshinev.github.io/go-sword-home/)

详细请参阅文档
[文档 https://go-sword-doc.osinger.com/](https://go-sword-doc.osinger.com/)


![59384a43cbc382dec53dd76d169a5d001587995174.jpg](https://github.com/sunshinev/remote_pics/raw/master/59384a43cbc382dec53dd76d169a5d001587995174.jpg)

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

