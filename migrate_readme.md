container_migrate
====

1.简介
-----
***
容器热迁移项目由华中科技大学服务计算技术和系统教育部重点实验室倾力打造，主要面向的是数据中心。Docker+k8s已经成为了云计算环境的一种广泛使用方式，以容器作为服务载体的服务方式也成为主流。在资源调度过程中，会出现容器重新调度到新的节点上的情况，而现在服务变得轻量化，实时性要求变高，无法忍受长时间的服务中断问题，因此在线迁移成为了需求。本项目提出了一种Docker在线迁移尝试与解决方案，并将其应用于k8s中，实现对pod的在线迁移的尝试。

2.文件目录
-----
***
### 2.1 docker_migrator

此文件夹包含docker容器的热迁移工具，更多关于docker容器热迁移环境的部署和操作情况请查看docker_migrator内提供的文档资料。

### 2.2 docker-ce

此文件夹包含修改后的docker-ce源码，在原有的docker-ce源码上添加了容器热迁移部分功能源码，需要使用该源码编译的二进制文件替换系统中原有的docker。

### 2.3 k8s_migrate

此文件夹包含修改过的k8s源码，实现了在docker+k8s环境下pod的热迁移，详细情况请参考k8s_migrate内提供的文档资料。

3.技术架构
-----
***
* 容器引擎：Docker17.03.0-ce
* 集群管理工具：kubernetes v1.9
* 冻结与恢复工具：CRIU
* 内核版本：Linux 3.5以上
* 开发语言：Go
* 运行环境：Go1.6.2、Python2.7

4.相关文档
-----
***
我们提供了关于docker-migrator详细的资料[document](/docker_migrator/doc/document.md)，讲述如何安装、使用docker-migrator以及如何进行容器的热迁移。

5.技术支持和帮助
-----
***
* 关于docker-migrator和docker-ce的问题和bug，请联系786748095@qq.com
* 关于k8s的问题和bug，请联系953361637@qq.com

6.Licensing
-----
***
Copyright (c) 2014-2021 [CGCL Labs](http://grid.hust.edu.cn/)

Licensed under the  MulanPSL License, Version 2.0. See [LICENSE](/LICENSE) for the full license text.
