package cmd

// import (
// 	"log"

// 	api "github.com/segfault42/binance-api"
// 	"gopkg.in/tucnak/telebot.v2"
// 	tb "gopkg.in/tucnak/telebot.v2"
// )

// func notifNewFiatDeposit(client api.ApiInfo, b *telebot.Bot, m *tb.Message) {
// 	res, err := client.GetAccountService()
// 	if err != nil {
// 		log.Println("notifNewFiatDeposit	()", err)
// 		return
// 	}

// 	for _, elem := range res.Balances {
// 		if elem.Asset == "EUR"
// 	}

// 	db, err := connectToDB()
// 	defer db.Close()

// 	err = createTable(db, "fiat_deposit_history")
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}

// 	InsertToDB(db, )
// 	// for _, elem := range res.Balances {
// 	// 	if elem.Asset == "EUR" {
// 	// 		b.Send(m.Sender, msg)
// 	// 	}
// 	// }
// }
