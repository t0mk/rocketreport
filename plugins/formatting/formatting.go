package formatting

import (
	"fmt"
	"math"
	"time"

	"github.com/t0mk/rocketreport/prices"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

const (
	ColorReset string = "\033[0m"
	ColorRed   string = "\033[31m"
	ColorGreen string = "\033[32m"
	ColorBlue  string = "\033[34m"
	ColorBlack string = "\033[1;30m"
	ColorBold  string = ""
)

func Str(i interface{}) string {
	return i.(string)
}

func FloatSuffix(ndecs int, suffix string) func(interface{}) string {
	return func(i interface{}) string {
		f := message.NewPrinter(language.English)
		replacedSuffix := prices.FindAndReplaceAllCurrencyOccurencesBySign(suffix)
		return f.Sprintf("%.*f %s", ndecs, i.(float64), replacedSuffix)
	}
}

func SmartFloatSuffix(suffix string) func(interface{}) string {
	return func(i interface{}) string {
		n := SmartFloat(i)
		replacedSuffix := prices.FindAndReplaceAllCurrencyOccurencesBySign(suffix)
		return fmt.Sprintf("%s %s", n, replacedSuffix)
	}
}

func Time (i interface{}) string {
	return i.(time.Time).Format("2006-01-02 15:04:05")
}


func SmartFloat(i interface{}) string {
	f := i.(float64)
	pr := message.NewPrinter(language.English)
	absVal := math.Abs(f)
	if absVal == 0.	{	
		return pr.Sprintf("%.0f", f)
	}
	if absVal < 1 {
		return pr.Sprintf("%.6f", f)
	}
	if absVal < 2.5 {
		return pr.Sprintf("%.4f", f)
	}
	if absVal < 100 {
		return pr.Sprintf("%.2f", f)
	}
	return pr.Sprintf("%.0f", f)
}

func Uint(i interface{}) string {
	return fmt.Sprintf("%d", i.(uint64))
}

func Duration(i interface{}) string {
	duration := i.(time.Duration)
	days := duration / (time.Hour * 24)
	hours := (duration % (time.Hour * 24)) / time.Hour
	minutes := (duration % time.Hour) / time.Minute

	formatted := ""
	if days > 0 {
		formatted += fmt.Sprintf("%dd ", days)
	}
	if hours > 0 {
		formatted += fmt.Sprintf("%dh ", hours)
	}
	if minutes > 0 {
		formatted += fmt.Sprintf("%dmin", minutes)
	}

	return formatted
}
