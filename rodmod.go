package main

import (
//	"crypto/md5"
	"fmt"
	"strings"
	"time"
	"io/ioutil"
	"github.com/tidwall/gjson"

	"github.com/go-rod/rod"
	"github.com/go-rod/stealth"
	"github.com/go-rod/rod/lib/utils"
	"github.com/go-rod/rod/lib/launcher"
)

//func init() {
//	launcher.NewBrowser().MustGet()
//}
//
func main() {

	// Read cookie
	//chatgpt_json, err := ioutil.ReadFile("./4.json")
	//if err != nil {
	//	err = ioutil.WriteFile("./4.json", []byte(""), 0644)
	//}
	//var chatgpt_js string
	//chatgpt_js = gjson.Parse(string(chatgpt_json)).String()

	// Read user/password
	u_json, _ := ioutil.ReadFile("user.json")
	var chatgpt_user string
	var chatgpt_password string
	chatgpt_user = gjson.Get(string(u_json), "chatgpt.user").String()
	chatgpt_password = gjson.Get(string(u_json), "chatgpt.password").String()
	fmt.Println(chatgpt_user)
	fmt.Println(chatgpt_password)
        //time.Sleep(1000 * time.Second)


	//proxy_url := launcher.New().Proxy("socks5://127.0.0.1:7890").Delete("use-mock-keychain").MustLaunch()
	proxy_url := launcher.New().Proxy("socks5://127.0.0.1:7890").MustLaunch()

	//browser := rod.New().
	//           Trace(true).
	//	   ControlURL(proxy_url).
	//           //Timeout(time.Minute).
	//           Timeout(3 * time.Minute).
	//	   MustConnect()
        var browser *rod.Browser
	if "1" != "1" {
	browser = rod.New().
	           Trace(true).
		   ControlURL(proxy_url).
	           Timeout(3 * time.Minute).
		   MustConnect()
		  } else {
	browser = rod.New().
	           Trace(true).
	           Timeout(3 * time.Minute).
		   MustConnect()

		  }
	//defer browser.MustClose()

	// You can also use stealth.JS directly without rod
	//fmt.Printf("js: %x\n\n", md5.Sum([]byte(stealth.JS)))

	// Read cookie
	//chatgpt_json, err := ioutil.ReadFile("cookies/chatgpt.json")
	//if err != nil {
	//	err = ioutil.WriteFile("cookies/chatgpt.json", []byte(""), 0644)
	//}
	//var chatgptjs string
	//chatgptjs = gjson.Parse(string(chatgpt_json)).String()

	page := stealth.MustPage(browser)
	//page.Call("Network.setCookies", cdp.Object{ "cookies": []cdp.Object{{

	//page.MustNavigate("https://bot.sannysoft.com")
	page.MustNavigate("https://chat.openai.com")

	page.MustElementX("//div[contains(text(), 'Welcome to ChatGPT')] | //h2[contains(text(), 'Get started')]").MustWaitVisible()
	page.MustElementX("//div[not(contains(@class, 'mb-4')) and contains(text(), 'Log in')]").MustClick()
	utils.Sleep(1.5)
	page.MustElementX("//input[@id='username']").MustWaitVisible().MustInput(chatgpt_user)
	utils.Sleep(1.5)
	page.MustElementX("//button[contains(text(), 'Continue')]").MustClick()
	utils.Sleep(1.5)
	page.MustElementX("//input[@id='password']").MustWaitVisible().MustInput(chatgpt_password)
	utils.Sleep(1.5)
	//page.MustElementX("//input[@id='password']").MustInput(chatgpt_password)
	page.MustElementX("//button[not(contains(@aria-hidden, 'true')) and contains(text(), 'Continue')]").MustClick()
	page.MustElementX("//h4[contains(text(), 'This is a free research preview.')]").MustWaitVisible()
	utils.Sleep(1.5)
	page.MustElementX("//button/div[contains(text(), 'Next')]").MustClick()
	page.MustElementX("//h4[contains(text(), 'How we collect data')]").MustWaitVisible()
	utils.Sleep(1.5)
	page.MustElementX("//button/div[contains(text(), 'Next')]").MustClick()
	page.MustElementX("//h4[contains(text(), 'love your feedback!')]").MustWaitVisible()
	utils.Sleep(1.5)
	page.MustElementX("//button/div[contains(text(), 'Done')]").MustClick()
	utils.Sleep(1.5)
	page.MustElementX("//a[contains(text(), 'New chat')]").MustWaitVisible().MustClick()
	page.MustElementX("//textarea[@id='prompt-textarea']").MustWaitVisible()
	utils.Sleep(1.5)
	page.MustElementX("//textarea[@id='prompt-textarea']").MustInput("hello")
	utils.Sleep(1.5)
	//utils.Pause()
	sends := page.Timeout(200 * time.Second).MustElements("button:last-of-type svg path[d='M.5 1.163A1 1 0 0 1 1.97.28l12.868 6.837a1 1 0 0 1 0 1.766L1.969 15.72A1 1 0 0 1 .5 14.836V10.33a1 1 0 0 1 .816-.983L8.5 8 1.316 6.653A1 1 0 0 1 .5 5.67V1.163Z']")
	sends[len(sends)-1].MustClick()
        page.Timeout(20000 * time.Second).MustElement("svg:last-of-type path[d='M20.49 9A9 9 0 0 0 5.64 5.64L1 10m22 4l-4.64 4.36A9 9 0 0 1 3.51 15']").MustWaitVisible()
	//page.MustScreenshot("")
	//page.MustScreenshot("")
	fmt.Println("Retry icon show")
	content := page.MustElementX("(//div[contains(@class, 'group w-full')])[last()]").MustText()
	fmt.Println(content)
	page.MustScreenshot("")
	utils.Pause()
//	utils.Pause()
	//printReport(page)

	/*
		Output:

		js: 173d23e3db48bf47441b2f4735bbc631

		User Agent (Old): true

		WebDriver (New): missing (passed)

		WebDriver Advanced: passed

		Chrome (New): present (passed)

		Permissions (New): prompt

		Plugins Length (Old): 3

		Plugins is of type PluginArray: passed

		Languages (Old): en-US,en

		WebGL Vendor: Intel Inc.

		WebGL Renderer: Intel Iris OpenGL Engine

		Broken Image Dimensions: 16x16
	*/
}

func printReport(page *rod.Page) {
	el := page.MustElement("#broken-image-dimensions.passed")
	for _, row := range el.MustParents("table").First().MustElements("tr:nth-child(n+2)") {
		cells := row.MustElements("td")
		key := cells[0].MustProperty("textContent")
		if strings.HasPrefix(key.String(), "User Agent") {
			fmt.Printf("\t\t%s: %t\n\n", key, !strings.Contains(cells[1].MustProperty("textContent").String(), "HeadlessChrome/"))
		} else if strings.HasPrefix(key.String(), "Hairline Feature") {
			// Detects support for hidpi/retina hairlines, which are CSS borders with less than 1px in width, for being physically 1px on hidpi screens.
			// Not all the machine suppports it.
			continue
		} else {
			fmt.Printf("\t\t%s: %s\n\n", key, cells[1].MustProperty("textContent"))
		}
	}

	page.MustScreenshot("")
}
