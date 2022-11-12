package main

import (
	initial "selenium-test/Init"

	"github.com/tebeka/selenium"
)

func main() {
	wbo := initial.WebDriverOptions{
		ChromeDriverPath: "./chromedriver",
		Port:             9091,
		Url:              "http://www.baidu.com",
	}
	service := wbo.Init()
	defer service.Stop()

	wd := wbo.CreateWebDriver()
	we, err := wd.FindElement(selenium.ByCSSSelector, "#kw")

	if err != nil {
		panic(err)
	}

	err = we.SendKeys("hhh")
	if err != nil {
		panic(err)
	}

	for {

	}

}
