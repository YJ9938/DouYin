# DouYin
新建项目 
字节青训营抖音项目开发,希望项目成功完成！


# 文件夹说明
1.config配置文件，用于连接数据库等
2.处理流程
router层 接收url  -->  controller层 路由执行的函数  -->  service层对上一层屏蔽操作数据库层，处理多个表，返回相应的结果 -->  model层 数据库的增删改查

3.public 静态文件夹目录
用来存放视频video 和 封面cover 信息
通过生成的http连接能访问对应的视频和封面

**使用方法**
配置config.yaml文件 设置自己的数据库和ip port信息
go build 生成 DouYin.exe
./DouYin.exe运行
通过抖音APK或postman进行测试



