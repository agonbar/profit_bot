package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	bolt "go.etcd.io/bbolt"

	tb "gopkg.in/tucnak/telebot.v2"

	"github.com/jedib0t/go-pretty/table"
)

func initTG(db *bolt.DB) *tb.Bot {
	// INIT TELEGRAM BOT
	b, err := tb.NewBot(tb.Settings{
		Token:  os.Getenv("TOKEN"),
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
	}

	// HANDLE REQUESTS
	b.Handle("/space", func(m *tb.Message) {
		var nodes = getNodes(db)
		var res string
		for i := 0; i < len(nodes); i++ {
			res += "\n" + nodes[i][0] + ": " + getSpace(nodes[i][1])
		}
		b.Send(m.Sender, res)
	})

	b.Handle("/prices", func(m *tb.Message) {
		var nodes = getNodes(db)
		var res string
		for i := 0; i < len(nodes); i++ {
			res += "\n" + nodes[i][0] + ": " + getPrice(nodes[i][1])
		}
		_, err = b.Send(m.Sender, res)
		if err != nil {
			log.Println(err)
		}
	})

	b.Handle("/update", func(m *tb.Message) {
		updateResponse(db, b)
	})

	b.Handle("/screenshot", func(m *tb.Message) {
		var nodes = getNodes(db)
		var pictures []tb.InputMedia
		for i := 0; i < len(nodes); i++ {
			ctx, cancel := chromedp.NewContext(context.Background())
			defer cancel()
			var buf []byte
			if err := chromedp.Run(ctx, fullScreenshot(nodes[i][1], 90, &buf)); err != nil {
				log.Fatal(err)
			}
			p := &tb.Photo{File: tb.FromReader(bytes.NewReader(buf))}
			// Will upload the file from disk and send it to recipient
			pictures = append(pictures, p)
		}
		b.SendAlbum(m.Sender, pictures)
	})

	b.Handle("/subscribe", func(m *tb.Message) {
		saveUser(db, m.Sender.ID, "1")
	})

	b.Handle("/unsubscribe", func(m *tb.Message) {
		saveUser(db, m.Sender.ID, "0")
	})

	b.Handle("/newNode", func(m *tb.Message) {
		strParts := strings.Split(m.Text, " ")
		if len(strParts) == 3 {
			saveNode(db, strParts[1], strParts[2])
			b.Send(m.Sender, "Created! name: "+strParts[1]+" url: "+strParts[2])
		} else {
			b.Send(m.Sender, "Please specify a node separated by space and its url, example: /newNode test https://newnode.myurl.com:14002/")
		}
	})

	b.Handle("/delNode", func(m *tb.Message) {
		strParts := strings.Split(m.Text, " ")
		if len(strParts) == 2 {
			deleteNode(db, strParts[1])
			b.Send(m.Sender, "Deleted! name: "+strParts[1])
		} else {
			b.Send(m.Sender, "Please specify a node name separated by space, example: /delNode test")
		}
	})

	return b
}

func updateResponse(db *bolt.DB, b *tb.Bot) {
	t := table.NewWriter()
	t.AppendHeader(table.Row{"NODE", "SPACE", "PROFIT"})
	t.SetStyle(table.StyleLight)

	var users [][2]string = getUsers(db)
	var nodes = getNodes(db)
	var total [3]float32
	for i := 0; i < len(nodes); i++ {
		var node [3]string = getStatus(nodes[i][1])
		t.AppendRow([]interface{}{nodes[i][0], node[0] + "/" + node[1], node[2]})
		if size, err := strconv.ParseFloat(strings.Replace(node[0], "TB", "", -1), 32); err == nil {
			total[0] += float32(size)
		}
		if totalSize, err := strconv.ParseFloat(strings.Replace(node[1], "TB", "", -1), 32); err == nil {
			total[1] += float32(totalSize)
		}
		if price, err := strconv.ParseFloat(strings.Replace(node[2], "$", "", -1), 32); err == nil {
			total[2] += float32(price)
		}
	}
	t.AppendRow([]interface{}{"Total", fmt.Sprintf("%.2fTB/%.2fTB", total[0], total[1]), fmt.Sprintf("$%.2f", total[2])})
	for i := 0; i < len(users); i++ {
		if users[i][1] == "1" {
			var id, _ = strconv.Atoi(users[i][0])
			var _, err = b.Send(tb.ChatID(id), "<code>\n"+t.Render()+"\n</code>", tb.ModeHTML)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
