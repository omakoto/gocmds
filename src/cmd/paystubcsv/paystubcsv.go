///bin/true; exec /usr/bin/env go run "$0" "$@"
package main

import (
	"fmt"
	"github.com/omakoto/go-common/src/common"
	"github.com/omakoto/go-common/src/textio/textblock"
	"github.com/omakoto/go-common/src/utils"
	"github.com/shopspring/decimal"
	"math"
	"os"
	"regexp"
	"strings"
)

var (
	reEarnings   = regexp.MustCompile(`^\s*Earnings\b`)
	reDeductions = regexp.MustCompile(`\bDeductions\s*$`)
	reTaxes      = regexp.MustCompile(`^\s*Taxes\s*$`)
	reSummary    = regexp.MustCompile(`^\s*Pay Summary\s*$`)
)

func main() {
	common.RunAndExit(realMain)
}

func parseDollars(val string) decimal.Decimal {
	orig := val
	positive := true

	if val[0] == '(' {
		common.OrPanicf(val[len(val)-1] == ')', "Invalid amount: %s", orig)
		val = val[1 : len(val)-1]
		positive = false
	}
	if val[0] == '$' {
		val = val[1:]
	}
	val = strings.Replace(val, ",", "", -1)

	ret, err := decimal.NewFromString(val)
	common.CheckPanicf(err, "Invalid amount: %s", orig)
	if positive {
		return ret
	}
	return ret.Neg()
}

func printCommaSeparated(values ...string) {
	b := strings.Builder{}
	for i, v := range values {
		if i > 0 {
			b.WriteByte(',')
		}
		if strings.ContainsRune(v, ',') {
			b.WriteByte('"')
			b.WriteString(strings.Replace(v, `"`, `""`, -1))
			b.WriteByte('"')
		} else {
			b.WriteString(v)
		}
	}
	fmt.Println(b.String())
}

func extractLabeledData(line string, label string) (bool, string) {
	labelPos := strings.Index(line, label)
	if labelPos < 0 {
		return false, ""
	}
	return true, strings.TrimSpace(line[labelPos+len(label):])
}

func extractField(line string, label string, skip int) (bool, string) {
	found, data := extractLabeledData(line, label)
	if !found {
		return false, ""
	}

	values := strings.Fields(data)

	if skip >= len(values) {
		return true, ""
	}
	return true, values[skip]
}

func extractDollars(line string, label string, skip int) (bool, decimal.Decimal) {
	found, data := extractLabeledData(line, label)
	if !found {
		return false, decimal.Zero
	}

	values := strings.Fields(data)

	if skip >= len(values) {
		return true, decimal.Zero
	}
	return true, parseDollars(values[skip])
}

func extractFirstField(lines []string, label string, skip int) string {
	for _, line := range lines {
		found, data := extractField(line, label, skip)
		if found {
			return data
		}
	}
	return ""
}

func extractDollarsSum(lines []string, label string, skip int, sum bool) decimal.Decimal {
	ret := decimal.Zero
	for _, line := range lines {
		found, dollars := extractDollars(line, label, skip)
		if found {
			ret = ret.Add(dollars)
			if !sum {
				break
			}
		}
	}
	return ret
}

func dumpPdf(file string) error {
	fmt.Fprintf(os.Stderr, "Scanning: %s\n", file)

	data, err := utils.ReadPdfAsText(file, true)
	common.Checkf(err, "Reading from %s failed.", file)

	b := textblock.NewFromBuffer(data)

	earningsX, earningsY := b.FindFirst(reEarnings)
	deductionsX, deductionsY := b.FindFirst(reDeductions)
	taxesX, taxesY := b.FindFirst(reTaxes)
	summaryX, summaryY := b.FindFirst(reSummary)

	//dry.Nop(taxesX, taxesY, earningsX, earningsY, deductionsX, deductionsY)

	for y, line := range b.Lines() {
		common.Debugf("%3d: %s\n", y, string(line))
	}

	header := b.Copy(0, 0, math.MaxInt32, earningsY)
	earnings := b.Copy(0, earningsY+1, deductionsX, taxesY)
	deductions := b.Copy(deductionsX, deductionsY+1, math.MaxInt32, taxesY)
	taxes := b.Slice(taxesY+1, summaryY)
	summary := b.Slice(summaryY+1, b.Size())

	common.Debugf("- Header")
	common.Debugf("%s\n", header.Bytes())

	common.Debugf("- Earnings: [%d, %d]\n", earningsX, earningsY)
	common.Debugf("%s\n", earnings.Bytes())

	common.Debugf("- Deductions: [%d, %d]\n", deductionsX, deductionsY)
	common.Debugf("%s\n", deductions.Bytes())

	common.Debugf("- Taxes: [%d, %d]\n", taxesX, taxesY)
	common.Debugf("%s\n", taxes.Bytes())

	common.Debugf("- Pay Summary: [%d, %d]\n", summaryX, summaryY)
	common.Debugf("%s\n", summary.Bytes())

	//dataStr := string(data)
	//
	//lines := strings.Split(dataStr, "\n")
	//
	//printCommaSeparated(
	//	filepath.Base(file),
	//
	//	extractFirstField(lines, "Period Start Date", 0),
	//	extractFirstField(lines, "Period End Date", 0),
	//	extractFirstField(lines, "Pay Date", 0),
	//	extractFirstField(lines, "Document", 0),
	//	extractDollarsSum(lines, "Net Pay", 0, false).String(),
	//
	//	extractDollarsSum(lines, "Regular Pay", 2, true).String(),
	//	extractDollarsSum(lines, "Annual Bonus", 2, true).String(),
	//	extractDollarsSum(lines, "Group Term Life", 2, true).String(),
	//	extractDollarsSum(lines, "Goog Stock", 5, true).String(), // TODO -- Depends on job title
	//	extractDollarsSum(lines, "Peer Bonus", 2, true).String(),
	//
	//)

	return nil
}

func printHeader() {
	printCommaSeparated(
		"File",

		"Start Date",
		"End Date",
		"Pay Date",
		"Document",
		"Net Pay",

		"Regular Pay",
		"Annual Bonus",
		"Group Term Life",
		"GSU",
		"Peer Bonus",

		"",
		"",
	)
}

func realMain() int {
	printHeader()
	for _, file := range os.Args[1:] {
		dumpPdf(file)
	}

	return 0
}
