# rate-limit

使用 Redis 实现限流组件  

针对日常使用的限流大概可以分成三种     

### 1、接口一定时间内限制请求的次数；  

这种限流是最常见的场景  

```go
TimesLimiter(ctx context.Context, key string, maxThreads, expireTime int64) error
```

这种场景的使用，一般我们会提供一个限流的 Key ,然后就是限流对应的线程数和限流的时长  

### 2、限制用户单接口的单次访问，用户第一个请求处理好了，第二次请求才能发起；   

这种场景可能不太好理解，举个栗子  

如果用户的手速很快，或者接口被恶意请求，那么同一个用户的在同一个时刻，可能使用相同的请求数据发起很多次的数据请求。  

这样请求只有一个有效，但是这些相同的请求，相互之间就会产生竞争，比如有时候数据库存在读写延迟，这种异常的请求就可能出现问题。   

所以可以在入加个限流，对于同一个用户的请求，只有第一次完成了，后面才能在发起   

```go
UserSingleRequestLimiter(ctx context.Context, accountId string, expireTime int64) error
```

然后使用的时候可以考虑，把用户ID，和对应的接口名一起组装成对应的限流 KEY。   

### 3、用户请求次数的限制，避免用户的恶意请求。   

这种纯属是应对恶意请求，发起者使用同一个用户信息，反复对多个接口发起恶意请求  

那么我们可以认为没这个用户在一段时间内的请求，超过某个值吗，后面的请求就可以拦截了   

```go
UserTimesLimiter(ctx context.Context, accountId string, maxThreads, expireTime int64) error
```

未实现部分  

key 设置过期时间，之后可以使用 LUA 保证命令的原子性  


