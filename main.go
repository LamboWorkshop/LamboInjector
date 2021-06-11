package main

import (
	"log"
	"os"
	"time"

	"lamboInjector/cmd"

	"github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
)

const (
	order   = 0
	twitter = 1
)

type sChan struct {
	stopChan    chan struct{}
	stoppedChan chan struct{}
}

func main() {

	orderNotif := false
	// twitterNotif := false
	var notif [2]sChan

	notif[order].stopChan = make(chan struct{})
	notif[order].stoppedChan = make(chan struct{})

	f, client, err := cmd.InitServices()
	if err != nil {
		log.Println("initServices(): ", err)
		return
	}
	defer f.Close()

	b, err := tb.NewBot(tb.Settings{
		Token:  os.Getenv("LAMBO_TELEGRAM_TOKEN"),
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Println(err)
		return
	}

	b.Handle("/order_notifications_on", func(m *tb.Message) {
		if orderNotif == false {
			orderNotif = true
			go cmd.OrderNotifications(client, b, m, notif[order].stoppedChan, notif[order].stopChan)
			b.Send(m.Sender, "🔔 Order notifications on")
		} else {
			b.Send(m.Sender, "Order notifications already turned on !")
		}
	})

	b.Handle("/order_notifications_off", func(m *tb.Message) {
		if orderNotif == true {
			orderNotif = false
			close(notif[order].stopChan)
			<-notif[order].stoppedChan
			b.Send(m.Sender, "🔕 Order notifications off")
		} else {
			b.Send(m.Sender, "Order notifications already turned off !")
		}
	})

	b.Handle("/balance", func(m *tb.Message) {
		log.Println(m.Sender.Username + " typed : " + m.Text)
		balances, err := cmd.GetTotalBalanceInDollar(client)
		if err != nil {
			log.Println("getTotalBalanceInDollar() error : ", err.Error())
			b.Send(m.Sender, "Oops ! Something went wrong !")
		} else {
			b.Send(m.Sender, cmd.FormatMessage(balances))
		}
	})

	// b.Handle("twitter_notifications_on", func(m *tb.Message) {
	// 	if twitterNotif == false {
	// 		twitterNotif = true
	// 		go cmd.TwitterNotifications(b, m, notif[twitter].stoppedChan, notif[twitter].stopChan)
	// 		b.Send(m.Sender, "🔔 Twitter notifications on")
	// 	} else {
	// 		b.Send(m.Sender, "Twitter notifications already turned on !")
	// 	}
	// })

	b.Handle("/lamboTrunk", func(m *tb.Message) {
		log.Println(m.Sender.Username + " typed : " + m.Text)

		res, err := cmd.LamboTrunk(m)
		if err != nil {
			logrus.Error("lamboTrunk() error : ", err.Error())
			b.Send(m.Sender, "Usage :\n/lamboTrunk botstatus on\n/lamboTrunk botstatus off\n/lamboTrunk ETH/BUSD on 100")
		} else {
			b.Send(m.Sender, res)
		}
	})

	b.Handle("/help", func(m *tb.Message) {
		log.Println(m.Sender.Username + " typed : " + m.Text)
		b.Send(m.Sender, `Available commands :\n\n
		- /balance (This command print all crypto available on your binance account)\n
		- /order_notifications_on (Turn on notiications for the oder list)\n
		- /order_notifications_off (Turn off notifications for the order list)
		- /lamboTrunk botstatus on/off
		- /lamboTrunk ETH/BUSD On 100
		`)
	})

	log.Println("🏎 lamboInjector started !")

	b.Start()
}
