package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/mmcdole/gofeed"
)

func init() {
	InitDB("./data/assembly.db")
}

type config struct {
	assemblyURL string
	jobPrefix   []string
	webhookURL  string
}

func main() {
	fmt.Println("Read config.")
	cfg, err := getConfig()
	if err != nil {
		log.Panic(err)
	}

	ch := make(chan *gofeed.Item)

	slack := &Slack{cfg.webhookURL}

	start := make(chan bool)
	stop := make(chan bool)

	fmt.Println("Update list of assemblies.")
	UpdateAssemblyList(cfg.assemblyURL, cfg.jobPrefix)

	go func() {
		start <- true
	}()

	upd := time.Tick(24 * time.Hour)
	for {
		select {
		case <-upd:
			fmt.Println("Stop listen rss.")
			stop <- true

			fmt.Println("Update list of assemblies.")
			UpdateAssemblyList(cfg.assemblyURL, cfg.jobPrefix)

			fmt.Println("Start listen rss.")
			go func() {
				start <- true
			}()

		case <-start:
			for _, url := range GetActualAssembly() {
				fmt.Printf("Start listen %s\n", url)
				go pollFeed(fmt.Sprintf("%srssAll", url), ch, stop)
			}

		case item := <-ch:
			if err := slack.Post(item); err != nil {
				log.Printf("could not post message to slack %v: %v", item, err)
			} else {
				_, err = db.Exec("INSERT INTO entries(url) VALUES(?)", item.Link)
				if err != nil {
					log.Printf("could not insert to entries %s: %v", item.Link, err)
					log.Panic(err)
				}
			}
		}
	}
}

func getConfig() (cfg *config, err error) {
	buf, err := ioutil.ReadFile("./data/config.json")
	if err != nil {
		return nil, fmt.Errorf("could not read config: %v", err)
	}

	cfg = &config{}
	err = json.Unmarshal(buf, cfg)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal config data: %v", err)
	}

	return cfg, nil
}

func pollFeed(url string, ch chan *gofeed.Item, stop chan bool) {
	fp := gofeed.NewParser()
	for {
		select {
		case <-stop:
			fmt.Printf("Stop listen %s\n", url)
			return
		default:
			fmt.Printf("Pull %s\n", url)
			feed, err := fp.ParseURL(url)
			if err != nil {
				log.Printf("could not parse rss feed %s: %v", url, err)
				return
			}

			itemsCount := len(feed.Items)
			if itemsCount > 0 {
				item := feed.Items[0]
				row := db.QueryRow("SELECT id FROM entries WHERE url = ?", item.Link)
				var id int64
				err = row.Scan(&id)
				if err == sql.ErrNoRows {
					ch <- item
				}
			}

			<-time.After(time.Duration(10 * time.Minute))
		}
	}
}
