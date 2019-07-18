package exchange

import (
	"fmt"
	"log"
	"math/big"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"

	"github.com/KyberNetwork/reserve-data/common"
	commonv3 "github.com/KyberNetwork/reserve-data/v3/common"
	"github.com/KyberNetwork/reserve-data/v3/storage"
)

const (
	binanceEpsilon float64 = 0.0000001 // 10e-7
	batchSize      int     = 4
)

type Binance struct {
	interf  BinanceInterface
	storage BinanceStorage
	sr      storage.SettingReader
}

func (bn *Binance) TokenAddresses() (map[string]ethereum.Address, error) {
	result, err := bn.sr.GetDepositAddresses(uint64(common.Binance))
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (bn *Binance) MarshalText() (text []byte, err error) {
	return []byte(bn.ID()), nil
}

// Address returns the deposit address of given token.
func (bn *Binance) Address(token common.Token) (ethereum.Address, bool) {
	liveAddress, err := bn.interf.GetDepositAddress(token.ID)
	if err != nil || liveAddress.Address == "" {
		log.Printf("WARNING: Get Binance live deposit address for token %s failed: err: (%v) or the address repplied is empty . Use the currently available address instead", token.ID, err)
		addrs, uErr := bn.sr.GetDepositAddresses(uint64(common.Binance))
		if uErr != nil {
			log.Printf("WARNING: get address of token %s in Binance exchange failed:(%s), it will be considered as not supported", token.ID, err.Error())
			return ethereum.Address{}, false
		}
		depositAddr, ok := addrs[token.ID]
		return depositAddr, ok
	}
	log.Printf("Got Binance live deposit address for token %s, attempt to update it to current setting", token.ID)
	addrs := common.NewExchangeAddresses()
	addrs.Update(token.ID, ethereum.HexToAddress(liveAddress.Address))

	// TODO: update deposit address
	//if err = bn.setting.UpdateDepositAddress(common.Binance, *addrs, common.GetTimepoint()); err != nil {
	//	log.Printf("WARNING: cannot update deposit address for token %s on Binance: (%s)", token.ID, err.Error())
	//}
	return ethereum.HexToAddress(liveAddress.Address), true
}

func (bn *Binance) precisionFromStepSize(stepSize string) int {
	re := regexp.MustCompile("0*$")
	parts := strings.Split(re.ReplaceAllString(stepSize, ""), ".")
	if len(parts) > 1 {
		return len(parts[1])
	}
	return 0
}

// GetLiveExchangeInfo queries the Exchange Endpoint for exchange precision and limit of a certain pair ID
// It return error if occurs.
func (bn *Binance) GetLiveExchangeInfos(pairs []commonv3.TradingPairSymbols) (common.ExchangeInfo, error) {
	result := make(common.ExchangeInfo)
	exchangeInfo, err := bn.interf.GetExchangeInfo()
	if err != nil {
		return result, err
	}
	symbols := exchangeInfo.Symbols
	for _, pair := range pairs {
		exchangePrecisionLimit, ok := bn.getPrecisionLimitFromSymbols(pair, symbols)
		if !ok {
			return result, fmt.Errorf("binance Exchange Info reply doesn't contain token pair %d", pair.ID)
		}
		result[pair.ID] = exchangePrecisionLimit
	}
	return result, nil
}

// getPrecisionLimitFromSymbols find the pairID amongs symbols from exchanges,
// return ExchangePrecisionLimit of that pair and true if the pairID exist amongs symbols, false if otherwise
func (bn *Binance) getPrecisionLimitFromSymbols(pair commonv3.TradingPairSymbols, symbols []BinanceSymbol) (common.ExchangePrecisionLimit, bool) {
	var result common.ExchangePrecisionLimit
	pairName := strings.ToUpper(fmt.Sprintf("%s%s", pair.BaseSymbol, pair.QuoteSymbol))
	for _, symbol := range symbols {
		if strings.ToUpper(symbol.Symbol) == pairName {
			//update precision
			result.Precision.Amount = symbol.BaseAssetPrecision
			result.Precision.Price = symbol.QuotePrecision
			// update limit
			for _, filter := range symbol.Filters {
				if filter.FilterType == "LOT_SIZE" {
					// update amount min
					minQuantity, _ := strconv.ParseFloat(filter.MinQuantity, 64)
					result.AmountLimit.Min = minQuantity
					// update amount max
					maxQuantity, _ := strconv.ParseFloat(filter.MaxQuantity, 64)
					result.AmountLimit.Max = maxQuantity
					result.Precision.Amount = bn.precisionFromStepSize(filter.StepSize)
				}

				if filter.FilterType == "PRICE_FILTER" {
					// update price min
					minPrice, _ := strconv.ParseFloat(filter.MinPrice, 64)
					result.PriceLimit.Min = minPrice
					// update price max
					maxPrice, _ := strconv.ParseFloat(filter.MaxPrice, 64)
					result.PriceLimit.Max = maxPrice
					result.Precision.Price = bn.precisionFromStepSize(filter.TickSize)
				}

				if filter.FilterType == "MIN_NOTIONAL" {
					minNotional, _ := strconv.ParseFloat(filter.MinNotional, 64)
					result.MinNotional = minNotional
				}
			}
			return result, true
		}
	}
	return result, false
}

// ID must return the exact string or else simulation will fail
func (bn *Binance) ID() common.ExchangeID {
	return common.ExchangeID(common.Binance.String())
}

func (bn *Binance) TokenPairs() ([]commonv3.TradingPairSymbols, error) {
	pairs, err := bn.sr.GetTradingPairs(uint64(common.Binance))
	if err != nil {
		return nil, err
	}
	return pairs, nil
}

func (bn *Binance) Name() string {
	return "binance"
}

func (bn *Binance) QueryOrder(symbol string, id uint64) (done float64, remaining float64, finished bool, err error) {
	result, err := bn.interf.OrderStatus(symbol, id)
	if err != nil {
		return 0, 0, false, err
	}
	done, _ = strconv.ParseFloat(result.ExecutedQty, 64)
	total, _ := strconv.ParseFloat(result.OrigQty, 64)
	return done, total - done, total-done < binanceEpsilon, nil
}

func (bn *Binance) Trade(tradeType string, base common.Token, quote common.Token, rate float64, amount float64, timepoint uint64) (id string, done float64, remaining float64, finished bool, err error) {
	result, err := bn.interf.Trade(tradeType, base, quote, rate, amount)
	if err != nil {
		return "", 0, 0, false, err
	}
	done, remaining, finished, err = bn.QueryOrder(
		base.ID+quote.ID,
		result.OrderID,
	)
	id = strconv.FormatUint(result.OrderID, 10)
	return id, done, remaining, finished, err
}

func (bn *Binance) Withdraw(token common.Token, amount *big.Int, address ethereum.Address, timepoint uint64) (string, error) {
	tx, err := bn.interf.Withdraw(token, amount, address)
	return tx, err
}

func (bn *Binance) CancelOrder(id string, base, quote string) error {
	idNo, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return err
	}
	symbol := base + quote
	_, err = bn.interf.CancelOrder(symbol, idNo)
	if err != nil {
		return err
	}
	return nil
}

