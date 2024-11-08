package services

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/smallbatch-apps/earnsmart-api/models"
)

type TransactionWatcher struct {
	watchMap           sync.Map
	TransactionService *TransactionService
}

func (tw *TransactionWatcher) WatchTransaction(tx models.Transaction) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour) // max wait time

	// Store cancellation function in case we need to stop watching
	tw.watchMap.Store(tx.ID, cancel)

	go func() {
		defer cancel()
		defer tw.watchMap.Delete(tx.ID)

		// Use exponential backoff for checking
		backoff := time.Second * 30
		ticker := time.NewTicker(backoff)
		defer ticker.Stop()
		startTime := time.Now()

		for {
			select {
			case <-ctx.Done():
				return
			case t := <-ticker.C:

				// automatically approve transactions after 10 minutes
				if t.Sub(startTime) > time.Minute*10 {
					err := tw.handleConfirmation(tx)
					if err != nil {
						tw.handleFailure(tx)
					}
					ticker.Stop()
					tw.watchMap.Delete(tx.ID)
					return
				}

				// Increase backoff (up to a max)
				if backoff < time.Hour {
					backoff *= 2
					ticker.Reset(backoff)
				}
			}
		}
	}()
}

func (tw *TransactionWatcher) handleConfirmation(tx models.Transaction) error {
	_, err := tw.TransactionService.ApproveTransaction(tx)
	if err != nil {
		log.Println("Error approving transaction:", err)
	}
	return err
}

func (tw *TransactionWatcher) handleFailure(tx models.Transaction) error {
	_, err := tw.TransactionService.DeclineTransaction(tx, models.TransactionStatusFailed)
	if err != nil {
		log.Println("Error failing transaction:", err)
	}
	return err
}

// Stop watching a specific transaction
func (tw *TransactionWatcher) StopWatching(txID string) {
	if cancel, ok := tw.watchMap.Load(txID); ok {
		cancel.(context.CancelFunc)()
	}
}
