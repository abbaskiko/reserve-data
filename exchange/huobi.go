package exchange

import (
	"errors"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"strings"
	"sync"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/KyberNetwork/reserve-data/common"
	"github.com/KyberNetwork/reserve-data/common/blockchain"
	huobiblockchain "github.com/KyberNetwork/reserve-data/exchange/huobi/blockchain"
	huobihttp "github.com/KyberNetwork/reserve-data/exchange/huobi/http"
	commonv3 "github.com/KyberNetwork/reserve-data/v3/common"
	"github.com/KyberNetwork/reserve-data/v3/storage"
)

const (
	huobiEpsilon float64 = 0.0000000001 // 10e-10
)

type Huobi struct {
	interf     HuobiInterface
	blockchain HuobiBlockchain
	storage    HuobiStorage
	sr         storage.SettingReader
}

func (h *Huobi) TokenAddresses() (map[string]ethereum.Address, error) {
	result, err := h.sr.GetDepositAddresses(uint64(common.Huobi))
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (h *Huobi) MarshalText() (text []byte, err error) {
	return []byte(h.ID()), nil
}

// RealDepositAddress return the actual Huobi deposit address of a token
// It should only be used to send 2nd transaction.
func (h *Huobi) RealDepositAddress(tokenID string) (ethereum.Address, error) {
	liveAddress, err := h.interf.GetDepositAddress(tokenID)
	if err != nil || liveAddress.Address == "" {
		log.Printf("WARNING: Get Huobi live deposit address for token %s failed: (%v) or the replied address is empty. Check the currently available address instead", tokenID, err)
		addrs, uErr := h.sr.GetDepositAddresses(uint64(common.Huobi))
		if uErr != nil {
			return ethereum.Address{}, uErr
		}
		result, supported := addrs[tokenID]
		if !supported {
			return result, fmt.Errorf("real deposit address of token %s is not available", tokenID)
		}
		return result, nil
	}
	return ethereum.HexToAddress(liveAddress.Address), nil
}

// Address return the deposit address of a token in Huobi exchange.
// Due to the logic of Huobi exchange, every token if supported will be
// deposited to an Intermediator address instead.
func (h *Huobi) Address(token common.Token) (ethereum.Address, bool) {
	result := h.blockchain.GetIntermediatorAddr()
	_, err := h.RealDepositAddress(token.ID)
	//if the realDepositAddress can not be querried, that mean the token isn't supported on Huobi
	if err != nil {
		return result, false
	}
	return result, true
}

// GetLiveExchangeInfos querry the Exchange Endpoint for exchange precision and limit of a list of tokenPairIDs
// It return error if occurs.
func (h *Huobi) GetLiveExchangeInfos(tokenPairIDs []common.TokenPairID) (common.ExchangeInfo, error) {
	result := make(common.ExchangeInfo)
	exchangeInfo, err := h.interf.GetExchangeInfo()
	if err != nil {
		return result, err
	}
	for _, pairID := range tokenPairIDs {
		exchangePrecisionLimit, ok := h.getPrecisionLimitFromSymbols(pairID, exchangeInfo)
		if !ok {
			return result, fmt.Errorf("huobi Exchange Info reply doesn't contain token pair %s", string(pairID))
		}
		result[pairID] = exchangePrecisionLimit
	}
	return result, nil
}

// getPrecisionLimitFromSymbols find the pairID amongs symbols from exchanges,
// return ExchangePrecisionLimit of that pair and true if the pairID exist amongs symbols, false if otherwise
func (h *Huobi) getPrecisionLimitFromSymbols(pair common.TokenPairID, symbols HuobiExchangeInfo) (common.ExchangePrecisionLimit, bool) {
	var result common.ExchangePrecisionLimit
	pairName := strings.ToUpper(strings.Replace(string(pair), "-", "", 1))
	for _, symbol := range symbols.Data {
		symbolName := strings.ToUpper(symbol.Base + symbol.Quote)
		if symbolName == pairName {
			result.Precision.Amount = symbol.AmountPrecision
			result.Precision.Price = symbol.PricePrecision
			result.MinNotional = 0.02
			return result, true
		}
	}
	return result, false
}

// ID must return the exact string or else simulation will fail
func (h *Huobi) ID() common.ExchangeID {
	return common.ExchangeID(common.Huobi.String())
}

func (h *Huobi) TokenPairs() ([]commonv3.TradingPairSymbols, error) {
	pairs, err := h.sr.GetTradingPairSymbols(uint64(common.Binance))
	if err != nil {
		return nil, err
	}
	return pairs, nil
}

func (h *Huobi) Name() string {
	return "huobi"
}

func (h *Huobi) QueryOrder(symbol string, id uint64) (done float64, remaining float64, finished bool, err error) {
	result, err := h.interf.OrderStatus(symbol, id)
	if err != nil {
		return 0, 0, false, err
	}
	if result.Data.ExecutedQty != "" {
		done, err = strconv.ParseFloat(result.Data.ExecutedQty, 64)
		if err != nil {
			return 0, 0, false, err
		}
	}
	var total float64
	if result.Data.OrigQty != "" {
		total, err = strconv.ParseFloat(result.Data.OrigQty, 64)
		if err != nil {
			return 0, 0, false, err
		}
	}
	return done, total - done, total-done < huobiEpsilon, nil
}

func (h *Huobi) Trade(tradeType string, base common.Token, quote common.Token, rate float64, amount float64, timepoint uint64) (id string, done float64, remaining float64, finished bool, err error) {
	result, err := h.interf.Trade(tradeType, base, quote, rate, amount, timepoint)

	if err != nil {
		return "", 0, 0, false, err
	}
	var orderID uint64
	if result.OrderID != "" {
		orderID, err = strconv.ParseUint(result.OrderID, 10, 64)
		if err != nil {
			return "", 0, 0, false, err
		}
	}
	done, remaining, finished, err = h.QueryOrder(
		base.ID+quote.ID,
		orderID,
	)
	if err != nil {
		log.Printf("Huobi Query order error: %s", err.Error())
	}
	return result.OrderID, done, remaining, finished, err
}

//Withdraw return withdraw id from huobi
func (h *Huobi) Withdraw(token common.Token, amount *big.Int, address ethereum.Address, timepoint uint64) (string, error) {
	withdrawID, err := h.interf.Withdraw(token, amount, address)
	if err != nil {
		return "", err
	}
	return withdrawID, err
}

func (h *Huobi) CancelOrder(id, base, quote string) error {
	idNo, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return err
	}
	symbol := base + quote
	result, err := h.interf.CancelOrder(symbol, idNo)
	if err != nil {
		return err
	}
	if result.Status != "ok" {
		return errors.New("Huobi Couldn't cancel order id " + id)
	}
	return nil
}

