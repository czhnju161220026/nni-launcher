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

### 准备

+ 需要Kubernetes集群

+ 需要拉取镜像：

  ``` shell
  docker pull czh1998/nni-launcher:1.0
  docker pull czh1998/nnidemo:0.3
  ```

+ 创建命名空间

  ``` shell
  kubectl create namespace nni-resource
  kubectl create namespace nni-exp # 提交的nni任务的pod都在这个空间里
  ```

+ 部署nni-launcher

  ``` shell
  # 找到位于template中的yaml文件
  kubectl apply -f dep.yaml
  kubectl apply -f service.yaml
  
  ```

  创建service后，执行如下命令，查看服务从哪个端口暴露

  ``` shell
  kubectl get svc -n nni-resource
  ```

  

+ 授权

  ```shell
  kubectl apply -f clusterrole.yaml
  ```



### 使用

+ 检测是否正常：

  ``` shell
  curl http://localhost:port/api/v1/nni-exp/hello
  Hello world
  ```

  

+ 提交任务

  目前用于测试的任务构建成了容器czh1998/nni-demo:0.3，提交一个任务需要向 http://localhost:port/api/v1/nni-exp发送post请求，并且request body中有如下信息:

  ``` json
  {
      "user": "cuizihan",
      "workspace": "test",
      "trailConcurrency": 2,
      "num": 8,
      "target": "maximize",
      "command": "python3 train.py",
      "gpuNum": 0,
      "search_space": {
          "solver": {
              "_type": "choice",
              "_value": [
                  "svd",
                  "cholesky",
                  "sparse_cg",
                  "sag",
                  "lsqr"
              ]
          },
          "alpha": {
              "_type": "choice",
              "_value": [
                  1,
                  0.5,
                  0.1,
                  0.01,
                  5
              ]
          },
          "max_iter": {
              "_type": "choice",
              "_value": [
                  1,
                  10,
                  100,
                  1000,
                  10000,
                  300000
              ]
          }
      }
  }
  ```

  user和workspace以及其余和任务本身无关的参数可以改动

+ 获取数据

  发送get请求到http://localhost:port/api/v1/nni-exp，请求体中需包含工作区和用户名

  ``` json
  {
      "workspace": "test",
      "user": "cuizihan"
  }
  ```

  