package main

import (
	"bufio"
	"crypto/tls"
	"flag"
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

var WhiteDomainList []string

func GetList(list []string) []io.Reader {
	bodys := make([]io.Reader, 0)
	for _, l := range list {
		resp, err := req.Get(l)
		if err != nil {
			panic(err)
		}
		bodys = append(bodys, resp.Body)
	}
	return bodys
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
	for _, w := range WhiteDomainList {
		if strings.EqualFold(original, w) {
			return false
		}
	}
	return true
}

func Parase(list map[string]struct{}, writer *bufio.Writer, params []string) {
	domainSuffix := []string{".com.cn", ".net.cn", ".org.cn", ".gov.cn"}
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
			for _, sufix := range domainSuffix {
				if strings.Contains(k, sufix) {
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

func Resolve(bodys []io.Reader, list map[string]struct{}) {
	for _, body := range bodys {
		reader := bufio.NewReader(body)
		for {
			o, _, c := reader.ReadLine()
			if c == io.EOF {
				break
			}
			original := string(o)
			// 第一个字符为 # 时跳过
			if strings.IndexRune(original, '#') == 0 {
				continue
			}
			// 为空行时跳过
			if strings.TrimSpace(original) == "" {
				continue
			}
			// 用于 https://hosts-file.net/ad_servers.txt
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
			if strings.IndexRune(original, '#') == 0 {
				continue
			}
			if strings.ContainsRune(newOrg, '#') {
				newOrg = newOrg[:strings.IndexRune(newOrg, '#')]
			}
			// dnsmasq-list
			newOrg = strings.ReplaceAll(newOrg, "server=/", "")
			newOrg = strings.ReplaceAll(newOrg, "/114.114.114.114", "")
			if !formatter(newOrg) {
				continue
			}
			newOrg = strings.TrimSpace(newOrg)
			// 检测是否有端口号，有则移除端口号
			if strings.ContainsRune(newOrg, ':') {
				newOrg = newOrg[:strings.IndexRune(newOrg, ':')]
			}
			urlStr, err := url.Parse(newOrg)
			if err != nil {
				log.Println(newOrg)
				continue
			}
			list[urlStr.String()] = struct{}{}
		}
	}
}

func GetDomainList(path, suffix string, domans, params []string) error {
	file, err := os.Create(path + "/" + suffix)
	if err != nil {
		return err
	}
	domainList := make(map[string]struct{}, 0)
	Resolve(GetList(domans), domainList)
	log.Printf("%s total: %d", suffix, len(domainList))
	bw := bufio.NewWriter(file)
	Parase(domainList, bw, params)
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

func main() {
	file := flag.String("c", "", "")
	pars := flag.String("p", "", "")
	flag.Parse()
	params = make([]string, 0)
	params = strings.Split(strings.TrimSpace(*pars), ",")
	f, err := os.Open(strings.TrimSpace(*file))
	if err != nil {
		log.Println(err)
		return
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
	err = GetDomainList(route, "newDomain.txt", domainList, params)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("success")
}
