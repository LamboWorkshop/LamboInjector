package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/dariubs/percent"
	api "github.com/segfault42/binance-api"
	"gopkg.in/tucnak/telebot.v2"
	tb "gopkg.in/tucnak/telebot.v2"
)

// For /balance command
type coin struct {
	Total float64
	Asset string
}

// For daily benefit
type SCoins struct {
	Asset            string  `json:"asset"`
	Total            float64 `json:"total"`
	PercentVariation float64 `json:"percentVariation"`
}
type SBenefit struct {
	Coins                 []SCoins `json:"coins"`
	Total                 float64  `json:"total"`
	TotalPercentVariation float64  `json:totalPercentVariation`
}

func GetTotalBalanceInDollar(client api.ApiInfo) ([]coin, error) {

	var balances []coin
	// totalBalanceInDollar := 0.0

	res, err := client.GetBalances()
	if err != nil {
		return balances, err
	}

	for _, elem := range res {
		var tmp coin

		res1, _ := strconv.ParseFloat(elem.Free, 64)
		res2, _ := strconv.ParseFloat(elem.Locked, 64)

		if res1 > 0 || res2 > 0 {
			total := res1 + res2
			tmp.Asset = elem.Asset
			tmp.Total = total
			balances = append(balances, tmp)
		}
	}

	for i, elem := range balances {
		price := ""
		if elem.Asset == "USDT" {
			price, err = client.GetTickerPrice("BUSD" + elem.Asset)
		} else {
			price, err = client.GetTickerPrice(elem.Asset + "USDT")
		}
		if err != nil {
			return balances, err
		} else {
			floatPrice, _ := strconv.ParseFloat(price, 64)
			floatPrice *= elem.Total
			balances[i].Total = floatPrice
		}
	}

	return balances, nil
}

func FormatMessage(balances []coin) string {
	total := 0.0
	message := ""

	for _, elem := range balances {
		// Don't print 0 balance
		if elem.Total >= 1.0 {
			message += "💵 " + elem.Asset + " : " + fmt.Sprintf("%.2f$", elem.Total) + "\n"
			total += elem.Total
		}
	}

	message += "\n💰 Total : " + fmt.Sprintf("%.2f$", total)

	return message
}

func dailyBenefit(api api.ApiInfo, yesterdayBenefit SBenefit) (SBenefit, error) {

	var todayBenef SBenefit

	balance, err := GetTotalBalanceInDollar(api)
	if err != nil {
		return todayBenef, err
	}

	for _, elem := range balance {
		var tmp SCoins

		tmp.Asset = elem.Asset
		tmp.Total = elem.Total
		for _, elem2 := range yesterdayBenefit.Coins {
			if elem2.Asset == elem.Asset {
				tmp.PercentVariation = percent.ChangeFloat(elem2.Total, elem.Total)
				break
			}
		}
		todayBenef.Coins = append(todayBenef.Coins, tmp)
		todayBenef.Total += elem.Total
	}
	todayBenef.TotalPercentVariation = percent.ChangeFloat(yesterdayBenefit.Total, todayBenef.Total)

	return todayBenef, nil

}

func getYesterdayBenefit() (SBenefit, error) {
	var result SBenefit

	jsonFile, err := os.Open("benefit.json")
	if err != nil {
		return result, err
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return result, err
	}

	json.Unmarshal([]byte(byteValue), &result)

	return result, nil
}

func notifDailyBenefit(api api.ApiInfo, b *telebot.Bot, m *tb.Message) error {

	currentTime := time.Now()

	timeStampString := currentTime.Format("2006-01-02 15:04:05")
	layOut := "2006-01-02 15:04:05"
	timeStamp, err := time.Parse(layOut, timeStampString)
	if err != nil {
		return err
	}
	hr, min, sec := timeStamp.Local().Clock()
	if hr == 00 && min == 00 && sec < 10 {
		yb, err := getYesterdayBenefit()
		if err != nil {
			return err
		}

		todayBenef, err := dailyBenefit(api, yb)
		if err != nil {
			return err
		}

		writeJsonToFile(todayBenef)
		message := formatDailyBenefit(todayBenef, yb)
		b.Send(m.Sender, message)
	}

	return nil
}

func writeJsonToFile(todayBenef SBenefit) error {
	file, err := json.MarshalIndent(todayBenef, "", " ")

	if err != nil {
		return err
	}
	err = ioutil.WriteFile("benefit.json", file, 0644)

	return err
}

func formatDailyBenefit(todayBenef, yesterdayBenefit SBenefit) string {
	message := "🎰 Benefit of the day :\n\n"

	for _, elem := range todayBenef.Coins {
		if elem.Total > 1.0 {
			message += "💵 " + elem.Asset + " : " + fmt.Sprintf("%.2f", elem.Total) + "$"
			for _, elem2 := range yesterdayBenefit.Coins {
				if elem.Asset == elem2.Asset {
					if elem.Total > elem2.Total {
						message += fmt.Sprintf(", +%.2f%% ", elem.PercentVariation) + "⬆️"
					} else if elem.Total < elem2.Total {
						message += fmt.Sprintf(", %.2f%% ", elem.PercentVariation) + "⬇️"
					} else {
						message += fmt.Sprintf(", %.2f%% ", elem.PercentVariation) + "➡️"
					}
				}
			}
			message += "\n"
		}
	}

	return message
}
