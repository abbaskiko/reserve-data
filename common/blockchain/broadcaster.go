package blockchain

import (
	"context"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Broadcaster takes a signed tx and try to broadcast it to all
// nodes that it manages as fast as possible. It returns a map of
// failures and a bool indicating that the tx is broadcasted to
// at least 1 node
type Broadcaster struct {
	clients map[string]*ethclient.Client
}

func (b Broadcaster) sendTx(client ethereum.TransactionSender, tx *types.Transaction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	err := client.SendTransaction(ctx, tx)
	cancel()
	return err
}

func (b Broadcaster) Broadcast(tx *types.Transaction) (map[string]error, bool) {
	failures := sync.Map{}
	wg := sync.WaitGroup{}
	for id, client := range b.clients {
		wg.Add(1)
		go func(cid string, c *ethclient.Client) {
			defer wg.Done()
			if err := b.sendTx(c, tx); err != nil {
				failures.Store(cid, err)
			}
		}(id, client)
	}
	wg.Wait()
	errorDetail := map[string]error{}
	failures.Range(func(key, value interface{}) bool {
		errorDetail[key.(string)] = value.(error) // as failures access scope is this func, we sure it's a string->error map
		return true
	})
	return errorDetail, len(errorDetail) != len(b.clients) && len(b.clients) > 0
}

func NewBroadcaster(clients map[string]*ethclient.Client) *Broadcaster {
	return &Broadcaster{
		clients: clients,
	}
}
