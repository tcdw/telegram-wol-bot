# Telegram WoL Bot

一个能把你家电脑打开的 Telegram Bot。

## 需求

* 有另一台在同一局域网且长期开启的设备（比如 NAS），且安装了 Go 语言环境
* 你的电脑支持 [Wake-on-LAN](https://en.wikipedia.org/wiki/Wake-on-LAN)（现代 PC 应该都支持的）

## 安装

1. `go get -u -v https://github.com/tcdw/telegram-wol-bot`
2. 编辑 JSON 格式的配置文件：

```json5
{
    "token": "",    // 你的 Telegram Bot Token
    "chatID": 0,    // 使用该 bot 的聊天 ID
    "computers": [  // 电脑列表
        {
            "name": "home",                   // 电脑昵称
            "mac": "11:22:33:44:55:66",       // 电脑网卡的 MAC 地址
            "broadcast": "255.255.255.255:9"  // （可选）广播地址和端口号
        }
    ]
}
```

3. 在 BIOS 里启动你家电脑的 WoL 功能
4. 用 `telegram-wol-bot -c 你的配置文件.json` 把 bot 跑起来  
如果需要使用 http 代理，直接指定环境变量 `HTTP_PROXY=http://blahblah` 即可
5. 测试无误后，可以用 pm2 等工具让 bot 持续运行
6. 开始享用吧！

## 命令

* `/list` - 列出所有可供开启的电脑
* `/boot <name>` - 打开指定的电脑