func (h *Huobi) FetchOnePairData(
	wg *sync.WaitGroup,
	baseID, quoteID string,
	data *sync.Map,
	timepoint uint64) {

	defer wg.Done()
	result := common.ExchangePrice{}

	timestamp := common.Timestamp(fmt.Sprintf("%d", timepoint))
	result.Timestamp = timestamp
	result.Valid = true
	respData, err := h.interf.GetDepthOnePair(baseID, quoteID)
	returnTime := common.GetTimestamp()
	result.ReturnTime = returnTime
	if err != nil {
		result.Valid = false
		result.Error = err.Error()
	} else {
		if respData.Status != "ok" {
			result.Valid = false
		} else {
			for _, buy := range respData.Tick.Bids {
				quantity := buy[1]
				rate := buy[0]
				result.Bids = append(
					result.Bids,
					common.NewPriceEntry(
						quantity,
						rate,
					),
				)
			}
			for _, sell := range respData.Tick.Asks {
				quantity := sell[1]
				rate := sell[0]
				result.Asks = append(
					result.Asks,
					common.NewPriceEntry(
						quantity,
						rate,
					),
				)
			}
		}
	}
	data.Store(common.NewTokenPairID(baseID, quoteID), result)
}

func (h *Huobi) FetchPriceData(timepoint uint64) (map[common.TokenPairID]common.ExchangePrice, error) {
	wait := sync.WaitGroup{}
	data := sync.Map{}
	pairs, err := h.TokenPairs()
	if err != nil {
		return nil, err
	}
	for _, pair := range pairs {
		wait.Add(1)
		baseID, quoteID := pair.BaseSymbol, pair.QuoteSymbol
		go h.FetchOnePairData(&wait, baseID, quoteID, &data, timepoint)
	}
	wait.Wait()
	result := map[common.TokenPairID]common.ExchangePrice{}
	data.Range(func(key, value interface{}) bool {
		tokenPairID, ok := key.(common.TokenPairID)
		//if there is conversion error, continue to next key,val
		if !ok {
			err = fmt.Errorf("key (%v) cannot be asserted to TokenPairID", key)
			return false
		}
		exPrice, ok := value.(common.ExchangePrice)
		if !ok {
			err = fmt.Errorf("value (%v) cannot be asserted to ExchangePrice", value)
			return false
		}
		result[tokenPairID] = exPrice
		return true
	})
	return result, err
}

