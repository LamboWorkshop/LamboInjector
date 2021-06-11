package cmd

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/adshao/go-binance/v2"
	api "github.com/segfault42/binance-api"
	"gopkg.in/tucnak/telebot.v2"
	tb "gopkg.in/tucnak/telebot.v2"
)

const (
	notifCloseOrderTime = 1 * time.Minute
)

func difference(slice1, slice2 []*binance.Order) []*binance.Order {
	var diff []*binance.Order

	for i := 0; i < 2; i++ {
		for _, s1 := range slice1 {
			found := false
			for _, s2 := range slice2 {
				if *s1 == *s2 {
					found = true
					break
				}
			}
			// String not found. We add it to return slice
			if !found {
				diff = append(diff, s1)
			}
		}
		// Swap the slices, only if it was the first loop
		if i == 0 {
			slice1, slice2 = slice2, slice1
		}
	}

	return diff
}

func OrderNotifications(api api.ApiInfo, b *telebot.Bot, m *tb.Message, stoppedchan, stopchan chan struct{}) {
	defer close(stoppedchan)
	for {
		select {
		default:
			if err := notifDailyBenefit(api, b, m); err != nil {
				log.Println("notifDailyBenefit() :", err)
			}
			// notifNewFiatDeposit(api, b, m)
			notifCloseOrder(api, b, m)
		case <-stopchan:
			return
		}
	}
}

func notifCloseOrder(api api.ApiInfo, b *telebot.Bot, m *tb.Message) {
	orders, err := api.GetOpenOrders()
	if err != nil {
		log.Println("GetOpenOrders() :", err)
		time.Sleep(notifCloseOrderTime)
		return
	}

	time.Sleep(notifCloseOrderTime)

	updateOrders, err := api.GetOpenOrders()
	if err != nil {
		log.Println("GetOpenOrders() :", err)
		time.Sleep(notifCloseOrderTime)
		return
	}

	diff := difference(orders, updateOrders)

	if len(diff) > 0 {
		for _, elem := range diff {
			if elem.Status == binance.OrderStatusTypeFilled ||
				elem.Status == binance.OrderStatusTypeNew ||
				elem.Status == binance.OrderStatusTypePartiallyFilled {

				price, _ := strconv.ParseFloat(elem.Price, 64)
				quantity, _ := strconv.ParseFloat(elem.OrigQuantity, 64)
				status := ""

				if len(updateOrders) > len(orders) && elem.Status == binance.OrderStatusTypeNew {
					status = "➕ New"
				} else if len(updateOrders) < len(orders) && elem.Status == binance.OrderStatusTypeNew {
					status = "❌ Cancelled"
				} else {
					status = "🎰 "

					//TODO : Sell
					log.Println(&updateOrders)
					log.Println(&orders)
					quantity, _ = strconv.ParseFloat(elem.ExecutedQuantity, 64)
					status = string(elem.Status)
				}
				msg := "⚠️ Order Update :\n\n" + status + " " + string(elem.Side) + " order : " + elem.Symbol + "\n" +
					"🥞 Quantity : " + fmt.Sprintf("%f", quantity) + "\n" +
					"🏷 Price : " + fmt.Sprintf("%f", price) + "$\n" +
					"💰 Total : " + fmt.Sprintf("%f", quantity*price) + "$"

				b.Send(m.Sender, msg)
			}
		}
	}
}

func TwitterNotifications(b *telebot.Bot, m *tb.Message, stoppedchan, stopchan chan struct{}) {

}