func (bn *Binance) FetchOnePairData(
	wg *sync.WaitGroup,
	pair commonv3.TradingPairSymbols,
	data *sync.Map,
	timepoint uint64) {

	defer wg.Done()
	result := common.ExchangePrice{}

	timestamp := common.Timestamp(fmt.Sprintf("%d", timepoint))
	result.Timestamp = timestamp
	result.Valid = true
	respData, err := bn.interf.GetDepthOnePair(pair.BaseSymbol, pair.QuoteSymbol)
	returnTime := common.GetTimestamp()
	result.ReturnTime = returnTime
	if err != nil {
		result.Valid = false
		result.Error = err.Error()
	} else {
		if respData.Code != 0 || respData.Msg != "" {
			result.Valid = false
			result.Error = fmt.Sprintf("Code: %d, Msg: %s", respData.Code, respData.Msg)
		} else {
			for _, buy := range respData.Bids {
				quantity, _ := strconv.ParseFloat(buy.Quantity, 64)
				rate, _ := strconv.ParseFloat(buy.Rate, 64)
				result.Bids = append(
					result.Bids,
					common.NewPriceEntry(
						quantity,
						rate,
					),
				)
			}
			for _, sell := range respData.Asks {
				quantity, _ := strconv.ParseFloat(sell.Quantity, 64)
				rate, _ := strconv.ParseFloat(sell.Rate, 64)
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
	data.Store(pair.ID, result)
}

func (bn *Binance) FetchPriceData(timepoint uint64) (map[uint64]common.ExchangePrice, error) {
	wait := sync.WaitGroup{}
	data := sync.Map{}
	pairs, err := bn.TokenPairs()
	if err != nil {
		return nil, err
	}
	var (
		i int
		x int
	)
	for i < len(pairs) {
		for x = i; x < len(pairs) && x < i+batchSize; x++ {
			wait.Add(1)
			pair := pairs[x]
			go bn.FetchOnePairData(&wait, pair, &data, timepoint)
		}
		wait.Wait()
		i = x
	}
	result := map[uint64]common.ExchangePrice{}
	data.Range(func(key, value interface{}) bool {
		//if there is conversion error, continue to next key,val
		tokenPairID, ok := key.(uint64)
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

func (bn *Binance) FetchEBalanceData(timepoint uint64) (common.EBalanceEntry, error) {
	result := common.EBalanceEntry{}
	result.Timestamp = common.Timestamp(fmt.Sprintf("%d", timepoint))
	result.Valid = true
	result.Error = ""
	respData, err := bn.interf.GetInfo()
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
		if respData.Code != 0 {
			result.Valid = false
			result.Error = fmt.Sprintf("Code: %d, Msg: %s", respData.Code, respData.Msg)
			result.Status = false
		} else {
			for _, b := range respData.Balances {
				tokenID := b.Asset
				_, err := bn.sr.GetAssetBySymbol(uint64(common.Binance), tokenID)
				if err == nil {
					avai, _ := strconv.ParseFloat(b.Free, 64)
					locked, _ := strconv.ParseFloat(b.Locked, 64)
					result.AvailableBalance[tokenID] = avai
					result.LockedBalance[tokenID] = locked
					result.DepositBalance[tokenID] = 0
				}
			}
		}
	}
	return result, nil
}

//FetchOnePairTradeHistory fetch trade history for one pair from exchange
func (bn *Binance) FetchOnePairTradeHistory(
	wait *sync.WaitGroup,
	data *sync.Map,
	pair commonv3.TradingPairSymbols) {

	defer wait.Done()
	result := []common.TradeHistory{}
	fromID, err := bn.storage.GetLastIDTradeHistory(pair.ID)
	if err != nil {
		log.Printf("Cannot get last ID trade history: %s", err.Error())
	}
	resp, err := bn.interf.GetAccountTradeHistory(pair.BaseSymbol, pair.QuoteSymbol, fromID)
	if err != nil {
		log.Printf("Binance Cannot fetch data for pair %s%s: %s", pair.BaseSymbol, pair.QuoteSymbol, err.Error())
	}
	for _, trade := range resp {
		price, _ := strconv.ParseFloat(trade.Price, 64)
		quantity, _ := strconv.ParseFloat(trade.Qty, 64)
		historyType := "sell"
		if trade.IsBuyer {
			historyType = "buy"
		}
		tradeHistory := common.NewTradeHistory(
			strconv.FormatUint(trade.ID, 10),
			price,
			quantity,
			historyType,
			trade.Time,
		)
		result = append(result, tradeHistory)
	}
	data.Store(pair.ID, result)
}

//FetchTradeHistory get all trade history for all tokens in the exchange
func (bn *Binance) FetchTradeHistory() {
	t := time.NewTicker(10 * time.Minute)
	go func() {
		for {
			result := common.ExchangeTradeHistory{}
			data := sync.Map{}
			pairs, err := bn.TokenPairs()
			if err != nil {
				log.Printf("Binance Get Token pairs setting failed (%s)", err.Error())
				continue
			}
			wait := sync.WaitGroup{}
			var i int
			var x int
			for i < len(pairs) {
				for x = i; x < len(pairs) && x < i+batchSize; x++ {
					wait.Add(1)
					pair := pairs[x]
					go bn.FetchOnePairTradeHistory(&wait, &data, pair)
				}
				i = x
				wait.Wait()
			}
			var integrity = true
			data.Range(func(key, value interface{}) bool {
				tokenPairID, ok := key.(uint64)
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
				log.Print("Binance fetch trade history returns corrupted. Try again in 10 mins")
				continue
			}
			if err := bn.storage.StoreTradeHistory(result); err != nil {
				log.Printf("Binance Store trade history error: %s", err.Error())
			}
			<-t.C
		}
	}()
}

func (bn *Binance) GetTradeHistory(fromTime, toTime uint64) (common.ExchangeTradeHistory, error) {
	return bn.storage.GetTradeHistory(fromTime, toTime)
}

func (bn *Binance) DepositStatus(id common.ActivityID, txHash, currency string, amount float64, timepoint uint64) (string, error) {
	startTime := timepoint - 86400000
	endTime := timepoint
	deposits, err := bn.interf.DepositHistory(startTime, endTime)
	if err != nil || !deposits.Success {
		return "", err
	}
	for _, deposit := range deposits.Deposits {
		if deposit.TxID == txHash {
			if deposit.Status == 1 {
				return common.ExchangeStatusDone, nil
			}
			return "", nil
		}
	}
	log.Printf("Binance Deposit is not found in deposit list returned from Binance. This might cause by wrong start/end time, please check again.")
	return "", nil
}

func (bn *Binance) WithdrawStatus(id, currency string, amount float64, timepoint uint64) (string, string, error) {
	startTime := timepoint - 86400000
	endTime := timepoint
	withdraws, err := bn.interf.WithdrawHistory(startTime, endTime)
	if err != nil || !withdraws.Success {
		return "", "", err
	}
	for _, withdraw := range withdraws.Withdrawals {
		if withdraw.ID == id {
			if withdraw.Status == 3 || withdraw.Status == 5 || withdraw.Status == 6 {
				return common.ExchangeStatusDone, withdraw.TxID, nil
			}
			return "", withdraw.TxID, nil
		}
	}
	log.Printf("Binance Withdrawal doesn't exist. This shouldn't happen unless tx returned from withdrawal from binance and activity ID are not consistently designed")
	return "", "", nil
}

func (bn *Binance) OrderStatus(id string, base, quote string) (string, error) {
	orderID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return "", fmt.Errorf("can not parse orderID (val %s) to uint", id)
	}
	symbol := base + quote
	order, err := bn.interf.OrderStatus(symbol, orderID)
	if err != nil {
		return "", err
	}
	if order.Status == "NEW" || order.Status == "PARTIALLY_FILLED" || order.Status == "PENDING_CANCEL" {
		return "", nil
	}
	return common.ExchangeStatusDone, nil
}

func NewBinance(
	interf BinanceInterface,
	storage BinanceStorage,
	sr storage.SettingReader) (*Binance, error) {
	binance := &Binance{
		interf:  interf,
		storage: storage,
		sr:      sr,
	}
	binance.FetchTradeHistory()
	return binance, nil
}