func (h *Huobi) OpenOrdersForOnePair(
	wg *sync.WaitGroup,
	pair common.TokenPair,
	data *sync.Map,
	timepoint uint64) {

	// defer wg.Done()

	// result, err := h.interf.OpenOrdersForOnePair(pair, timepoint)

	//TODO: complete open orders for one pair
}

func (h *Huobi) FetchOrderData(timepoint uint64) (common.OrderEntry, error) {
	result := common.OrderEntry{}
	result.Timestamp = common.Timestamp(fmt.Sprintf("%d", timepoint))
	result.Valid = true
	result.Data = []common.Order{}

	var (
		data = sync.Map{}
		err  error
	)

	result.ReturnTime = common.GetTimestamp()
	data.Range(func(key, value interface{}) bool {
		orders, ok := value.([]common.Order)
		if !ok {
			err = fmt.Errorf("cannot convert value (%v) to Order", value)
			return false
		}
		result.Data = append(result.Data, orders...)
		return true
	})
	return result, err
}

func (h *Huobi) FetchEBalanceData(timepoint uint64) (common.EBalanceEntry, error) {
	result := common.EBalanceEntry{}
	result.Timestamp = common.Timestamp(fmt.Sprintf("%d", timepoint))
	result.Valid = true
	result.Error = ""
	respData, err := h.interf.GetInfo()
	result.ReturnTime = common.GetTimestamp()
	if err != nil {
		result.Valid = false
		result.Error = err.Error()
		result.Status = false
	} else {
		result.AvailableBalance = map[string]float64{}
		result.LockedBalance = map[string]float64{}
		result.DepositBalance = map[string]float64{}
		result.Status = true
		if respData.Status != "ok" {
			result.Valid = false
			result.Error = fmt.Sprintf("Cannot fetch ebalance")
			result.Status = false
		} else {
			balances := respData.Data.List
			for _, b := range balances {
				tokenID := strings.ToUpper(b.Currency)
				_, err := h.sr.GetAssetBySymbol(uint64(common.Huobi), tokenID)
				if err == nil {
					balance, _ := strconv.ParseFloat(b.Balance, 64)
					if b.Type == "trade" {
						result.AvailableBalance[tokenID] = balance
					} else {
						result.LockedBalance[tokenID] = balance
					}
					result.DepositBalance[tokenID] = 0
				}
			}
		}
	}
	return result, nil
}

func (h *Huobi) FetchOnePairTradeHistory(
	wait *sync.WaitGroup,
	data *sync.Map,
	pair commonv3.TradingPairSymbols) {
	defer wait.Done()
	result := []common.TradeHistory{}
	resp, err := h.interf.GetAccountTradeHistory(pair.BaseSymbol, pair.QuoteSymbol)
	if err != nil {
		log.Printf("Cannot fetch data for pair %s%s: %s",
			pair.BaseSymbol, pair.QuoteSymbol, err.Error())
		return
	}
	pairString := pair.PairID()
	for _, trade := range resp.Data {
		price, _ := strconv.ParseFloat(trade.Price, 64)
		quantity, _ := strconv.ParseFloat(trade.Amount, 64)
		historyType := tradeTypeSell
		if trade.Type == "buy-limit" {
			historyType = tradeTypeBuy
		}
		tradeHistory := common.NewTradeHistory(
			strconv.FormatUint(trade.ID, 10),
			price,
			quantity,
			historyType,
			trade.Timestamp,
		)
		result = append(result, tradeHistory)
	}
	data.Store(pairString, result)
}

