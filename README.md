# Distributed-Systems
MIT课程[《Distributed Systems 》](http://nil.csail.mit.edu/6.824/2016/schedule.html)学习和翻译
+ 翻译和完成课程的实验代码，之后在代码里添加了注释说明，去掉了代码实现
+ 整理课程，编写简单的分布式入门教程


#### 资料推荐
+ [《大规模分布式存储系统》](https://book.douban.com/subject/25723658/)
+ [《分布式系统原理介绍》](http://pan.baidu.com/s/1geU1XAz)
+ [awesome-distributed-systems](https://github.com/kevinxhuang/awesome-distributed-systems)
+ [一名分布式存储工程师的技能树是怎样的？](https://www.zhihu.com/question/43687427/answer/96306564)
+ [袖珍分布式系统](http://www.jianshu.com/c/0cf64976a481)



#### pinewu's note

- 脚本运行的问题解决

  git clone这个项目后，如果未能正确的运行makefile，各种sh脚本，没关系，可以手动编译和运行与测试

  ```
  # lab1
  1. 直接去main目录下 go build wc.go 生成wc二进制文件
  2. 单worker 顺序性lab 测试
     直接在main目录下执行 ./wc master sequential pg-*
     执行sort -n -k2 mrtmp.wcseq | tail -10
     得出的结果大致如下，若相差不大，则可认为没问题
     he: 34077
     was: 37044
     that: 37495
     I: 44502
     in: 46092
     a: 60558
     to: 74357
     of: 79727
     and: 93990
     the: 154024
  3. 分布式的worker
     这个可以运行go test文件来测试
     在 6.824目录下 export "GOPATH=$PWD"
     在mapreduce目录下 go test -run TestBasic
  ```
