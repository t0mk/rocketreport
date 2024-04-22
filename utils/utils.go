package utils

import (
	"fmt"

	"io"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/params"
	"github.com/rocket-pool/smartnode/shared/types/api"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func GetHTTPResponseBodyFromUrl(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("http.Get: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadAll: %v", err)
	}
	return body, nil
}

func SingleClientStatusString(status api.ClientStatus) string {
	sentence := " "
	if !status.IsWorking {
		sentence += "not working,"
	}
	if status.IsSynced {
		sentence += "synced"
	} else {
		sentence += "not synced,"
	}
	if status.SyncProgress < 1 {
		sentence += fmt.Sprintf(" syncing, now at %d%%", int(100*(status.SyncProgress)))
	}
	if status.Error != "" {
		sentence += fmt.Sprintf(", Error: %s", status.Error)
	}
	return sentence
}

func EthClientStatusString(status *api.ClientManagerStatus) string {
	sentence := "Prim" + SingleClientStatusString(status.PrimaryClientStatus)
	if status.FallbackEnabled {
		sentence += ", FB " + SingleClientStatusString(status.FallbackClientStatus)
	} else {
		sentence += ", FB n/a"
	}
	return sentence
}

func WeiToEther(wei *big.Int) *big.Float {
	return new(big.Float).Quo(new(big.Float).SetInt(wei), big.NewFloat(params.Ether))
}

func FmtEth(p float64) string {
	return fmt.Sprintf("%.6f", p)
}

func FmtRplFiat(p float64) string {
	return fmt.Sprintf("%.2f", p)
}

func FmtRpl(p float64) string {
	return fmt.Sprintf("%.1f", p)
}

func FmtFiat(p float64) string {
	f := message.NewPrinter(language.English)
	i := int(p)
	return f.Sprintf("%d", i)
}

func ToIfSlice[T any](slice []T) []interface{} {
	convertedSlice := make([]interface{}, len(slice))
	for i, v := range slice {
		convertedSlice[i] = v
	}
	return convertedSlice
}
