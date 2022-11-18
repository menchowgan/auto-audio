package main

import (
	"log"
	initial "selenium-test/initial"
	"selenium-test/study"
	"time"
)

func main() {
	wbo := initial.WebDriverOptions{
		ChromeDriverPath: "./chromedriver",
		Port:             9091,
		Url:              "https://train.casicloud.com",
		CourseDetailUri:  "/#/train-new/class-detail/6ccdd860-be61-4578-b35d-825d80473481",
	}

	service := wbo.Init()
	defer service.Stop()

	wd := wbo.CreateWebDriver()
	defer wd.Quit()

	time.Sleep(5 * time.Second)

	err := wbo.Login(&wd)
	if err != nil {
		panic(err)
	}

	if err = wd.Get(wbo.Url + wbo.CourseDetailUri); err != nil {
		panic(err)
	}

	time.Sleep(10 * time.Second)

	activitiesToStart, err := wbo.FindActivitiesToStart(&wd)
	if err != nil {
		panic(err)
	}

	if len(activitiesToStart) == 0 {
		log.Println("课程已学完")
		return
	}

	learnner := study.Activities{
		ActivitiesToStart: activitiesToStart,
	}

	learnner.Learn(&wd)

}
