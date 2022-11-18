package main

import (
	"log"
	"os"
	initial "selenium-test/initial"
	"selenium-test/study"
	"time"
)

func main() {
	log.Println("start", os.Args[1:])
	wbo := initial.WebDriverOptions{
		ChromeDriverPath: "./chromedriver",
		Port:             9091,
		Url:              "https://train.casicloud.com",
		CourseDetailUri:  "/#/train-new/class-detail/6ccdd860-be61-4578-b35d-825d80473481",
		Username:         "15618137573",
		Password:         "gmc951120",
	}

	args := os.Args[1:]
	if len(args) > 0 {
		for i := 0; i < len(args)-1; i = i + 2 {
			arg := args[i]
			switch arg {
			case "-u":
				wbo.Username = args[i+1]
			case "-p":
				wbo.Password = args[i+1]
			}
		}
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
