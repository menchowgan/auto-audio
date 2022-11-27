package main

import (
	"log"
	initial "selenium-test/initial"
	"time"

	"github.com/tebeka/selenium"
)

func CheckLogged(wd *selenium.WebDriver, result chan bool, errCh chan error) {
	for {
		err := focus(wd)
		if err != nil {
			errCh <- err
			return
		}

		loggedText, err := (*wd).FindElement(selenium.ByCSSSelector, ".logged-text")
		log.Println("loggedText", loggedText, err)

		if err != nil {
			log.Println("not found 1")
			time.Sleep(10 * time.Second)
			continue
		}

		displayed, err := loggedText.IsDisplayed()
		if err != nil {
			log.Println("not found 2")
			time.Sleep(10 * time.Second)
			continue
		}
		log.Println("displayed", displayed)
		if displayed {
			result <- true
			return
		}
		time.Sleep(3 * time.Second)
	}
}

func focus(wd *selenium.WebDriver) error {
	wins, err := (*wd).WindowHandles()
	if err != nil {
		return err
	}
	err = (*wd).SwitchWindow(wins[len(wins)-1])
	if err != nil {
		return err
	}
	return nil
}

const (
	ArticleUrl = "/72ac54163d26d6677a80b8e21a776cfa/9a3668c13f6e303932b5e0e100fc248b.html"
	VideoUrl   = "/4426aa87b0b64ac671c96379a3a8bd26/db086044562a57b441c24f2af1c8e101.html#1novbsbi47k-5"
)

func FindArticleToRead(wd *selenium.WebDriver) []selenium.WebElement {
	as := []selenium.WebElement{}
	we, _ := (*wd).FindElements(selenium.ByCSSSelector, ".text-wrap")
	log.Println("ats list", we)
	for i, title := range we {
		t, _ := title.FindElement(selenium.ByCSSSelector, ".text")
		s, _ := t.Text()
		log.Println("title", s)
		if i < 6 {
			as = append(as, t)
		}
	}
	return as
}

func StartStudy(wd *selenium.WebDriver) {
	endArticles, err := studyArticles(wd)
	if err != nil {
		panic(err)
	}
	if endArticles {
		return
	}
}

func scrollWindow(wd *selenium.WebDriver, scrollCh chan int) {
	i := 0
	for ; i < 60; i++ {
		time.Sleep(3 * time.Second)
		(*wd).ExecuteScript("window.scrollBy(0, 250);", nil)
	}
	if i == 60 {
		scrollCh <- 1
	}
}

func studyArticles(wd *selenium.WebDriver) (bool, error) {
	ats := FindArticleToRead(wd)
	log.Println("articles to read", ats)
	for _, a := range ats {
		a.Click()
		time.Sleep(3 * time.Second)

		focus(wd)
		scrollCh := make(chan int)
		go scrollWindow(wd, scrollCh)

		select {
		case <-scrollCh:
			wins, err := (*wd).WindowHandles()
			if err != nil {
				return false, err
			}
			(*wd).CloseWindow(wins[len(wins)-1])
			focus(wd)
		}
	}
	return true, nil
}

func studyVideos(wd *selenium.WebDriver) (bool, error) {
	focus(wd)
	videos, err := (*wd).FindElements(selenium.ByCSSSelector, ".innerPic")
	if err != nil {
		return false, nil
	}
	log.Println("videos", videos)
	if len(videos) == 0 {
		return true, nil
	}

	videos = videos[0:9]
	for _, video := range videos {
		video.Click()
		time.Sleep(2 * time.Second)
		for {
			focus(wd)
			we, _ := (*wd).FindElement(selenium.ByCSSSelector, ".replay-btn")
			if we != nil {
				break
			}
		}
		log.Println("video stop")
		wins, err := (*wd).WindowHandles()
		if err != nil {
			panic(err)
		}
		(*wd).CloseWindow(wins[len(wins)-1])
		focus(wd)
	}

	return true, nil
}

func timeoutCheck(timeoutCh chan bool) {
	time.Sleep(30 * time.Second)
	timeoutCh <- true
}

func main() {
	wbo := initial.WebDriverOptions{
		ChromeDriverPath: "./chromedriver",
		Port:             9091,
		Url:              "https://www.xuexi.cn",
	}

	service := wbo.Init()
	defer service.Stop()

	wd := wbo.CreateWebDriver()

	resultCh := make(chan bool)
	errCh := make(chan error)
	timeoutCh := make(chan bool)

	wins, _ := wd.WindowHandles()
	wd.MaximizeWindow(wins[len(wins)-1])

	go CheckLogged(&wd, resultCh, errCh)
	go timeoutCheck(timeoutCh)

	select {
	case result := <-resultCh:
		if result {
			log.Println("logged")
			wins, err := wd.WindowHandles()
			if err != nil {
				panic(err)
			}
			wd.CloseWindow(wins[0])
			if err = wd.Get(wbo.Url + ArticleUrl); err != nil {
				panic(err)
			}
			time.Sleep(5 * time.Second)
			StartStudy(&wd)

			if err = wd.Get(wbo.Url + VideoUrl); err != nil {
				panic(err)
			}
			time.Sleep(5 * time.Second)
			studyVideos(&wd)
		}
	case err := <-errCh:
		if err != nil {
			log.Println("Be panic and stop")
			panic(err)
		}
	case isTimeout := <-timeoutCh:
		log.Println("not logged in 30 seconds,", isTimeout)
		return
	}
}
