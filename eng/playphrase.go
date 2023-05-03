package eng

import (
//	"fmt"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"log"
	"strings"
	"time"
)

func Play(words []string) {
	caps := selenium.Capabilities{
		"browserName": "chrome",
		"chromeOptions": map[string]interface{}{
			"excludeSwitches": [1]string{"enable-automation"},
		},
	}
	// Set up Chrome options
	chromeCaps := chrome.Capabilities{
		Args: []string{
			"--disable-gpu",
			"--no-sandbox",
			"--disable-infobars",
			"--disable-extensions",
			"--disable-web-security",
		},
		Prefs: map[string]interface{}{
			"profile.default_content_setting_values.notifications": 2,
		},
	}

	caps.AddChrome(chromeCaps)

	ops := []selenium.ServiceOption{}
	// Set up ChromeDriver service
	service, err := selenium.NewChromeDriverService("/usr/local/projects/aih/release/chromedriver", 8083, ops...)
	if err != nil {
		log.Fatal(err)
	}
	defer service.Stop()

	// Create a WebDriver object for Chrome browser
	webDriver, err := selenium.NewRemote(caps, "http://127.0.0.1:8083/wd/hub")
	if err != nil {
		log.Fatal(err)
	}

	// Maximize the window for keep the stream alive
	if err := webDriver.MaximizeWindow(""); err != nil {
		log.Fatal(err)
	}
	defer webDriver.Quit()
	// defer webDriver.Quit()

	if err := webDriver.Get("https://playphrase.me/"); err != nil {
		log.Fatal(err)
	}

	// Wait for the page to load
	if err := webDriver.SetImplicitWaitTimeout(10 * time.Second); err != nil {
		log.Fatal(err)
	}

	// List of words to search
	//words := []string{"Fingerstyle", "Fingerpaint", "Fingerling potatoes", "Finger food", "Finger cymbals"}

	// Click the body to start play video
	if body, err := webDriver.FindElement(selenium.ByXPATH, "//body"); err != nil {
		log.Fatal(err)
	} else {
		if err := body.Click(); err != nil {
			log.Fatal(err)
		}
	}

	//fmt.Println("clicked start")

	search_bar, err := webDriver.FindElement(selenium.ByID, "search-input")
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println("find search bar")

	for _, phrase := range words {

		for {
			value, err := search_bar.GetAttribute("value")
			if err != nil {
				log.Fatal(err)
			}
			if value == phrase {
				break
			}
			if elems, err := webDriver.FindElements(selenium.ByXPATH, "//i[contains(text(),'close')]"); err != nil {
				log.Fatal(err)
			} else if len(elems) > 0 {
				if err := elems[0].Click(); err != nil {
					log.Fatal(err)
				}
			}
			if err := search_bar.SendKeys(phrase); err != nil {
				log.Fatal(err)
			}
		}

		//fmt.Println("searching", phrase)

		time.Sleep(2 * time.Second)
		search_result_count, err := webDriver.FindElement(selenium.ByXPATH, "//li/div[@class='search-result-count']")
		if err != nil {
			log.Fatal(err)
		}
		text, err := search_result_count.Text()
		if err != nil {
			log.Fatal(err)
		}
		if text == "1/0" {
			//fmt.Println("Find nothing. Next!")
			continue
		}

		ch := make(chan bool)
		go func() {
			// Check if the page source contains the message "If you are not a sponsor you have a limit on our site."
			for {
				content, err := webDriver.PageSource()
				if err != nil {
					log.Fatal(err)
				}
				if strings.Contains(content, "If you are not a sponsor you have a limit on our site.") {
					//fmt.Println("Played 5 already. Next!")
					close(ch)
					break
				}
			}

		}()

		<-ch

		// Check if the search result count is "1/0"

	}
	fmt.Println("Clips play finished")

	// Close the Chrome browser
	if err := webDriver.Quit(); err != nil {
		log.Fatal(err)
	}

}
