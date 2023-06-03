package model

import (
	"log"
	"os"
	"sync"

	"github.com/playwright-community/playwright-go"
)

func CrawlPage(browser playwright.Browser, url string, fontFamilies *[]string, crawCrawlerOpts *CrawlerOpts) error {
	page, err := browser.NewPage()
	if err != nil {
		return err
	}

	log.Printf("Navigating to url: %s", url)
	_, err = page.Goto(url, playwright.PageGotoOptions{Timeout: playwright.Float(crawCrawlerOpts.GoToTimeout), WaitUntil: (*playwright.WaitUntilState)(playwright.LoadStateLoad)})
	if err != nil {
		log.Printf("Error navigating to url: %s\n\t%s", url, err)
		return err
	}

	err = UpdateFontFamilies(page, fontFamilies)
	if err != nil {
		log.Printf("Error processing font families: %s\n\t%s", url, err)
		return err
	}

	log.Printf("Waiting for network idle: %s", url)
	page.WaitForNavigation(playwright.PageWaitForNavigationOptions{WaitUntil: (*playwright.WaitUntilState)(playwright.LoadStateNetworkidle), Timeout: playwright.Float(crawCrawlerOpts.NetIdleTimeout)})

	err = UpdateFontFamilies(page, fontFamilies)
	if err != nil {
		log.Printf("Error processing font families: %s\n\t%s", url, err)
		return err
	}

	log.Printf("Closing page: %s", url)
	err = page.Close()
	if err != nil {
		log.Printf("Error closing page: %s\n\t%s", url, err)
		return err
	}

	return nil
}

func CrawlSite(browser playwright.Browser, siteUrl string, crawlerOpts *CrawlerOpts) (fontFamilies []string, err error) {
	urls, err := GetUrls(siteUrl)
	if err != nil {
		return nil, err
	}

	if len(urls) == 0 {
		urls = append(urls, siteUrl)
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, crawlerOpts.GoRoutineCount)

	log.Printf("Found %d urls on %s", len(urls), siteUrl)

	if len(urls) > int(crawlerOpts.MaxPageCount) {
		urls = urls[:crawlerOpts.MaxPageCount]
	}

	log.Printf("Processing %d urls on %s", len(urls), siteUrl)

	for _, url := range urls {
		sem <- struct{}{}
		wg.Add(1)
		go func(url string) {
			defer func() {
				<-sem
				wg.Done()
			}()

			CrawlPage(browser, url, &fontFamilies, crawlerOpts)
		}(url)
	}

	wg.Wait()
	log.Printf("Found %d font families on %s", len(fontFamilies), siteUrl)

	return fontFamilies, nil
}

func UpdateFontFamilies(page playwright.Page, fontFamilies *[]string) error {
	fonts, err := page.Evaluate(`() => {
		if (typeof getComputedStyle == "undefined") {
			getComputedStyle = function (elem) {
			  return elem.currentStyle;
			}
		  }
		
		  var style,
			thisNode,
			styleId,
			allStyles = [],
			nodes = document.body.getElementsByTagName('*');
		
		  for (var i = 0; i < nodes.length; i++) {
			thisNode = nodes[i];
			if (thisNode.style) {
			  styleId = '#' + (thisNode.id || thisNode.nodeName + '(' + i + ')');
			  style = thisNode.style.fontFamily || getComputedStyle(thisNode, '')["fontFamily"];
		
			  if (style) {
				if (allStyles.indexOf(style) == -1) {
				  allStyles.push(style);
				}
		
				thisNode.dataset.styleId = allStyles.indexOf(style);
			  }
		
			  styleBefore = getComputedStyle(thisNode, ':before')["fontFamily"];
			  if (styleBefore) {
				if (allStyles.indexOf(styleBefore) == -1) {
				  allStyles.push(styleBefore);
				}
		
				thisNode.dataset.styleId = allStyles.indexOf(styleBefore);
			  }
		
			  styleAfter = getComputedStyle(thisNode, ':after')["fontFamily"];
			  if (styleAfter) {
				if (allStyles.indexOf(styleAfter) == -1) {
				  allStyles.push(styleAfter);
				}
		
				thisNode.dataset.styleId = allStyles.indexOf(styleAfter);
			  }
			}
		  }
		  return allStyles;
	}`)
	if err != nil {
		return err
	}

	for _, font := range fonts.([]interface{}) {
		if !Contains(*fontFamilies, font.(string)) {
			*fontFamilies = append(*fontFamilies, font.(string))
		}
	}

	return nil
}

func StartCrawler(pw *playwright.Playwright, siteUrl string, crawlerOpts *CrawlerOpts) error {
	browser, err := pw.Firefox.Launch()
	if err != nil {
		return err
	}

	fontFamilies, err := CrawlSite(browser, siteUrl, crawlerOpts)
	if err != nil {
		return err
	}

	for _, fontFamily := range fontFamilies {
		log.Printf("Found font family %s", fontFamily)
	}

	WriteToFile("font-families.txt", fontFamilies)

	if err := browser.Close(); err != nil {
		return err
	}

	return nil
}

func WriteToFile(file string, content []string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}

	for _, c := range content {
		_, err := f.WriteString(c + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}

func Contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}

	return false
}
