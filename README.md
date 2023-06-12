# telegram-claude-bot
一个golang 写的 telegram bot<br>
调取claude slack api来获取应答<br>
免费、无限制
<br>
<br>
## 运行前请先配置config.yaml.bak文件
- botconfig部分:
    - 采用webhook注册模式,处理效率更高
    - 你需要有自己的域名和证书
- claudeconfig部分:
    - 需要先配置你的slack
    - 配置教程可参考 https://github.com/LlmKira/claude-in-slack-server
- 配置完成后把配置文件重命名为config.yaml

<br>

##  开始运行 
```
go run main.go
```

<br>

##  问答效果
https://github.com/ack233/telegram-claude-bot/assets/136236457/96918706-7902-4cb8-ac9d-48272edf4898

<br>

##  加telegram claude 群体验
https://t.me/claude00000



<br>
<br>

## 参考:
- https://github.com/go-telegram-bot-api/telegram-bot-api
- https://github.com/LlmKira/claude-in-slack-server
