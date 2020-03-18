# 域名列表解析

**对现有的域名列表进行格式化**

[点击这里](https://github.com/CalmLong/domain-parse/releases)下载最新版本

## 应用场景

例如将 `dnsmasq-list` 中的域名去掉前缀 `server=/` 以及后缀 `/114.114.114.114` 或者替换为其他的内容

## 功能

* 支持自定义域名前后缀
* 支持对主域名，子域名分别操作
* 支持自动移除域名列表中的注释(不支持网页)
* 支持识别 `http_proxy` 代理变量

## 参数

* `-c` 指定一个 `.txt` 文本路径，里面应当包含域名列表的 URL，每个一行

* `-p` 共四个参数，中间用英文的 `,` 隔开

`-p` 前两(0,1)个参数为子域名的前/后缀，后两个(2,3)为主域名的前后缀，
四个参数均为必填，不需要添加内容则用 `""` 代替

## 示例

* 替换为 dnsmasq 支持的格式：

`./domain-parse -c=url.txt -p=server=/,/114.114.114.114,server=/,/114.114.114.114`

* 替换为 V2Ray 支持的格式：

`./domain-parse -c=url.txt -p=full:,"",domain:,""`

* 替换为 AdGuardHome 支持的格式：

> 由于命令行本身的原因，对于一些特殊的字符你需要使用 "" 括起来表示字符串

`./domain-parse -c=url.txt -p="||,^,||,^"`

* 替换为 hosts 支持的格式：

`./domain-parse -c=url.txt -p="127.0.0.1 ,"",127.0.0.1 ,"""`

* 仅输出域名

`./domain-parse -c=url.txt -p="","","",""`

## 输出

工具输出一个无后缀名且为 `domain` 的文本，可用记事本等文本编辑器打开

## 贡献

遇到任何问题可以提交[issues](https://github.com/CalmLong/domain-parse/issues)

## 其他

`url.txt` 中的域名仅供测试(演示)使用，可能需要代理才可连接，与本项目无关。如有侵权可联系删除