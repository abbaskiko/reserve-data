package main

import (
	"flag"
	"log"
	"strings"

	"github.com/boltdb/bolt"

	"github.com/KyberNetwork/reserve-data/settings"
	"github.com/KyberNetwork/reserve-data/settings/storage"
)

func main() {
	var settingDB string
	var token string
	var exchanges string
	flag.StringVar(&settingDB, "setting-db", "settings.db", "path to setting db file")
	flag.StringVar(&token, "token", "", "token to unlist")
	flag.StringVar(&exchanges, "exchanges", "binance,huobi,stable_exchange", "a list of exchanges to unlist")

	flag.Parse()

	if token == "" {
		log.Println("need a token to process")
		bs, err := storage.NewBoltSettingStorage(settingDB)
		if err != nil {
			log.Panicln(err)
		}
		tokens, err := bs.GetAllTokens()
		if err != nil {
			log.Panicln(err)
		}
		log.Println("follow is list of all NOT internal token, can consider to be unlist")
		log.Printf("ID\tName\tAddress\tIsActive\n")
		for _, tk := range tokens {
			if !tk.Internal {
				log.Printf("%s\t%s\t%s\t%v\n", tk.ID, tk.Name, tk.Address, tk.Active)
			}
		}
		return
	}
	var err error
	var db *bolt.DB
	db, err = bolt.Open(settingDB, 0600, nil)
	if err != nil {
		log.Panicln(err)
	}
	defer func() {
		_ = db.Close()
	}()
	var exs []settings.ExchangeName
	for _, e := range strings.Split(exchanges, ",") {
		exh, ok := settings.ExchangeTypeValues()[e]
		if ok {
			exs = append(exs, exh)
		} else {
			log.Panicf("exchange not supported %s", e)
		}
	}

	err = db.Update(func(tx *bolt.Tx) error {
		return storage.RemoveTokensFromExchanges(tx, []string{token}, exs)
	})
	if err != nil {
		log.Panicln(err)
	}
	log.Println("unlist token", token, "successfully")
}