//FetchTradeHistory get all trade history for all pairs from huobi exchange
func (h *Huobi) FetchTradeHistory() {
	t := time.NewTicker(10 * time.Minute)
	go func() {
		for {
			result := map[common.TokenPairID][]common.TradeHistory{}
			data := sync.Map{}
			pairs, err := h.TokenPairs()
			if err != nil {
				log.Printf("Huobi fetch trade history failed (%s). This might due to pairs setting hasn't been init yet", err.Error())
				continue
			}
			wait := sync.WaitGroup{}
			for _, pair := range pairs {
				wait.Add(1)
				go h.FetchOnePairTradeHistory(&wait, &data, pair)
			}
			wait.Wait()
			var integrity = true
			data.Range(func(key, value interface{}) bool {
				tokenPairID, ok := key.(common.TokenPairID)
				//if there is conversion error, continue to next key,val
				if !ok {
					log.Printf("Key (%v) cannot be asserted to TokenPairID", key)
					integrity = false
					return false
				}
				tradeHistories, ok := value.([]common.TradeHistory)
				if !ok {
					log.Printf("Value (%v) cannot be asserted to []TradeHistory", value)
					integrity = false
					return false
				}
				result[tokenPairID] = tradeHistories
				return true
			})
			if !integrity {
				log.Print("Huobi fetch trade history returns corrupted. Try again in 10 mins")
				continue
			}
			if err := h.storage.StoreTradeHistory(result); err != nil {
				log.Printf("Store trade history error: %s", err.Error())
			}
			<-t.C
		}
	}()
}

func (h *Huobi) GetTradeHistory(fromTime, toTime uint64) (common.ExchangeTradeHistory, error) {
	return h.storage.GetTradeHistory(fromTime, toTime)
}

func (h *Huobi) Send2ndTransaction(amount float64, asset commonv3.Asset, exchangeAddress ethereum.Address) (*types.Transaction, error) {
	IAmount := common.FloatToBigInt(amount, int64(asset.Decimals))
	// Check balance, removed from huobi's blockchain object.
	// currBalance := h.blockchain.CheckBalance(token)
	// log.Printf("current balance of token %s is %d", token.ID, currBalance)
	// //h.blockchain.
	// if currBalance.Cmp(IAmount) < 0 {
	// 	log.Printf("balance is not enough, wait till next check")
	// 	return nil, errors.New("balance is not enough")
	// }
	var tx *types.Transaction
	var err error
	// TODO: add a check isETH that matching id instead of symbol
	if asset.Symbol == "ETH" {
		tx, err = h.blockchain.SendETHFromAccountToExchange(IAmount, exchangeAddress)
	} else {
		tx, err = h.blockchain.SendTokenFromAccountToExchange(IAmount, exchangeAddress, asset.Address)
	}
	if err != nil {
		log.Printf("ERROR: Can not send transaction to exchange: %v", err)
		return nil, err
	}
	log.Printf("Transaction submitted. Tx is: %v", tx)
	return tx, nil

}

