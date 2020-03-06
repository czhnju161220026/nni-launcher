# NNI Launcher

## 代码结构

主要分为如下部分：

+ handler：处理请求，在handler.go中处理创建一个任务(Post)，请求任务的结果(Get)的请求。
+ main：程序入口所在的包
+ typed：自定义的数据结构
+ vendor：依赖
+ template：用于创建服务的一些yaml文件，测试用



## 测试

目前nni-launcher暂时部署在n167集群中进行测试。在nni-resource的nni-launcher-svc下，通过nodePort 30497暴露服务。

### 测试提交功能

可以发送包含类似如下请求体的post请求到 http://210.28.132.167:30497/api/v1/nni-exp/submit 开启一个tfoperator的调参任务

``` json
{
    "user": "jack",
    "workspace": "test",
    "trailConcurrency": 2,
    "num": 8,
    "target": "maximize",
    "command": "python3 mnist_tf.py",
    "gpuNum": 1,
    "trainer": "tf",
    "search_space": {
        "dropout_rate": {
            "_type": "uniform",
            "_value": [
                0.5,
                0.9
            ]
        },
        "conv_size": {
            "_type": "choice",
            "_value": [
                2,
                3,
                5,
                7
            ]
        },
        "hidden_size": {
            "_type": "choice",
            "_value": [
                124,
                512,
                1024
            ]
        },
        "batch_size": {
            "_type": "choice",
            "_value": [
                1,
                4,
                8,
                16,
                32
            ]
        },
        "learning_rate": {
            "_type": "choice",
            "_value": [
                0.0001,
                0.001,
                0.01,
                0.1
            ]
        }
    }
}
```

或者发送

``` json 
{
    "user": "jack",
    "workspace": "test",
    "trailConcurrency": 2,
    "num": 8,
    "target": "maximize",
    "command": "python3 mnist_pt.py",
    "gpuNum": 1,
    "trainer": "pt",
    "search_space": {
        "epoch": {
            "_type": "choice",
            "_value": [
                1,
                3,
                5,
                7,
                9
            ]
        },
        "lr": {
            "_type": "choice",
            "_value": [
                0.1,
                0.01,
                0.001
            ]
        },
        "bz": {
            "_type": "choice",
            "_value": [
                16,
                32,
                64,
                128
            ]
        }
    }
}
```

开启一个pytorch operator的调参任务。



### 测试获取数据

get http://210.28.132.167:30497/api/v1/nni-exp/logs?workspace=test&user=jack

