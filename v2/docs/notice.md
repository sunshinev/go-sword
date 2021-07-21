## 重要须知

### 数据库相关

Go-sword 根据MySQL数据表的表结构进行解析，生成CRUD页面，所以对于数据表有如下要求

**必须包含如下三个字段**

1. id - 主键id
2. created_at - 创建时间
3. updated_at - 更新时间


否则生成的页面无法正常工作


### 核心概念

#### Go-sword 工具

Go-sword工具，是生成代码的服务，需要在项目中引入Go-sword工具包，加载配置文件。然后会通过协程启动新服务


#### Go-sword CRUD服务

Go-sword CRUD 服务，是工具生成的CRUD管理后台服务，同样需要加载配置文件，然后会通过协程启动新服务