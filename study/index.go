package study

import (
	"log"
	"strings"
	"time"

	"github.com/tebeka/selenium"
)

type Activities struct {
	ActivitiesToStart []selenium.WebElement
}

var goingCheck bool = false

func (a *Activities) Learn(wd *selenium.WebDriver) {
	for _, actvt := range a.ActivitiesToStart {
		log.Println("------------------not finished")
		log.Println("------------------Start Learning")
		log.Println("--------activity is", actvt)
		err := (*wd).SetImplicitWaitTimeout(time.Duration(15))
		if err != nil {
			panic(err)
		}

		link, _ := actvt.FindElement(selenium.ByCSSSelector, ".pointer")
		log.Println("link ", link)
		link.Click()
		time.Sleep(10 * time.Second)

		win, err := (*wd).WindowHandles()

		if err != nil {
			panic(err)
		}
		err = (*wd).SwitchWindow(win[len(win)-1])
		if err != nil {
			panic(err)
		}

		go going(wd)

		emuteVolume(wd)
		for {
			isEnd := keepLearning(&actvt, wd)
			if isEnd {
				goingCheck = false
				break
			}

			time.Sleep(10 * time.Second)
		}
		err = (*wd).CloseWindow(win[len(win)-1])
		if err != nil {
			panic(err)
		}
		err = (*wd).SwitchWindow(win[0])
		if err != nil {
			panic(err)
		}
	}
}

func keepLearning(activity *selenium.WebElement, wd *selenium.WebDriver) bool {
	sectionList := findSectionList(wd)
	if len(sectionList) == 0 {
		return true
	}

	sections := findToLearnSections(sectionList, wd)
	log.Println("sections", sections)
	if len(sections) == 0 {
		return true
	}

	for index, section := range sections {
		progress, statusText := getProgressAndStatus(&section, wd)
		elemClass, err := section.GetAttribute("class")
		if err != nil {
			panic(err)
		}

		if strings.Contains(elemClass, "focus") {
			break
		}
		log.Printf("学习状态：%s, 学习进度：%s", statusText, progress)
		section.Click()
		time.Sleep(5 * time.Second)
		log.Printf("section index %d is running\n", index)
	}

	return false
}

func findSectionList(wd *selenium.WebDriver) []selenium.WebElement {
	err := (*wd).SetImplicitWaitTimeout(time.Duration(15))
	if err != nil {
		panic(err)
	}
	tabs_cont_box, err := (*wd).FindElement(selenium.ByCSSSelector, ".tabs-cont-box")
	if err != nil {
		panic(err)
	}
	sectionArrow, err := tabs_cont_box.FindElement(selenium.ByCSSSelector, ".section-arrow")
	if err != nil {
		panic(err)
	}
	sectionList, err := sectionArrow.FindElements(selenium.ByCSSSelector, ".chapter-list-box")
	if err != nil {
		panic(err)
	}
	log.Println("sectionList", sectionList)
	return sectionList
}

func findToLearnSections(list []selenium.WebElement, wd *selenium.WebDriver) []selenium.WebElement {
	sections := []selenium.WebElement{}
	for _, section := range list {
		progress, statusText := getProgressAndStatus(&section, wd)

		if progress == "" && statusText == "" {
			continue
		}
		if statusText == "重新学习" || statusText == "已完成" {
			log.Println("本节已学完")
			continue
		}
		sections = append(sections, section)
	}
	return sections
}

func getProgressAndStatus(section *selenium.WebElement, wd *selenium.WebDriver) (string, string) {
	progress := ""
	statusText := ""
	err := (*wd).SetImplicitWaitTimeout(time.Duration(15))
	if err != nil {
		panic(err)
	}
	status, err := (*section).FindElement(selenium.ByCSSSelector, ".pointer")
	if err != nil {
		panic(err)
	}
	statusElems, err := status.FindElements(selenium.ByTagName, "span")
	if err != nil {
		panic(err)
	}
	if len(statusElems) == 1 {
		progress = "100%"
		statusText, _ = status.Text()
	} else if len(statusElems) == 2 {
		progress, _ = statusElems[0].Text()
		statusText, _ = statusElems[1].Text()
	} else {
		log.Println("无法获取进度")
	}
	return progress, statusText
}

func emuteVolume(wd *selenium.WebDriver) {
	err := (*wd).SetImplicitWaitTimeout(time.Duration(15))
	if err != nil {
		return
	}
	volumeButton, _ := (*wd).FindElement(selenium.ByCSSSelector, ".vjs-volume-menu-button")
	if volumeButton != nil {
		volumeButton.Click()
	}
}

func going(wd *selenium.WebDriver) {
	goingCheck = true
	defer func() {
		if i := recover(); i != nil {
			panic(i)
		}
	}()

	for {
		going, _ := (*wd).FindElement(selenium.ByCSSSelector, ".btn-ok")
		if going == nil {
			continue
		}

		if displayed, _ := going.IsDisplayed(); going != nil && displayed {
			err := going.Click()
			if err != nil {
				panic(err)
			}
		}
		time.Sleep(time.Second)
	}
}
