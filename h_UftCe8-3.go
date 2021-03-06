package main

// How many Siacoins will there be at block X?
// What is the expected date of block X?

import "fmt"
import "time"

func main() {
	// The times chosen are not perfect. For example,
	// the actual genesis block time was not used, and
	// the target block time of 10 mintes was not used.
	// This is because when the hashrate increases, a
	// series of faster blocks are found before the
	// difficulty adjusts - and so the actual average
	// block time is less than the target time.
	//
	// The equations used result in the estimated time
	// for block 16,000 being Sept. 1st, which is what
	// was observed on the real network.

	HEIGHT := 16 * 1000
	
	totalCoins := uint64(0)
	blockReward := uint64(300 * 1000)
	expectedDate := time.Date(2015, 05, 24, 0, 0, 0, 0, time.UTC)
	for i := 0; i < HEIGHT; i++ {
		totalCoins += blockReward
		if blockReward > 30*1000 {
			blockReward--
		}
		expectedDate = expectedDate.Add(time.Duration(time.Minute * 9))
	}

	fmt.Printf("Calculating Stats for Height: %v\n\n", HEIGHT)
	fmt.Printf("Expected Number of Coins: %.3f Billion\n", float64(totalCoins)/1e9)
	fmt.Printf("Expected Date:            %v", expectedDate)
}
