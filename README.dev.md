# Todo List
删除冗余的onex文件     
代码测试  
https方式 暂时不使用  
grpc方式   
包管理  
错误码生成    
参数校验生成:protoc-gen-validate 直接调用validate方法，在绑定json后校验，不使用kratos的中间件  
自定义校验 使用中间件 按照路由Switch   
参考iam生成swagger  也可参考makefile查询生成方式    
启动etcd的命令是什么  
biz层带有接口，是参考kratos的，暂不修改  
night watch修改成不依赖kratos，不使用wire  
usercenter比较复杂，暂不接入，先使用gin-jwt，参考miniblog，grpc方面参考iam，后续再考虑调用usercenter  
需要引入一个调用grpc的服务，可考虑使用batch来调用gateway  


# 需要的目录
部署dockerfile k8s     
文档  
sql  
# 其他
mongodb 分片后续可以使用多机器k8s部署


