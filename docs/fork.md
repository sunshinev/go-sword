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
