package main

import (
    "github.com/sclevine/agouti"
    "log"
    "time"
)

/*
$ go install
$ crontab -e
    # Launch phantom program to play sunday night slow jams
    # Start five minutes before the show begins to ensure the
    # audio is playing by the time the consuming program is initialized
    # https://crontab.guru/#55_21_*_*_SUN
    55 21 * * SUN phantom
*/
func main() {
    // https://agouti.org
    driver := agouti.Selenium()
    if err := driver.Start(); err != nil {
        log.Printf("Failed to start Selenium: %s", err)
    }
    page, err := driver.NewPage(agouti.Browser("chrome"))
    if err != nil {
        log.Printf("Failed to open page: %s", err)
    }
    if err := page.Navigate("https://www.iheart.com/live/z100-portland-1961/"); err != nil {
        log.Printf("Failed to navigate: %s", err)
    }
    loginURL, err := page.URL()
    if err != nil {
        log.Printf("Failed to get page URL: %s", err)
    }
    log.Println(loginURL)
    arguments := map[string]interface{}{}
    var result string
    // https://github.com/sclevine/agouti/blob/6ada53bb069e86f8baf0953d4f0ddac081bc7610/internal/integration/page_test.go#L54
    page.RunScript(`
    var jq = document.createElement('script');
    jq.src = "https://ajax.googleapis.com/ajax/libs/jquery/2.1.4/jquery.min.js";
    document.getElementsByTagName('head')[0].appendChild(jq);
    return jQuery.noConflict();
        `, arguments, &result)
    time.Sleep(5 * time.Second)
    page.RunScript("$(\".css-wwzdid\").click();", arguments, &result)
    // Let the slow jamz begin
    time.Sleep(4*time.Hour + 5*time.Minute)
    if err := driver.Stop(); err != nil {
        log.Printf("Failed to close pages and stop WebDriver: %s", err)
    }
    log.Println("Music of the night")
}
