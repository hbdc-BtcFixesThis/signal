package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const mempoolUTXO = "https://mempool.space/api/address/%s/utxo"
const mempoolBlocks = "https://mempool.space/api/v1/blocks"

type mempoolResult struct {
	Value uint64 `json:"value"`
}

func BtcAddressTotal(addr string) (uint64, error) {
	return uint64(100), nil
	addrURL := fmt.Sprintf(mempoolUTXO, addr)
	resp, err := http.Get(addrURL)
	if err != nil {
		return uint64(0), err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return uint64(0), err
	}

	if resp.StatusCode != http.StatusOK {
		return uint64(0), errors.New(string(body))
	}

	var results []mempoolResult
	if err := json.Unmarshal(body, &results); err != nil {
		return uint64(0), errors.New(string(body))
	}

	var total uint64
	for i := 0; i < len(results); i++ {
		total += results[i].Value
	}

	return total, nil
}

// I realize I could save space on the return but lets
// bee optimistic keep it a 64. Can you imagine?
func LatestBtcBlockHeight() (uint64, error) {
	type block struct {
		Height uint64 `json:"height"`
	}
	resp, err := http.Get(mempoolBlocks)
	if err != nil {
		return uint64(0), err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return uint64(0), err
	}

	if resp.StatusCode != http.StatusOK {
		return 0, errors.New(string(body))
	}

	var blocks []block
	if err := json.Unmarshal(body, &blocks); err != nil {
		return uint64(0), errors.New(string(body))
	}
	return blocks[0].Height, nil
}

/*
func main() {
	a, b := LatestBtcBlockHeight()
	fmt.Println(a, b)
}

/*func main() {
	goodAddress := "1KFHE7w8BhaENAswwryaoccDb6qcT6DbYY"
	goodResp, err1 := BtcAddressTotal(goodAddress)
	if err1 != nil {
		fmt.Println(err1)
	}
	fmt.Println(goodResp)

	badAddress := "1KFHE7w8BhaENAswwryaoccDb6qcT6DbYY_"
	badResp, err2 := BtcAddressTotal(badAddress)
	if err2 != nil {
		fmt.Println(err2, "Invalid Bitcoin address" == err2.Error())
	}
	fmt.Println(badResp)
}*/
