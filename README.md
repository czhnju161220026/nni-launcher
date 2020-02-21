# NNI Launcher

## 代码结构

主要分为如下部分：

+ handler：处理请求，在handler.go中处理创建一个任务(Post)，请求任务的结果(Get)，删除一个任务(Delete)的请求。
+ main：程序入口所在的包
+ typed：自定义的数据结构
+ test：测试使用
+ vendor：依赖
+ template：用于创建服务的一些yaml文件



## 用法

目前nni launcher已经打包成镜像 czh1998/nni-launcher:0.8，假设已经部署在了集群中。

目前，创建任务是从已有的镜像 czh1998/nni-demo:0.2创建，代码已经打包在了镜像中，真实的使用应该是代码和数据通过PV挂载到容器中。

### 创建一个任务

发送一个Post的请求到 http://service-ip:service-port/api/v1/, 需要在request body中包含创建任务需要的必要信息，例如：

``` json
{
    "user": "cuizihan",
    "workspace": "test",
    "trailConcurrency": 2,
    "num": 30,
    "target": "maximize",
    "command": "python3 train.py",
    "gpuNum": 0,
    "search_space": {
        "solver": {
            "_type": "choice",
            "_value": [
                "svd","cholesky","sparse_cg","sag","lsqr"
            ]
        },
        "alpha": {
            "_type": "choice",
            "_value": [
                1,0.5,0.1,0.01,5
            ]
        },
        "max_iter": {
            "_type": "choice",
            "_value": [
                1, 10, 15, 20
            ]
        }
    }
}
```
会在response中给出这个job的id。

### 查询调参过程

发送Get请求到 http://service-ip:service-port/api/v1/?workspace=xxx&id=xxx

在参数中给出工作区名和job的id。


 

### 删除某个任务

发送Delete请求到 http://service-ip:service-port/api/v1/?workspace=xxx&id=xxx

在参数中给出工作区名和job的id。



