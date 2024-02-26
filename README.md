# minichat

简洁的终端聊天应用

## 启动服务端
```bash
go run cmd/server/main.go
```

## 启动客户端
```bash
go run cmd/client/main.go
```

## 客户端命令
- `/who`
    - 显示在线用户名称
- `/rename 张三`
    - 更改用户名
- `/c hello,world.`
    - 公聊
- `/m username hello,world`
    - 私聊

## 公聊数据格式
```bash
[c][2024-02-25 08:58:28 127.0.0.1:51219] hello    #[(当前时间) (用户名)] (消息内容)
``` 

## 私聊数据格式
```bash
[m][2024-02-25 08:58:28 127.0.0.1:51219] hello    #[(当前时间) (用户名)] (消息内容)
``` 

