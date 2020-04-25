# base-fe

## 介绍

传统的前后端开发方式是由前端开发Vue单文件组件，作为每个页面。通过调用后端Api接口的方式来完成交互。

但是为了适应中小B端项目的敏捷开发，我们可能会需要后端开发来完成一些B端页面的编写，甚至是一键生成（请参考sunshinev/laravel-gii项目）

为此我们采用远程加载组件的方式，来实现编译页面的能力


## 项目结构

项目依赖如下组件:

```
"axios": "^0.19.0",
"core-js": "^3.3.2",
"http-vue-loader": "git+https://github.com/sunshinev/http-vue-loader.git",
"iview": "^3.5.3",
"view-design": "git+https://github.com/sunshinev/ViewUI.git",
"vue": "^2.6.10",
"vue-axios": "^2.1.5",
"vue-router": "^3.1.3"
```

`http-vue-loader`进行调整后，用于加载远程的Vue组件

`view-design`调整后用于开启所有的prefix前缀模式`i-*`


## 其他项目

[https://github.com/sunshinev/laravel-gii](https://github.com/sunshinev/laravel-gii)项目依赖该工作模式