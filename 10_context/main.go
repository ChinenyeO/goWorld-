package main

import (
	"context"
	"fmt"
	"time"
)

// ProcessOrder charges the customer, then reserves inventory.
func ProcessOrder(ctx context.Context, orderID string) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
    	defer wg.Done()
    	chargeCustomer(ctx, orderID)
	}()

	go func() {
    	defer wg.Done()
    	reserveInventory(ctx, orderID)
	}()

	wg.Wait()

	return nil
}

// chargeCustomer pretends to talk to a payment service.
func chargeCustomer(ctx context.Context, orderID string) error {
    select {
    case <-time.After(5 * time.Second):
        fmt.Println("chargeCustomer", orderID)
		return nil

    case <-ctx.Done():
        return ctx.Err()
    }
}

// reserveInventory pretends to talk to a warehouse service.
func reserveInventory(ctx context.Context,orderID string) error {
	fmt.Println("reserveInventory: started for", orderID)
	time.Sleep(2 * time.Second) // simulate slow work
	fmt.Println("reserveInventory: finished for", orderID)
	return nil
}

func main() {
	ctx := context.Background()

	start := time.Now()
	err := ProcessOrder(ctx, "order-123")
	if err != nil {
		fmt.Println("order failed:", err)
	} else {
		fmt.Println("order succeeded")
	}
	fmt.Println("took", time.Since(start).Round(time.Millisecond))
}
