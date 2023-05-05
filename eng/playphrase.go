package eng

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
        "github.com/chromedp/cdproto/input"
	"github.com/chromedp/cdproto/cdp"
	"log"
	"strings"
	"time"
)

func Play(words []string) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Set up Chrome options
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		// chromedp.Flag("disable-gpu", true),
		// chromedp.Flag("no-sandbox", true),
		// chromedp.Flag("disable-infobars", true),
		// chromedp.Flag("disable-extensions", true),
		// chromedp.Flag("disable-web-security", true),
		// chromedp.Flag("mute-audio", true),
		chromedp.Flag("mute-audio", false),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	// Create a WebDriver object for Chrome browser
	ctx, cancel = chromedp.NewContext(allocCtx)
	defer cancel()

	// Maximize the window for keep the stream alive
	// This function has been removed as it returns an error, and it is not being used anywhere else
	// if err := chromedp.Run(ctx, chromedp.MaximizeWindow()); err != nil {
	//	log.Fatal(err)
	// }

	if err := chromedp.Run(ctx, chromedp.Navigate("https://playphrase.me/")); err != nil {
		fmt.Println(err)
	}

	if err := chromedp.Run(ctx, chromedp.WaitVisible(`//i[@class="material-icons-outlined"]`, chromedp.BySearch)); err != nil {
		fmt.Println(err)
	}
	//   chromedp.WaitVisible("body"),
	//chromedp.Sleep(2 * time.Second),
	if err := chromedp.Run(ctx, chromedp.Click("//body", chromedp.BySearch)); err != nil {
		fmt.Println(err)
	}

	search_bar := "#search-input"

	for _, phrase := range words {

		for {
			var value string
			if err := chromedp.Run(ctx, chromedp.Value(search_bar, &value)); err != nil {
				log.Fatal(err)
			}
			if value == phrase {
				break
			}
			if err := chromedp.Run(ctx, chromedp.Click(`//i[contains(text(),"close")]`, chromedp.BySearch)); err != nil {
				log.Fatal(err)
			}

		//	if err := chromedp.Run(ctx, chromedp.SendKeys(search_bar, phrase)); err != nil {
		//		log.Fatal(err)
		//	}
                        var nodes []*cdp.Node
			_ = chromedp.Run(ctx, chromedp.Nodes("#search-input", &nodes, chromedp.ByQuery))
			_ = chromedp.Run(ctx, chromedp.MouseClickNode(nodes[0]))
			_ = chromedp.Run(ctx, input.InsertText(phrase))


			//if err := chromedp.Run(ctx, chromedp.Sleep(1*time.Second)); err != nil {
			//	log.Fatal(err)
			//}
		}

		// Check if the search result count is "1/0"
		if err := chromedp.Run(ctx, chromedp.Sleep(5*time.Second)); err != nil {
			log.Fatal(err)
		}
		var search_result_count string
		if err := chromedp.Run(ctx, chromedp.Text("li div.search-result-count", &search_result_count)); err != nil {
			log.Fatal(err)
		}
		if search_result_count == "1/0" {
			//fmt.Println("Find nothing. Next!")
			continue
		}

		ch := make(chan bool)
		go func() {
			// Check if the page source contains the message "If you are not a sponsor you have a limit on our site."
			for {
				var content string
				if err := chromedp.Run(ctx, chromedp.InnerHTML("body", &content)); err != nil {
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

	}

	// Close the Chrome browser
	// Shutdown() function is not present in chromedp.Context. Using Stop() function instead.
	if err := chromedp.Run(ctx, chromedp.Stop()); err != nil {
		log.Fatal(err)
	}

}
