package main

import (
	"os"
	"fmt"
	"flag"
	"time"
	"strconv"
	"strings"
	"net/url"
	"math/rand"
	//"io/ioutil"
)

type resource struct {
	url		string
	target	string
	start	int
	end		int
}

var uaList = []string {
	"Mozilla/5.0 (Windows NT 6.2; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.90 Safari/537.36",
	"Mozilla/5.0 (compatible; MSIE 10.0; Windows NT 6.1; Trident/6.0; Touch; MASMJS)",
	"Mozilla/5.0 (X11; Linux i686) AppleWebKit/535.21 (KHTML, like Gecko) Chrome/19.0.1041.0 Safari/535.21",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.2999.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; U; Intel Mac OS X 10.4; en-US; rv:1.9.2.2) Gecko/20100316 Firefox/3.6.2",
}

func ruleResource () []resource {
	var res []resource
	//首页
	r1 := resource{
		url: "http://localhost:8888/",
		target: "",
		start: 0,
		end: 0,
	}
	//列表页 （数据库id 1-21）
	r2 := resource{
		url: "http://localhost:8888/list/{$id}.html",
		target: "{$id}",
		start: 1,
		end: 21,
	}
	//详情页 （数据库id 1-12924）
	r3 := resource{
		url: "http://localhost:8888/movie/{$id}.html",
		target: "{$id}",
		start: 1,
		end: 12924,
	}
	res = append(append(append(res, r1), r2), r3)
	return res
}

//地址处理
func buildUrl(res []resource) []string {
	var list []string //返回的数据

	for _,resItem := range res {
		//先处理首页
		if len(resItem.target) == 0 {
			list = append(list, resItem.url)
		} else {
			for i:=resItem.start; i<=resItem.end; i++ {
				urlStr := strings.Replace(resItem.url, resItem.target, strconv.Itoa(i), -1)
				list = append(list, urlStr)
			}
		}
	}

	return list
}

//获取随机数
func randInt(min, max int) int {
	//实例化一个对象， 传递种子值 rand.NewSource( time.Now().UnixNano() )
	r := rand.New( rand.NewSource( time.Now().UnixNano() ) )
	if min > max {
		return max
	}
	return r.Intn(max-min) + min
}

func main() {
	total		:= flag.Int("total", 100, "指定生成的行数");
	filePath	:= flag.String("filePath", "/Users/syl/nginx_access.log", "日志文件位置");
	flag.Parse()

	//构造真实网站的url
	res		:= ruleResource() //调用方法生成
	list	:= buildUrl(res)

	//按照规定格式，一次性生成日志字符串
	logStr := ""
	for i:=1; i <= *total; i++ {
		currentUrl 	:= list[randInt(0, len(list)-1)]
		referUrl	:= list[randInt(0, len(list)-1)]
		ua			:= uaList[randInt(0, len(uaList)-1)]
		logStr = logStr + makeLog(currentUrl, referUrl, ua) + "\n"
		//ioutil.WriteFile(*filePath, []byte(logStr), 0644) //覆盖写
	}
	fd,_ := os.OpenFile(*filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	fd.Write([]byte(logStr))
	fd.Close()

	fmt.Println("done \n")
}

//拼接日志
func makeLog(current, refer, ua string) string {
	u := url.Values{}
	u.Set("time", "1")
	u.Set("url", current)
	u.Set("refer", refer)
	u.Set("ua", ua)
	paramsStr := u.Encode()
	logTmp := "10.100.14.104 - - [19/Mar/2021:15:19:01 +0800] \"OPTIONS /nginx_access.log?{$paramsStr} HTTP/1.1\" 200 43 \"-\" \"{$ua}\" \"-\""
	log := strings.Replace(logTmp,"{$paramsStr}", paramsStr, -1)
	log = strings.Replace(log, "{$ua}", ua, -1)
	return log
}