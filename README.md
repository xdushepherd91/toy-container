### 简介

本项目是在阅读[runc](https://github.com/opencontainers/runc)源代码的基础上，开始的一个基于go语言的容器玩具项目。
目标：
     1. 检验源码阅读成果
     2. 深入理解容器原理
     3. 深入理解linux操作系统



### 第一阶段

2019/11/15 11:13  
目标:  
1. 使用overlay2构造一个rootfs
2. 开发toy-container命令行工具    
     1. 添加run 子命令，可以开始运行一个容器
     2. 添加exec 子命令，可以进入一个容器的命名空间
3. 使用namespace
4. 使用cgroup限制其资源占用
5. 其他待定