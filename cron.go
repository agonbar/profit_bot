package main

import (
	"github.com/robfig/cron/v3"
	bolt "go.etcd.io/bbolt"
	tb "gopkg.in/tucnak/telebot.v2"
)

func initCron(db *bolt.DB, tb *tb.Bot) {
	c := cron.New()
	c.AddFunc("0 20 * * *", func() { updateResponse(db, tb) })
	c.Start()
}
