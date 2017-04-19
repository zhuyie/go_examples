# 本示例用于测试time.Ticker的个数和时间间隔对程序CPU占用率的影响

## 原理
runtime.timeproc需要在所有Ticker（及其它Timer）所设置的最小超时值的时刻被唤醒，然后向Ticker中的channel写入值。这个写操作进而会唤醒当前正等待对channel进行读操作的goroutine，来执行程序自己的定时触发逻辑。以一次触发为例，会导致如下事件发生：
* runtime.timeproc被调度执行
* runtime.timeproc执行完毕被yield走
* 用户goroutine被调度执行
* 用户goroutine执行完毕被yield走
若用户goroutine中的定时触发逻辑非常简单，则CPU资源会主要消耗在goroutine调度和channel读写上。所以当Ticker的时间间隔设置的非常小，导致这样的情况要频繁发生时，则整个进程会持续的发生大量的goroutine上下文切换。

## 测试结果
```
  PID USER      PRI  NI  VIRT   RES S CPU% MEM%   TIME+  Command
 6226 zhuyie     24   0  530G  4956 ?  7.4  0.0  0:17.00 ./ticker -num_goroutines 1  -duration 1
 6461 zhuyie     24   0  530G  5336 ? 33.4  0.0  0:15.16 ./ticker -num_goroutines 32 -duration 1
 6521 zhuyie     24   0  530G  5520 ? 68.2  0.0  0:19.27 ./ticker -num_goroutines 64 -duration 1
 6566 zhuyie     24   0  530G  5532 ?  4.8  0.0  0:02.19 ./ticker -num_goroutines 64 -duration 20
```

## pprof分析
在1/1的配置中，pprof的top情况如下：
```
go tool pprof http://127.0.0.1:12345/debug/pprof/profile
Fetching profile from http://127.0.0.1:12345/debug/pprof/profile
Please wait... (30s)
Saved profile in /Users/zhuyie/pprof/pprof.127.0.0.1:12345.samples.cpu.001.pb.gz
Entering interactive mode (type "help" for commands)
(pprof) top
2350ms of 2360ms total (99.58%)
Dropped 2 nodes (cum <= 11.80ms)
Showing top 10 nodes out of 46 (cum >= 490ms)
      flat  flat%   sum%        cum   cum%
     930ms 39.41% 39.41%      930ms 39.41%  runtime.mach_semaphore_signal
     600ms 25.42% 64.83%      600ms 25.42%  runtime.mach_semaphore_wait
     340ms 14.41% 79.24%      340ms 14.41%  runtime.mach_semaphore_timedwait
     260ms 11.02% 90.25%      260ms 11.02%  runtime.usleep
     220ms  9.32% 99.58%      220ms  9.32%  runtime.selectgoImpl
         0     0% 99.58%      220ms  9.32%  main.testFunc
         0     0% 99.58%      390ms 16.53%  runtime.chansend
         0     0% 99.58%       50ms  2.12%  runtime.entersyscallblock
         0     0% 99.58%       50ms  2.12%  runtime.entersyscallblock_handoff
         0     0% 99.58%      490ms 20.76%  runtime.exitsyscall
```

在64/1的配置中，pprof的top情况如下：
```
~ go tool pprof http://127.0.0.1:12345/debug/pprof/profile
Fetching profile from http://127.0.0.1:12345/debug/pprof/profile
Please wait... (30s)
Saved profile in /Users/zhuyie/pprof/pprof.127.0.0.1:12345.samples.cpu.005.pb.gz
Entering interactive mode (type "help" for commands)
(pprof) top
16.93s of 17.12s total (98.89%)
Dropped 14 nodes (cum <= 0.09s)
Showing top 10 nodes out of 47 (cum >= 1.30s)
      flat  flat%   sum%        cum   cum%
     4.45s 25.99% 25.99%      4.45s 25.99%  runtime.mach_semaphore_wait
     3.85s 22.49% 48.48%      3.85s 22.49%  runtime.usleep
     3.81s 22.25% 70.74%      3.81s 22.25%  runtime.mach_semaphore_signal
     2.54s 14.84% 85.57%      2.54s 14.84%  runtime.mach_semaphore_timedwait
     2.25s 13.14% 98.71%      2.28s 13.32%  runtime.selectgoImpl
     0.01s 0.058% 98.77%      2.13s 12.44%  runtime.sysmon
     0.01s 0.058% 98.83%     10.28s 60.05%  runtime.systemstack
     0.01s 0.058% 98.89%      5.31s 31.02%  runtime.timerproc
         0     0% 98.89%      2.28s 13.32%  main.testFunc
         0     0% 98.89%      1.30s  7.59%  runtime.chansend
```

在64/20的配置中，pprof的top情况如下：
```
~ go tool pprof http://127.0.0.1:12345/debug/pprof/profile
Fetching profile from http://127.0.0.1:12345/debug/pprof/profile
Please wait... (30s)
Saved profile in /Users/zhuyie/pprof/pprof.127.0.0.1:12345.samples.cpu.006.pb.gz
Entering interactive mode (type "help" for commands)
(pprof) top
1360ms of 1400ms total (97.14%)
Showing top 10 nodes out of 58 (cum >= 20ms)
      flat  flat%   sum%        cum   cum%
     490ms 35.00% 35.00%      490ms 35.00%  runtime.usleep
     380ms 27.14% 62.14%      380ms 27.14%  runtime.mach_semaphore_signal
     140ms 10.00% 72.14%      140ms 10.00%  runtime.mach_semaphore_wait
     130ms  9.29% 81.43%      600ms 42.86%  runtime.selectgoImpl
      60ms  4.29% 85.71%       60ms  4.29%  runtime.duffcopy
      50ms  3.57% 89.29%      650ms 46.43%  main.testFunc
      40ms  2.86% 92.14%       40ms  2.86%  runtime.mach_semaphore_timedwait
      40ms  2.86% 95.00%      100ms  7.14%  runtime.unlock
      20ms  1.43% 96.43%       20ms  1.43%  runtime.procyield
      10ms  0.71% 97.14%       20ms  1.43%  runtime.gopark
```
