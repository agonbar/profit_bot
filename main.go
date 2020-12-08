package main

func main() {

	//INIT DB
	var db = initDB()

	//INIT TELEGRAM BOT
	var tg = initTG(db)

	// INIT CRON JOBS
	initCron(db, tg)

	tg.Start()
}
