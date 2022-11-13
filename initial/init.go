package initial

import (
	"fmt"
	"log"
	"time"

	"github.com/tebeka/selenium"
)

type WebDriverOptions struct {
	ChromeDriverPath string
	Port             int
	Url              string
	CourseDetailUri  string
}

func (wbo *WebDriverOptions) Init() *selenium.Service {
	opts := []selenium.ServiceOption{
		// selenium.Output(os.Stderr),
		selenium.ChromeDriver(wbo.ChromeDriverPath),
	}

	selenium.SetDebug(false)
	service, err := selenium.NewChromeDriverService(wbo.ChromeDriverPath, wbo.Port, opts...)
	if err != nil {
		panic(err)
	}

	return service
}

func (wbo *WebDriverOptions) CreateWebDriver() selenium.WebDriver {
	caps := selenium.Capabilities{"broswerName": "chrome"}
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", wbo.Port))
	if err != nil {
		panic(err)
	}

	if err = wd.Get(wbo.Url); err != nil {
		panic(err)
	}
	return wd
}

func (wbo *WebDriverOptions) Login(wd *selenium.WebDriver) error {
	we, err := (*wd).FindElement(selenium.ByID, "D38username")
	if err != nil {
		return err
	}
	err = we.SendKeys("15618137573")
	if err != nil {
		return err
	}

	psi, err := (*wd).FindElement(selenium.ByID, "D38pword")
	if err != nil {
		return err
	}
	err = psi.SendKeys("gmc951120")
	if err != nil {
		return err
	}

	lb, err := (*wd).FindElement(selenium.ByID, "D38login")
	if err != nil {
		return err
	}
	if lb != nil {
		lb.Click()
	}

	time.Sleep(3 * time.Second)
	return nil
}

func (wbo *WebDriverOptions) FindActivitiesToStart(wd *selenium.WebDriver) ([]selenium.WebElement, error) {
	err := (*wd).SetImplicitWaitTimeout(time.Duration(15))
	if err != nil {
		return nil, err
	}

	activityLinks, err := (*wd).FindElement(selenium.ByCSSSelector, ".today-activity")
	if err != nil {
		return nil, err
	}

	log.Println("activity links", activityLinks)
	aciitvityList, err := activityLinks.FindElement(selenium.ByTagName, "ul")
	if err != nil {
		return nil, err
	}
	activities, err := aciitvityList.FindElements(selenium.ByTagName, "li")
	if err != nil {
		return nil, err
	}
	log.Println("-----------------activities", activities)

	activitiesToStart := []selenium.WebElement{}
	for _, activity := range activities {
		a, err := activity.FindElement(selenium.ByCSSSelector, ".ms-train-state ")
		if err != nil {
			return nil, err
		}
		text, _ := a.Text()
		if text == "未完成" {
			activitiesToStart = append(activitiesToStart, activity)
		}
	}

	log.Println("activitiesToStart", activitiesToStart)

	return activitiesToStart, nil
}
