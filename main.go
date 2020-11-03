package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"
)

var params []string

var domainSuffix = []string{
	".com.cn", ".net.cn", ".org.cn", ".gov.cn", ".ah.cn", ".bj.cn", ".cq.cn", ".fj.cn",
	".gd.cn", ".gs.cn", ".gx.cn", ".gz.cn", ".ha.cn", ".hb.cn", ".he.cn", ".hi.cn", ".hk.cn", ".hn.cn", ".jl.cn",
	".js.cn", ".jx.cn", ".ln.cn", ".mo.cn", ".nm.cn", ".nx.cn", ".qh.cn", ".sc.cn", ".sd.cn", ".sh.cn", ".sn.cn",
	".sx.cn", ".tj.cn", ".tw.cn", ".xj.cn", ".yn.cn", ".zj.cn",
}

var localList = []string{
	"localhost",
	"ip6-localhost",
	"localhost.localdomain",
	"local",
	"broadcasthost",
	"ip6-loopback",
	"ip6-localnet",
	"ip6-mcastprefix",
	"ip6-allnodes",
	"ip6-allrouters",
	"ip6-allhosts",
	"0.0.0.0",
}

func GetList(list []string) []io.Reader {
	rs := make([]io.Reader, 0)
	for _, l := range list {
		log.Println("getting", l)
		resp, err := req.Get(l)
		if err != nil {
			panic(err)
		}
		rs = append(rs, resp.Body)
	}
	return rs
}

func DetectPath() (string, error) {
	str, err := os.Getwd()
	if err != nil {
		return "", err
	}
	dir := str + "/"
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return "", err
	}
	dir = strings.Replace(dir, "\\", "/", -1)
	return dir, nil
}

func formatter(original string) bool {
	for _, l := range localList {
		if l == original {
			return false
		}
	}
	return true
}

func Parse(list map[string]struct{}, writer *bufio.Writer, params []string) {
	domains := make([]string, 0)
	fulls := make([]string, 0)
	for k, _ := range list {
		if err := net.ParseIP(k); err != nil {
			continue
		}
		switch strings.Count(k, ".") {
		case 1:
			domains = append(domains, k)
		case 2:
			fulls = append(fulls, k)
			for _, suffix := range domainSuffix {
				if strings.Contains(k, suffix) {
					domains = append(domains, k)
					break
				}
			}
		default:
			fulls = append(fulls, k)
		}
	}
	sort.Strings(fulls)
	sort.Strings(domains)
	for i, f := range fulls {
		if i == 0 {
			continue
		}
		_, _ = writer.WriteString(params[0] + f + params[1] + "\n")
	}
	for _, d := range domains {
		_, _ = writer.WriteString(params[2] + d + params[3] + "\n")
	}
}

func deleteStr(newOrg string, str []string) string {
	for _, s := range str {
		newOrg = strings.ReplaceAll(newOrg, s, "")
	}
	return newOrg
}

func Resolve(reader []io.Reader, list map[string]struct{}) {
	const j = '#'
	for _, body := range reader {
		reader := bufio.NewReader(body)
		for {
			o, _, c := reader.ReadLine()
			if c == io.EOF {
				break
			}
			original := string(o)
			// 第一个字符为 # 或 ! 时跳过
			if strings.IndexRune(original, j) == 0 || strings.IndexRune(original, '!') == 0 {
				continue
			}
			// 为空行时跳过
			if strings.TrimSpace(original) == "" {
				continue
			}
			// 用于中间包含特殊空格的
			if strings.ContainsRune(original, '\t') {
				original = strings.ReplaceAll(original, "\t", " ")
			}
			newOrg := original
			// 移除前缀为 0.0.0.0 或者 127.0.0.1 (移除第一个空格前的内容)
			index := strings.IndexRune(original, ' ')
			if index > -1 {
				newOrg = strings.ReplaceAll(original, original[:index], "")
			}
			// 移除行中的空格
			newOrg = strings.TrimSpace(newOrg)
			// 再一次验证第一个字符为 # 时跳过
			if strings.IndexRune(original, j) == 0 {
				continue
			}
			if strings.ContainsRune(newOrg, j) {
				newOrg = newOrg[:strings.IndexRune(newOrg, j)]
			}
			// dnsmasq-list
			newOrg = deleteStr(newOrg, []string{"server=/", "/114.114.114.114"})
			if !formatter(newOrg) {
				continue
			}
			// adblock
			if strings.ContainsRune(newOrg, '^') {
				// 子域名包含 * 的不会被解析
				if strings.ContainsRune(newOrg, '*') {
					continue
				}
				// 表达式不会被解析
				if strings.Contains(newOrg, "/^") {
					continue
				}
				// 基础白名单规则会被一同解析
				newOrg = deleteStr(newOrg, []string{"||", "^", "@@"})
			}
			// Surge
			sg := len(strings.Split(newOrg, ","))
			switch sg {
			case 2, 3:
				newOrg = newOrg[strings.IndexRune(newOrg, ',')+1:]
				if sg == 3 {
					newOrg = newOrg[:strings.IndexRune(newOrg, ',')]
				}
			}
			newOrg = strings.TrimSpace(newOrg)
			// 检测是否有端口号，有则移除端口号
			if strings.ContainsRune(newOrg, ':') {
				newOrg = newOrg[:strings.IndexRune(newOrg, ':')]
			}
			urlStr, err := url.Parse(newOrg)
			if err != nil {
				log.Println("parse failed: ", newOrg)
				continue
			}
			urlString := urlStr.String()
			// 如果为 IP 则跳过
			if err := net.ParseIP(urlString); err != nil {
				continue
			}
			if strings.IndexRune(urlString, '.') == 0 {
				urlString = urlString[1:]
			}
			list[urlString] = struct{}{}
		}
	}
}

