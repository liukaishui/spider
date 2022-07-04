package main

import (
	"app/internal/common"
	"app/internal/models"
	"errors"
	"fmt"
	"github.com/gocolly/colly/v2"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

func main() {
	var result models.Sites
	err := common.DB.Where("id = 1").Take(&result).Error
	if err != nil {
		log.Fatal(err)
	}
	handle(&result)
}

// 处理 site
func handle(result *models.Sites) {
	// 执行时间
	executionTime := time.Now()

	// 解析链接
	u, err := url.Parse(result.Url)
	if err != nil {
		log.Fatal(err)
	}

	// 新建 colly
	c := colly.NewCollector(
		//colly.UserAgent("colly"),
		colly.AllowedDomains(u.Hostname()),
		//colly.MaxDepth(2),
		//colly.AllowURLRevisit(),
		colly.Async(true),
		//colly.URLFilters(
		//	regexp.MustCompile("http://httpbin\\.org/(|e.+)$"),
		//	regexp.MustCompile("http://httpbin\\.org/h.+"),
		//),
	)

	// 优化连接
	c.WithTransport(&http.Transport{
		DisableKeepAlives: true,
	})

	// 限制速度
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 9,
		Delay:       5 * time.Second,
	})

	// 请求
	//c.OnRequest(func(r *colly.Request) {
	//	fmt.Println(common.Func.NowDate(), "Request", r.URL)
	//	r.Headers.Set("User-Agent", randomString())
	//})

	// 错误
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(common.Func.NowDate(), "Error", r.Request.URL, err.Error())

		// 非详情页不存储
		if ok, _ := isDetailsUrl(r.Request.URL.Path); !ok {
			return
		}

		// 入库
		timeNow := time.Now()
		sqlErr := common.DB.Model(&models.Urls{}).Create(map[string]interface{}{
			"site_id":             result.ID,
			"url":                 r.Request.URL.Path,
			"info":                "{}",
			"content":             err.Error(),
			"status_code":         r.StatusCode,
			"last_execution_time": executionTime,
			"created_at":          timeNow,
			"updated_at":          timeNow,
		}).Error
		if sqlErr != nil {
			fmt.Println(common.Func.NowDate(), "DB Error", err)
		}
	})

	// 响应
	//c.OnResponse(func(r *colly.Response) {
	//	fmt.Println(common.Func.NowDate(), "Response", r.Request.URL)
	//})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		// 如果是播放页 就不再请求
		if ok, _ := isPlayUrl(link); ok {
			//e.Response.Ctx.Put("is_play_url", "1")
			return
		}

		e.Request.Visit(link)
		//c.Visit(e.Request.AbsoluteURL(link))
	})

	c.OnHTML("html", func(e *colly.HTMLElement) {
		if ok, _ := isDetailsUrl(e.Request.URL.Path); ok {
			e.Response.Ctx.Put("is_details_url", "1")
		} else {
			return
		}

		title := e.ChildText("title")
		keywords := e.ChildAttr(`meta[name="keywords"]`, "content")
		description := e.ChildAttr(`meta[name="description"]`, "content")
		content, err := e.DOM.Find("div.ct-c").First().Html()
		if err != nil {
			fmt.Println("Content Error", err)
		}

		e.Response.Ctx.Put("title", strings.TrimSpace(title))
		e.Response.Ctx.Put("keywords", strings.TrimSpace(keywords))
		e.Response.Ctx.Put("description", strings.TrimSpace(description))
		e.Response.Ctx.Put("content", strings.TrimSpace(content))
		html, _ := e.DOM.Html()
		e.Response.Ctx.Put("html", strings.TrimSpace(html))
	})

	// 抓取
	c.OnScraped(func(r *colly.Response) {
		fmt.Println(common.Func.NowDate(), "Scraped", r.Request.URL)

		//isDetailsUrl := r.Ctx.Get("is_details_url")
		//if isDetailsUrl != "1" {
		//	return
		//}

		if ok, _ := isDetailsUrl(r.Request.URL.Path); !ok {
			return
		}

		title := r.Ctx.Get("title")
		keywords := r.Ctx.Get("keywords")
		description := r.Ctx.Get("description")
		content := r.Ctx.Get("content")
		html := r.Ctx.Get("html")

		// 入库
		timeNow := time.Now()
		sqlErr := common.DB.Model(&models.Urls{}).Create(map[string]interface{}{
			"site_id": result.ID,
			"url":     r.Request.URL.Path,
			"info": common.Func.JSONEncode(map[string]interface{}{
				"title":       title,
				"keywords":    keywords,
				"description": description,
				"content":     content,
			}),
			"content":             html,
			"status_code":         r.StatusCode,
			"last_execution_time": executionTime,
			"created_at":          timeNow,
			"updated_at":          timeNow,
		}).Error
		if sqlErr != nil {
			fmt.Println("DB Error", sqlErr)
		}
	})

	c.Visit(result.Url)
	c.Wait()

	result.LastExecutionTime = executionTime
	common.DB.Select("last_execution_time").Save(&result)
}

// 随机字符串
func randomString() string {
	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, rand.Intn(10)+10)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// 是否详细页
func isDetailsUrl(href string) (bool, error) {
	if href == "" {
		return false, errors.New("empty url")
	}

	u, err := url.Parse(href)
	if err != nil {
		return false, err
	}

	r := regexp.MustCompile(`^[/]?[A-Za-z]+/[0-9]+/[0-9]+\.html$`)
	if r.MatchString(u.Path) {
		return true, nil
	} else {
		return false, nil
	}
}

// 是否播放页
func isPlayUrl(href string) (bool, error) {
	if href == "" {
		return false, errors.New("empty url")
	}
	u, err := url.Parse(href)
	if err != nil {
		return false, err
	}

	r := regexp.MustCompile(`^[/]?vod-play-id-[A-Za-z0-9\-]+\.html$`)
	if r.MatchString(u.Path) {
		return true, nil
	} else {
		return false, nil
	}
}
