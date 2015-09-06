# goframe

一个基于golang语言的后台服务框架。通过简单的修改以及添加模块可以实现后台守护进程的实现。


+ srv
 - srv下提供了简单的http服务，可以查看运行时的一些状态。
 通过　`go tool pprof http://ip:port/debug/pprof/heap`可以查看内存使用情况。
 或者通过浏览器"http://ip:port/debug/pprof"可以查看运行的内存、channel、线程等的使用情况。


+ util
 - logs 提供了按天切割的日志。
 - minuteTopn 5分钟内根据某个值的top10排序。通过共享内存或者channel可以让srv下httpsrv模块提供查询接口。
