package utils

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"io"
	"math/big"
	"net/http"

	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
	"github.com/rocket-pool/smartnode/shared/types/api"
	"github.com/t0mk/rocketreport/config"
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

func IfSliceToString(slice []interface{}) string {
	convertedSlice := make([]string, len(slice))
	for i, v := range slice {
		strval, err := toString(v)
		if err != nil {
			convertedSlice[i] = fmt.Sprintf("Error converting value: %v", err)
			continue
		}
		convertedSlice[i] = strval
	}
	return strings.Join(convertedSlice, " ")
}

func toString(value interface{}) (string, error) {
	switch v := value.(type) {
	case string:
		return v, nil
	case float64:
		// Convert float64 to string using strconv.FormatFloat.
		// The 'f' indicates a decimal without exponent,
		// -1 specifies the smallest number of digits necessary,
		// and 64 means it's a float64.
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	case int64:
		// Convert int64 to string using strconv.FormatInt.
		return strconv.FormatInt(v, 10), nil
	default:
		return "", fmt.Errorf("unsupported type %T", value)
	}
}

func ToIfSlice[T any](slice []T) []interface{} {
	convertedSlice := make([]interface{}, len(slice))
	for i, v := range slice {
		convertedSlice[i] = v
	}
	return convertedSlice
}

func AddressBalanceString(address string) (float64, error) {
	return AddressBalance(common.HexToAddress(address))
}

func AddressBalance(address common.Address) (float64, error) {
	balanceRaw, err := config.RP().Client.BalanceAt(context.Background(), address, nil)
	if err != nil {
		return 0, fmt.Errorf("error getting balance: %w", err)
	}
	balance, _ := WeiToEther(balanceRaw).Float64()
	return balance, nil
}

func ValidateAndParseAddress(address string) (*common.Address, bool) {
	if !common.IsHexAddress(address) {
		return nil, false
	}
	addr := common.HexToAddress(address)
	return &addr, true
}

type AddressBalanceEtherscanResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  string `json:"result"`
}

func AddressBalanceEtherscan(address common.Address) (float64, error) {
	url := fmt.Sprintf("https://api.etherscan.io/api?module=account&action=balance&address=%s&tag=latest", address.Hex())
	body, err := GetHTTPResponseBodyFromUrl(url)
	if err != nil {
		return 0, fmt.Errorf("error getting Etherscan balance: %w", err)
	}
	var response AddressBalanceEtherscanResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return 0, fmt.Errorf("error unmarshalling Etherscan response: %w", err)
	}
	if response.Status != "1" {
		return 0, fmt.Errorf("error Etherscan response status: %s", response.Message)
	}
	balanceRaw, ok := new(big.Int).SetString(response.Result, 10)
	if !ok {
		return 0, fmt.Errorf("error parsing Etherscan balance: %s", response.Result)
	}
	balance, _ := WeiToEther(balanceRaw).Float64()

	return balance, nil
}