func GetDomainList(path, suffix string, domain, params []string) error {
	file, err := os.Create(path + "/" + suffix)
	if err != nil {
		return err
	}
	domainList := make(map[string]struct{}, 0)
	Resolve(GetList(domain), domainList)
	log.Printf("%s total: %d", suffix, len(domainList))
	bw := bufio.NewWriter(file)
	Parse(domainList, bw, params)
	_ = bw.Flush()
	return file.Close()
}

var req = &http.Client{Transport: transport()}

func transport() *http.Transport {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false,
		},
		DialContext: (&net.Dialer{
			Timeout: 180 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout: 60 * time.Second,
		ForceAttemptHTTP2:   true,
		DisableKeepAlives:   true,
		MaxIdleConnsPerHost: 50,
	}
}

var vTools = []string{
	dnsmasq, v2ray, hosts, adblock, surge, only,
}

const (
	dnsmasq = "dnsmasq"
	v2ray   = "v2ray"
	hosts   = "hosts"
	adblock = "adblock"
	surge   = "surge"
	only    = "only"
)

const dnsIP = "114.114.114.114"

func domainFormat(value, exp string) []string {
	vParams := make([]string, 0)
	switch value {
	case dnsmasq:
		v := "/" + dnsIP
		if exp != "" {
			v = "/" + exp
		}
		vParams = append(vParams, "server=/", v, "server=/", v)
	case v2ray:
		vParams = append(vParams, "full:", "", "domain:", "")
	case adblock:
		vParams = append(vParams, "||", "^", "||", "^")
	case hosts:
		v := "0.0.0.0 "
		if exp != "" {
			v = exp + " "
		}
		vParams = append(vParams, v, "", v, "")
	case surge:
		v := ",REJECT"
		if exp != "" {
			v = "," + exp
		}
		vParams = append(vParams, "DOMAIN-SUFFIX,", v, "DOMAIN-SUFFIX,", v)
	case only:
		vParams = append(vParams, "", "", "", "")
	default:
		panic(fmt.Sprintln(value, " is an unsupported format"))
	}
	return vParams
}

func main() {
	file := flag.String("c", "url.txt", "")
	val := flag.String("v", "", "dnsmasq\nv2ray\nhosts\nadblock\ncoredns\nsurge\nonly")
	exp := flag.String("e", "", "-v takes effect when\ndnsmasq\nhosts\ncoredns\nsurge")
	pars := flag.String("p", "", "customize")
	flag.Parse()
	params = make([]string, 0)
	params = strings.Split(strings.TrimSpace(*pars), ";")
	f, err := os.Open(strings.TrimSpace(*file))
	if err != nil {
		log.Println(err)
		return
	}
	vl := strings.TrimSpace(*val)
	if vl != "" {
		for _, v := range vTools {
			if v == *val {
				params = domainFormat(v, strings.TrimSpace(*exp))
			}
		}
		if len(params) < 2 {
			panic(fmt.Sprintln(vl, " is an unsupported format"))
			return
		}
	}
	body := bufio.NewReader(f)
	domainList := make([]string, 0)
	for {
		i, _, e := body.ReadLine()
		if e == io.EOF {
			break
		}
		d := string(i)
		if strings.TrimSpace(d) == "" {
			continue
		}
		u, err := url.Parse(d)
		if err != nil {
			log.Println(err)
			return
		}
		domainList = append(domainList, u.String())
	}
	route, err := DetectPath()
	if err != nil {
		log.Println(err)
		return
	}
	err = GetDomainList(route, "domain", domainList, params)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("success")
}
