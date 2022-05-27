# DouYin
新建项目 
字节青训营抖音项目开发,希望项目成功完成！


# 文件夹说明
1.config配置文件，用于连接数据库等
2.处理流程
dao层连接数据库等
router层 接收url  -->  controller层 路由执行的函数  -->  logic层 可能查多个表，数据整合，逻辑判断  -->  model层 数据库的增删改查

# 0527更新
router 接收url
controller 路由执行函数，并返回响应消息 response
// 现在删除dao层，实际整合到model层，调用数据库处理函数是通过model层中dao的对象
service层 接收controller传递的参数，并调用model处理数据库
model层 操作数据库数据，进行增删改查

public 静态文件夹目录
用来存放视频video 和 封面cover 信息
通过生成的http连接能访问对应的视频和封面

**使用方法**
配置config.yaml文件 设置自己的数据库和ip port信息
go build 生成 DouYin.exe
./DouYin.exe运行
通过postman进行测试



