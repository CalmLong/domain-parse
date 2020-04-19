# 域名列表解析

[点击这里](https://github.com/CalmLong/domain-parse/releases)下载最新版本

## 应用场景

将现有域名格式化为目标工具所支持的格式

例如将 `dnsmasq-list` 中的域名转换为其他工具支持的格式

其他详情可见[常见问题](https://github.com/CalmLong/domain-parse/issues/2#issue-585661994)

## 功能

* 自定义域名前后缀
* 支持对主域名，子域名分别操作
* 自动移除域名列表中的注释
* 自动移除重复域名
* 支持 `http_proxy` 变量

## 支持格式

* Hosts
* dnsmasq
* Adblock
* Pi-Hole
* Surge

## 参数

`-c` 一个 `.txt` 文本路径；里面应当包含域名列表的 URL，每个一行；参数不存在时默认加载同级目录中的 `url.txt`

`-v` 输出常见应用程序所支持的格式，并设定了一些默认值

 * `hosts`
 * `dnsmasq`
 * `v2ray`
 * `adblock`
 * `coredns`
 * `surge`
 * `only`
 
 其中
  
 * `dnsmasq` 和 `coredns` DNS 地址为 `114.114.114.114`
 * `hosts` 默认 IP 为 `0.0.0.0`
 * `surge` 默认规则为 `REJECT`
 * `only` 特殊选项，表示仅输出域名
 
`-e` 修改 `-v` 参数输出的默认值，支持

* `hosts`
* `dnsmasq`
* `coredns`
* `surge`

`-p` 自定义输出域名格式

> 前两个参数(0,1)为子域名的前/后缀，后两个参数(2,3)为主域名的前后缀；
> 中间用英文的 `;` 隔开;
> 四个参数均为必填，不需要添加内容则用 `""` 代替

## 示例

* 输出自定义格式的域名；通等于 `-v` 中的 `dnsmasq`

`./domain-parse -c=url.txt -p server=/;/114.114.114.114;server=/;/114.114.114.114`

* 输出适用于 Surge 等工具支持的格式

`./domain-parse -v surge`

* 输出 V2Ray 支持的格式

`./domain-parse -v v2ray`

* 指定 `dnsmasq` 解析域名的 IP

`./domain-parse -v dnsmasq -e 119.29.29.29`

> 输出结果为 `server=/example.com/119.29.29.29`

* 仅输出域名

`./domain-parse -v only`

## 输出

转换完成后输出一个无后缀名且为 `domain` 的文本，可用记事本等文本编辑器打开

## 其他

`url.txt` 中的域名仅供测试(演示)使用，可能需要代理才可连接，与本项目无关