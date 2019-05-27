/*
Levis-MacBook-Pro:phantom levischoen$ which selenium-server
/usr/local/bin/selenium-server
$ go install
$ crontab -e
    # https://stackoverflow.com/questions/10129381/crontab-path-and-user
    SHELL=/bin/sh
    PATH=/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin:/Users/levischoen/go/bin/:/Users/levischoen/.rvm/bin
    # Launch phantom program to play sunday night slow jams
    # Start five minutes before the show begins to ensure the
    # audio is playing by the time the consuming program is initialized
    # https://crontab.guru/#55_21_*_*_SUN
    55 21 * * SUN phantom
*/
package main

import (
    "github.com/sclevine/agouti"
    "log"
    "time"
)

const (
    headStartMinutes   = 5
    slowJamzSeconds    = 14400 // 60 seconds/minute * 60 minutes/hour * 4 hours
    logIntervalSeconds = 30
)

func TailLogs(page *agouti.Page) {
    for {
        logTypes, err := page.LogTypes()
        if err != nil {
            log.Printf("error %s retrieving page %+v logtypes\n", err, page)
        }
        log.Printf("available log types: %v\n", logTypes)
        for _, logType := range logTypes {
            logs, err := page.ReadNewLogs(logType)
            if err != nil {
                log.Printf("error %s retrieving page %+v %s logs\n", err, page, logType)
            }
            log.Printf("%s logs: %v\n", logType, logs)
        }
        time.Sleep(logIntervalSeconds * time.Second)
    }
}

func main() {
    timeout := agouti.Timeout(slowJamzSeconds)
    // https://agouti.org
    driver := agouti.Selenium(timeout)
    if err := driver.Start(); err != nil {
        log.Printf("Failed to start Selenium: %s", err)
    }
    page, err := driver.NewPage(agouti.Browser("chrome"), timeout)
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
    // Do some busy work to ensure the session doesn't prematurely timeout
    go TailLogs(page)
    // Let the slow jamz begin
    time.Sleep(slowJamzSeconds*time.Second + headStartMinutes*time.Minute)
    if err := driver.Stop(); err != nil {
        log.Printf("Failed to close pages and stop WebDriver: %s", err)
    }
    log.Println("Music of the night")
}
