package initial

import (
	"fmt"
	"os"

	"github.com/tebeka/selenium"
)

type WebDriverOptions struct {
	ChromeDriverPath string
	Port             int
	Url              string
}

func (wbo *WebDriverOptions) Init() *selenium.Service {
	opts := []selenium.ServiceOption{
		selenium.Output(os.Stderr),
		selenium.ChromeDriver(wbo.ChromeDriverPath),
	}

	selenium.SetDebug(true)
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
