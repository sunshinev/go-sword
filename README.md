# Go-sword 🗡️

[【点我访问中文文档_cn_zh】](https://github.com/sunshinev/go-sword/blob/master/README_zh.md)

Go-sword is a visual web management background generation tool based on Go language

![GitHub last commit](https://img.shields.io/github/last-commit/sunshinev/go-sword)
![GitHub](https://img.shields.io/github/license/sunshinev/go-sword)
![GitHub repo size](https://img.shields.io/github/repo-size/sunshinev/go-sword)
![GitHub stars](https://img.shields.io/github/stars/sunshinev/go-sword?style=social)
![GitHub forks](https://img.shields.io/github/forks/sunshinev/go-sword?style=social)

The goal is to quickly create a CRUD visual background

According to the table structure of MySQL, create a complete management background interface, developers no longer need to manually create CRUD-capable pages
Just click the button to generate a complete management background

![136e8b44d5d4acf00d5a63125928bd731587996269.jpg](https://github.com/sunshinev/remote_pics/raw/master/136e8b44d5d4acf00d5a63125928bd731587996269.jpg)

## Features
1. One-key generation without writing a line of code
2. Support add, delete, edit, list, batch delete, paging, search
3. The page is based on Vue.js + iView
4. A separate logical file is generated for each data table, and developers can seek to use Vue or iView to implement more feature-rich pages

![1626ee1d3300ac6db6669d63721d96381587996351.jpg](https://github.com/sunshinev/remote_pics/raw/master/1626ee1d3300ac6db6669d63721d96381587996351.jpg)

## Start

### Installation
```
go get -u github.com/sunshinev/go-sword
```
After the installation is complete, make sure that the `go-sword` command is in the` GOPATH/bin` directory, executable


### Start the service
```
go-sword -db {db_database} -password {db_password} -user {db_uesr} -module {module_name}
```

For example: `go-sword -db blog -password 123456 -user root -module go-sword-app`

The above command is to connect to the database `blog`, the username` root`, the password `12345`, and create the project` go-sword-app` in the current directory of the go-sword command

Tips for successful startup
```
Go-Sword will create new project named go-sword-app in current directory

[Server info]
Server port: 8080
Project module: go-sword-app

[db info]
MySQL host: localhost
MySQL port: 3306
MySQL user: root
MySQL password: 123456

Start successful, server is running ...
Please request: http://localhost: 8080
```


#### Parameter Description
```
+ ------------------------------------------------- -+
| |
| Welcome to use Go-Sword |
| |
| Visualized tool |
| Fastest to create CRUD background |
| https://github.com/sunshinev/go-sword |
| |
+ ------------------------------------------------- -+
Usage of go-sword:
 //Database information to be connected
  -db string
      MySQL database
  -host string
      MySQL Host (default "localhost")
 //Important: module parameters are explained separately
  -module string
      New project module, the same as 'module' in go.mod file. (Default "go-sword-app /")
 //The default port where the go-sword service starts
  -p string
      Go-sword Server port (default "8080")
  -password string
      MySQL password
  -port int
      MySQL port (default 3306)
  -user string
      MySQL user
```

#### Parameters: -module
The `-module` parameter is the name of the project to be created, and it is also the value of the` module` field in the `go.mod` file of the new project. Please make sure that this is consistent.

#### Note
The new project will directly create the `module` directory under the current directory where the` go-sword` command is run as the new project

### Start using the service

```
Start successful, server is running ...
Please request: http://localhost: 8080
```

According to the prompt of service startup, directly click `http://localhost: 8080` to enter the web visualization tool page

![59384a43cbc382dec53dd76d169a5d001587995174.jpg](https://github.com/sunshinev/remote_pics/raw/master/59384a43cbc382dec53dd76d169a5d001587995174.jpg)

#### Important: Introduction to page functions
1. First, select the MySQL table by drop-down, and then click the `Preview` button to render the file to be created
2. The first time you create a new project file you need to click `select all` to select all, the first time you create a core file that contains the necessary project start
3. Click the `Generate` button, you can see the prompt file successfully created
4. So far, our background has been successfully created

note:
1. For the first time, all files need to be selected
2. If you create a second management page, you can just select the `select diff & new` button and click the` Generate` button
3. Every time a new management interface is generated, please restart the newly created project

## Start using a new project
Go to our newly created project directory
```
➜ test tree -L 2
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
For example, if we just run the `go-sword` command in the` test` directory, the created project is `test/go-sword-app`

We enter the `test/go-sword-app` directory and start the project according to the following command

### Initialize new project go mod init

Use `go mod` to initialize the project, and the` module` here is the same as the project name we mentioned earlier! !

```
go mod init {module}
```

###Startup project

```
go run main.go
```

Then you will see the following prompt, click `http://localhost: 8082` to enter the background management interface

```
Enjoy your system ^ ^
Generated by Go-sword
https://github.com/sunshinev/go-sword

[Server info]
Server port: 8082

[db info]
MySQL host: localhost
MySQL port: 3306
MySQL user: root
MySQL password: 123456

Start successful, server is running ...
Please request: http://localhost: 8082
```

### Manage background effects

1. Back-end error notification
2. Add, delete, edit, list, bulk delete, paging, search

![1626ee1d3300ac6db6669d63721d96381587996351.jpg](https://github.com/sunshinev/remote_pics/raw/master/1626ee1d3300ac6db6669d63721d96381587996351.jpg)

## some problems
1. Because of the out-of-order traversal of golang's map structure, the field order of some pages cannot be guaranteed to be consistent with the database field order
2. There may be better solutions regarding the parameters of `module`
3. The ability to register and log in is not provided, which is not in line with the original intention. At the beginning, I wanted to do more basic things and quickly create pages
4. The generated project code has a lot of room for optimization

## Page function display

#### List
![ea1f86ebc1b5c88aaf6484fa078584951587997286.jpg](https://github.com/sunshinev/remote_pics/raw/master/ea1f86ebc1b5c88aaf6484fa078584951587997286.jpg)

#### Delete
![70279af696d9a230001f821cdf3a1ac21587997368.jpg](https://github.com/sunshinev/remote_pics/raw/master/70279af696d9a230001f821cdf3a1ac21587997368.jpg)

#### Preview
![2d1871a645acc3d3544ad7f77a0d6fca1587997398.jpg](https://github.com/sunshinev/remote_pics/raw/master/2d1871a645acc3d3544ad7f77a0d6fca1587997398.jpg)

#### Edit
![a9255db26b2af0365655840f6afd27851587997440.jpg](https://github.com/sunshinev/remote_pics/raw/master/a9255db26b2af0365655840f6afd27851587997440.jpg)



## Go-sword fork
If you want to customize, then you need to note that the Go-sword project can be packaged into a single command to execute, because all static files are also packaged

The static file compression commands are as follows:
```
go-bindata -o assets/resource/dist.go -pkg resource resource/dist/...
```

```
go-bindata -o assets/stub/stub.go -pkg stub stub/...
```

```
go-bindata -o assets/view/view.go -pkg view view/...
```