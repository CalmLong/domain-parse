# 域名列表解析

**对现有的域名列表进行格式化**

[点击这里](https://github.com/CalmLong/domain-parse/releases)下载最新版本

## 应用场景

将现有域名格式化为目标工具所支持的格式

**通常需要解析标准或者非标准的 `hosts` 域名列表才应是使用本工具的主要原因**

例如将 [dnsmasq-list](https://github.com/felixonmars/dnsmasq-china-list/blob/master/accelerated-domains.china.conf) 中的域名去掉前缀 `server=/` 以及后缀 `/114.114.114.114` 
或者替换为你想要的内容亦或者仅输出纯粹的域名

其他详情可见[常见问题](https://github.com/CalmLong/domain-parse/issues/2#issue-585661994)

## 功能

* 自定义域名前后缀
* 对主域名，子域名分别操作
* 自动移除域名列表中的注释(不支持网页)
* 自动移除重复域名
* 可识别 `http_proxy` 代理变量

## 支持格式

可被解析的域名格式：

* hosts(标准或者非标准)
* dnsmasq
* adblock(基础)

输出的格式：

可为域名前后添加任意内容

> 只能为域名前后添加任意字符，不支持域名内添加字符；
> 例如 `github.com` 不能被输出为 `git*hub.com`

## 参数

`-c` 指定一个 `.txt` 文本路径，里面应当包含域名列表的 URL，每个一行

`-v` 输出常见应用程序所支持的格式

 * hosts
 * dnsmasq
 * v2ray
 * adblock
 * coredns
 * only(仅输出域名)
 
`-e` 修改 `-v` 参数输出的默认值，仅支持 `dnsmasq` 和 `hosts`

`-p` 自定义输出域名格式，中间用英文的 `;` 隔开

> `-p` 前两(0,1)个参数为子域名的前/后缀，后两个(2,3)为主域名的前后缀，
> 四个参数均为必填，不需要添加内容则用 `""` 代替

## 示例

* 输出自定义格式的域名；通等于 `-v` 中的 `dnsmasq`

`./domain-parse -c=url.txt -p=server=/;/114.114.114.114;server=/;/114.114.114.114`

* 输出适用于 Kitsunebi 

`./domain-parse -p "DOMAIN-SUFFIX,;,Reject;DOMAIN-SUFFIX,;,Reject"`

* 输出 V2Ray 支持的格式：

`./domain-parse -c=url.txt -v v2ray`

* 指定 `dnsmasq` 解析域名的 IP

`./domain-parse -c url.txt -v dnsmasq -e 119.29.29.29`

输出结果为 `server=/example.com/119.29.29.29`

* 仅输出域名

`./domain-parse -c=url.txt -v only`

## 输出

工具输出一个无后缀名且为 `domain` 的文本，可用记事本等文本编辑器打开

## 贡献

遇到任何问题可以提交[issues](https://github.com/CalmLong/domain-parse/issues)

## 其他

`url.txt` 中的域名仅供测试(演示)使用，可能需要代理才可连接，与本项目无关。如有侵权可联系删除