# timer
crontab类型及delay类型的定时器服务，并带有重试机制。
比linux的crontab功能更多，也更方便管理及控制。  
  
  
## fixed
1、实现Linux crontab的功能；  
2、可以指定开始时间；(而 linux crontab只能从加入那刻开始计算)  
3、每一个固定的时间点都会回调任务，即使服务异常停服后重启，
都会补齐每一个执行时间点的任务    

## delay
1、延时任务，延时的时间不小于设定的时间；  
2、服务重启后，会补齐没有执行的延时服务；  

## callback retry
1、按照一定的时间规律重试，超过20天仍然没有成功，则可能不再重试  
2、回调服务需要幂等  

```
callback 
api：从配置文件读取，http协议
request：
{
  time_point: 回调的时间点，秒级时间戳
  id: 回调的任务id，具体参见 manage task 中的id
}
response: http 200 --- 即表示成功，否则失败

```

## manage task
api 调用即可管理  

```
增加 delay 任务
无论系统中是否存在Id任务，都是直接添加进系统中
 
api: AddDelay
request: 
{
  d: 延时时间间隔，单位 s
  id：任务id
}
response: {} 空json
```

```
增加 fixed 定时任务
1、所有数据都相同的重复添加，不会对系统有任何的状态改变，并不会从StartTime重新执行此任务，并且返回OK
2、任何新添加的任务(此次添加前系统中不存在)，都是从StartTime开始执行此任务，返回OK
3、Id相同，但是其他数据不同的任务，添加失败，返回IdConflict
4、Id 不能为空字符串
 
api: AddFixed
request: 
{
  start: 开始生效的时间，UNIX时间戳，单位 s；0---从now开始
  cron：* * * * *  (linux crontab)
  id：任务id
}
response: 
{
  status: 0 --- OK; 1 --- IdConflict
}
```

```
删除 fixed 定时任务
 
api: DelFixed
request: 
{
  ids：[id]; 全部ids, 数组类型数据
}
response: {} 空json
```

```
查找 fixed 任务是否存在

api: ExistFixed
request: 
{
  ids：[id]; 全部ids, 数组类型数据
}
response: 
{
  ids：[id]; 系统存在的所有fixed任务的ids, 数组类型数据
}
```

```
遍历 fixed 任务

api: VisitFixed
request: 
{
  id：开始遍历的任务id；空字符串表示从头开始
}
response: 
{
  results：[{
    start: 开始生效的时间，UNIX时间戳，单位 s；0---从now开始
    cron：* * * * *  (linux crontab)
    id：任务id
    }] 本次遍历的fixed任务, 数组类型数据
  nextId：下一次遍历可用的起始点，空字符串表示遍历完成
}
```
