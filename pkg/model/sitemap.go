package model

import (
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/mucahitkurtlar/fft/pkg/config"
	"github.com/temoto/robotstxt"
)

type URL struct {
	Loc string `xml:"loc"`
}

type Sitemap struct {
	URLs []URL `xml:"url"`
}

type SitemapIndex struct {
	URLs []URL `xml:"sitemap"`
}

func GetRobotsTxt(url string) (robotsData *robotstxt.RobotsData, err error) {
	if url[len(url)-1:] != "/" {
		url = url + "/"
	}
	robotsTxtUrl := fmt.Sprintf("%s%s", url, "robots.txt")
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("GET", robotsTxtUrl, nil)
	if err != nil {
		return nil, err
	}

	for header, value := range config.Headers {
		req.Header.Add(header, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	robotsData, err = robotstxt.FromResponse(resp)
	if err != nil {
		return nil, err
	}

	return
}

func GetSiteMap(sitemapUrl string) (urls []string, err error) {
	resp, err := http.Get(sitemapUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if strings.Contains(string(body), "<urlset") {
		sm := Sitemap{}
		err = xml.Unmarshal(body, &sm)
		if err != nil {
			return nil, err
		}

		for _, url := range sm.URLs {
			urls = append(urls, url.Loc)
		}

	} else if strings.Contains(string(body), "<sitemapindex") {
		smi := SitemapIndex{}
		err = xml.Unmarshal(body, &smi)
		if err != nil {
			return nil, err
		}

		for _, sitemap := range smi.URLs {
			innerUrls, err := GetSiteMap(sitemap.Loc)
			if err != nil {
				return nil, err
			}
			urls = append(urls, innerUrls...)
		}
	}

	return
}

func GetUrls(siteUrl string) (urls []string, err error) {
	robotsData, err := GetRobotsTxt(siteUrl)
	if err != nil {
		return nil, err
	}

	for _, sitemap := range robotsData.Sitemaps {
		sitemapUrls, err := GetSiteMap(sitemap)
		if err != nil {
			return nil, err
		}
		urls = append(urls, sitemapUrls...)
	}

	return
}