func (h *Huobi) PendingIntermediateTxs() (map[common.ActivityID]common.TXEntry, error) {
	result, err := h.storage.GetPendingIntermediateTXs()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (h *Huobi) FindTx2InPending(id common.ActivityID) (common.TXEntry, bool) {
	pendings, err := h.storage.GetPendingIntermediateTXs()
	if err != nil {
		log.Printf("can't get pendings tx2 records: %v", err)
		return common.TXEntry{}, false
	}
	for actID, txentry := range pendings {
		if actID == id {
			return txentry, true
		}
	}
	return common.TXEntry{}, false
}

//FindTx2 : find Tx2 Record associates with activity ID, return
func (h *Huobi) FindTx2(id common.ActivityID) (tx2 common.TXEntry, found bool) {
	found = true
	//first look it up in permanent bucket
	tx2, err := h.storage.GetIntermedatorTx(id)
	if err != nil {
		//couldn't look for it in permanent bucket, look for it in pending bucket
		tx2, found = h.FindTx2InPending(id)
	}
	return tx2, found
}

func (h *Huobi) exchangeDepositStatus(id common.ActivityID, tx2Entry common.TXEntry, currency string, sentAmount float64) (string, error) {
	assets, err := h.sr.GetAssets()
	if err != nil {
		log.Printf("Huobi ERROR: Can not get list of assets from setting (%s)", err)
		return "", err
	}

	// make sure the size is enough for storing all deposit history
	deposits, err := h.interf.DepositHistory(len(assets) * 2)
	if err != nil || deposits.Status != "ok" {
		log.Printf("Huobi Getting deposit history from huobi failed, error: %v, status: %s", err, deposits.Status)
		return "", nil
	}
	//check tx2 deposit status from Huobi
	for _, deposit := range deposits.Data {
		log.Printf("deposit tx is %s, with token %s", deposit.TxHash, deposit.Currency)
		if deposit.TxHash[0:2] != "0x" {
			deposit.TxHash = "0x" + deposit.TxHash
		}
		if deposit.TxHash == tx2Entry.Hash {
			if deposit.State == "safe" || deposit.State == "confirmed" {
				data := common.NewTXEntry(tx2Entry.Hash,
					h.Name(),
					currency,
					"mined",
					exchangeStatusDone,
					sentAmount,
					common.GetTimestamp(),
				)
				if err = h.storage.StoreIntermediateTx(id, data); err != nil {
					log.Printf("Huobi Trying to store intermediate tx to huobi storage, error: %s. Ignore it and try later", err.Error())
					return "", nil
				}
				return exchangeStatusDone, nil
			}
			//TODO : handle other states following https://github.com/huobiapi/API_Docs_en/wiki/REST_Reference#deposit-states
			log.Printf("Huobi Tx %s is found but the status was not safe but %s", deposit.TxHash, deposit.State)
			return "", nil
		}
	}
	log.Printf("Huobi Deposit doesn't exist. Huobi hasn't recognized the deposit yet or in theory, you have more than %d deposits at the same time.", len(assets)*2)
	return "", nil
}

func (h *Huobi) process1stTx(id common.ActivityID, tx1Hash, currency string, sentAmount float64) (string, error) {
	status, blockno, err := h.blockchain.TxStatus(ethereum.HexToHash(tx1Hash))
	if err != nil {
		log.Printf("Huobi Can not get TX status (%s)", err.Error())
		return "", nil
	}
	log.Printf("Huobi Status for Tx1 was %s at block %d ", status, blockno)
	if status == common.MiningStatusMined {
		//if it is mined, send 2nd tx.
		log.Printf("Found a new deposit status, which deposit %f %s. Procceed to send it to Huobi", sentAmount, currency)
		//check if the asset is supported, the asset can be active or inactivee
		asset, err := h.sr.GetAssetBySymbol(uint64(common.Binance), currency)
		if err != nil {
			return "", err
		}
		exchangeAddress, err := h.RealDepositAddress(currency)
		if err != nil {
			return "", err
		}
		tx2, err := h.Send2ndTransaction(sentAmount, asset, exchangeAddress)
		if err != nil {
			log.Printf("Huobi Trying to send 2nd tx failed, error: %s. Will retry next time", err.Error())
			return "", nil
		}
		//store tx2 to pendingIntermediateTx
		data := common.NewTXEntry(
			tx2.Hash().Hex(),
			h.Name(),
			currency,
			common.MiningStatusSubmitted,
			"",
			sentAmount,
			common.GetTimestamp(),
		)
		if err = h.storage.StorePendingIntermediateTx(id, data); err != nil {
			log.Printf("Trying to store 2nd tx to pending tx storage failed, error: %s. It will be ignored and can make us to send to huobi again and the deposit will be marked as failed because the fund is not efficient", err.Error())
		}
		return "", nil
	}
	//No need to handle other blockchain status of TX1 here, since Fetcher will handle it from blockchain Status.
	return "", nil
}

func (h *Huobi) DepositStatus(id common.ActivityID, tx1Hash, currency string, sentAmount float64, timepoint uint64) (string, error) {
	var data common.TXEntry
	tx2Entry, found := h.FindTx2(id)
	//if not found, meaning there is no tx2 yet, process 1st Tx and send 2nd Tx.
	if !found {
		return h.process1stTx(id, tx1Hash, currency, sentAmount)
	}
	// if there is tx2Entry, check it blockchain status and handle the status accordingly:
	miningStatus, _, err := h.blockchain.TxStatus(ethereum.HexToHash(tx2Entry.Hash))
	if err != nil {
		return "", err
	}
	switch miningStatus {
	case common.MiningStatusMined:
		log.Println("Huobi 2nd Transaction is mined. Processed to store it and check the Huobi Deposit history")
		data = common.NewTXEntry(
			tx2Entry.Hash,
			h.Name(),
			currency,
			common.MiningStatusMined,
			"",
			sentAmount,
			common.GetTimestamp())
		if uErr := h.storage.StorePendingIntermediateTx(id, data); uErr != nil {
			log.Printf("Huobi Trying to store intermediate tx to huobi storage, error: %s. Ignore it and try later", uErr.Error())
			return "", nil
		}
		return h.exchangeDepositStatus(id, tx2Entry, currency, sentAmount)
	case common.MiningStatusFailed:
		data = common.NewTXEntry(
			tx2Entry.Hash,
			h.Name(),
			currency,
			common.MiningStatusFailed,
			common.ExchangeStatusFailed,
			sentAmount,
			common.GetTimestamp(),
		)
		if err = h.storage.StoreIntermediateTx(id, data); err != nil {
			log.Printf("Huobi Trying to store intermediate tx failed, error: %s. Ignore it and treat it like it is still pending", err.Error())
			return "", nil
		}
		return common.ExchangeStatusFailed, nil
	case common.MiningStatusLost:
		elapsed := common.GetTimepoint() - tx2Entry.Timestamp.MustToUint64()
		if elapsed > uint64(15*time.Minute/time.Millisecond) {
			data = common.NewTXEntry(
				tx2Entry.Hash,
				h.Name(),
				currency,
				common.MiningStatusLost,
				common.ExchangeStatusLost,
				sentAmount,
				common.GetTimestamp(),
			)
			if err = h.storage.StoreIntermediateTx(id, data); err != nil {
				log.Printf("Huobi Trying to store intermediate tx failed, error: %s. Ignore it and treat it like it is still pending", err.Error())
				return "", nil
			}
			log.Printf("Huobi The tx is not found for over 15mins, it is considered as lost and the deposit failed")
			return common.ExchangeStatusFailed, nil
		}
		return "", nil
	}
	return "", nil
}

//WithdrawStatus return withdraw status from huobi
func (h *Huobi) WithdrawStatus(
	id, currency string, amount float64, timepoint uint64) (string, string, error) {
	withdrawID, _ := strconv.ParseUint(id, 10, 64)
	assets, err := h.sr.GetAssets()
	if err != nil {
		return "", "", fmt.Errorf("huobi Can't get list of assets from setting (%s)", err)
	}
	// make sure the size is enough for storing all huobi withdrawal history
	withdraws, err := h.interf.WithdrawHistory(len(assets) * 2)
	if err != nil {
		return "", "", fmt.Errorf("can't get withdraw history from huobi: %s", err.Error())
	}
	log.Printf("Huobi Withdrawal id: %d", withdrawID)
	for _, withdraw := range withdraws.Data {
		if withdraw.ID == withdrawID {
			if withdraw.State == "confirmed" {
				if withdraw.TxHash[0:2] != "0x" {
					withdraw.TxHash = "0x" + withdraw.TxHash
				}
				return common.ExchangeStatusDone, withdraw.TxHash, nil
			}
			return "", withdraw.TxHash, nil
		}
	}
	return "", "", errors.New("huobi Withdrawal doesn't exist. This shouldn't happen unless tx returned from withdrawal from huobi and activity ID are not consistently designed")
}

func (h *Huobi) OrderStatus(id string, base, quote string) (string, error) {
	orderID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return "", err
	}
	symbol := base + quote
	order, err := h.interf.OrderStatus(symbol, orderID)
	if err != nil {
		return "", err
	}
	if order.Data.State == "pre-submitted" || order.Data.State == "submitting" || order.Data.State == "submitted" || order.Data.State == "partial-filled" || order.Data.State == "partial-canceled" {
		return "", nil
	}
	return common.ExchangeStatusDone, nil
}

//NewHuobi creates new Huobi exchange instance
func NewHuobi(
	interf HuobiInterface,
	blockchain *blockchain.BaseBlockchain,
	signer blockchain.Signer,
	nonce blockchain.NonceCorpus,
	storage HuobiStorage,
	sr storage.SettingReader,
) (*Huobi, error) {

	bc, err := huobiblockchain.NewBlockchain(blockchain, signer, nonce)
	if err != nil {
		return nil, err
	}

	huobiObj := Huobi{
		interf:     interf,
		blockchain: bc,
		storage:    storage,
		sr:         sr,
	}
	huobiObj.FetchTradeHistory()
	huobiServer := huobihttp.NewHuobiHTTPServer(&huobiObj)
	go huobiServer.Run()
	return &huobiObj, nil
}
