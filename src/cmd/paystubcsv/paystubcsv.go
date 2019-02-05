///bin/true; exec /usr/bin/env go run "$0" "$@"

// Parse paystub PDF into CSV
package main

import (
	"bytes"
	"fmt"
	"github.com/omakoto/go-common/src/common"
	"github.com/omakoto/go-common/src/textio/textblock"
	"github.com/omakoto/go-common/src/utils"
	"github.com/shopspring/decimal"
	"math"
	"os"
	"path/filepath"
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

	val = strings.TrimRight(val, "\x00")

	if val[0] == '(' {
		common.OrPanicf(val[len(val)-1] == ')', "Invalid amount: %#v", orig)
		val = val[1 : len(val)-1]
		positive = false
	}
	if val[0] == '$' {
		val = val[1:]
	}
	val = strings.Replace(val, ",", "", -1)

	ret, err := decimal.NewFromString(val)
	common.CheckPanicf(err, "Invalid amount: %#v", orig)
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

func extractLabeledData(line string, labelPat string) (bool, string) {
	p := regexp.MustCompile(labelPat)
	loc := p.FindStringIndex(line)

	if loc == nil {
		return false, ""
	}
	return true, strings.TrimSpace(line[loc[1]:])
}

func extractField(line string, labelPat string, skip int) (bool, string) {
	found, data := extractLabeledData(line, labelPat)
	if !found {
		return false, ""
	}

	values := strings.Fields(data)

	if skip >= len(values) {
		return true, ""
	}
	return true, values[skip]
}

func extractDollars(line string, labelPat string, skip int) (bool, decimal.Decimal) {
	found, data := extractLabeledData(line, labelPat)
	if !found {
		return false, decimal.Zero
	}

	values := strings.Fields(data)

	if skip >= len(values) {
		return true, decimal.Zero
	}
	return true, parseDollars(values[skip])
}

func extractFirstField(lines []string, labelPat string, skip int) string {
	for _, line := range lines {
		found, data := extractField(line, labelPat, skip)
		if found {
			return data
		}
	}
	return ""
}

func extractDollarsSum(lines []string, labelPat string, skip int, sum bool) decimal.Decimal {
	ret := decimal.Zero
	for _, line := range lines {
		found, dollars := extractDollars(line, labelPat, skip)
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

	data = bytes.Replace(data, []byte("\x00"), []byte(""), -1)

	// Extract the blocks.
	b := textblock.NewFromBuffer(data)

	earningsX, earningsY := b.FindFirst(reEarnings)
	deductionsX, deductionsY := b.FindFirst(reDeductions)
	taxesX, taxesY := b.FindFirst(reTaxes)
	summaryX, summaryY := b.FindFirst(reSummary)

	for y, line := range b.Lines() {
		common.Debugf("%3d: %s\n", y, string(line))
	}

	header := b.Copy(0, 0, math.MaxInt32, earningsY).LineStrings()
	earnings := b.Copy(0, earningsY+1, deductionsX, taxesY).LineStrings()
	deductions := b.Copy(deductionsX, deductionsY+1, math.MaxInt32, taxesY).LineStrings()
	taxes := b.Slice(taxesY+1, summaryY).LineStrings()
	summary := b.Slice(summaryY+1, b.Size()).LineStrings()

	common.Debugf("- Header")
	common.Debugf("%s\n", strings.Join(header, "\n"))

	common.Debugf("- Earnings: [%d, %d]\n", earningsX, earningsY)
	common.Debugf("%s\n", strings.Join(earnings, "\n"))

	common.Debugf("- Deductions: [%d, %d]\n", deductionsX, deductionsY)
	common.Debugf("%s\n", strings.Join(deductions, "\n"))

	common.Debugf("- Taxes: [%d, %d]\n", taxesX, taxesY)
	common.Debugf("%s\n", strings.Join(taxes, "\n"))

	common.Debugf("- Pay Summary: [%d, %d]\n", summaryX, summaryY)
	common.Debugf("%s\n", strings.Join(summary, "\n"))

	//dataStr := string(data)
	//
	//lines := strings.Split(dataStr, "\n")
	//
	printCommaSeparated(
		filepath.Base(file),

		extractFirstField(header, `Period Start Date`, 0),
		extractFirstField(header, `Period End Date`, 0),
		extractFirstField(header, `Pay Date`, 0),
		extractFirstField(header, `Document`, 0),
		extractDollarsSum(header, `Net Pay`, 0, false).String(),

		extractDollarsSum(earnings, `Regular Pay`, 2, true).String(),
		extractDollarsSum(earnings, `Annual Bonus`, 0, true).String(),
		extractDollarsSum(earnings, `Group Term( Life)?`, 2, true).String(),
		extractDollarsSum(earnings, `Goog Stock( Unit)?`, 5, true).String(), // TODO -- Depends on job title
		extractDollarsSum(earnings, `Peer Bonus`, 0, true).String(),

		extractDollarsSum(deductions, `401K After-Tax`, 2, false).String(),
		extractDollarsSum(deductions, `401K Pretax`, 2, false).String(),
		extractDollarsSum(deductions, `Class C Offset`, 2, false).String(),
		extractDollarsSum(deductions, `Dental`, 2, false).String(),
		extractDollarsSum(deductions, `Group Term Life`, 2, false).String(),
		extractDollarsSum(deductions, `GSU C Refund`, 2, false).String(),
		extractDollarsSum(deductions, `Internet Reim`, 2, false).String(),
		extractDollarsSum(deductions, `LegalAccess`, 2, false).String(),
		extractDollarsSum(deductions, `LongTerm Dis`, 2, false).String(),
		extractDollarsSum(deductions, `Medical`, 2, false).String(),
		extractDollarsSum(deductions, `Vision`, 2, false).String(),
		extractDollarsSum(deductions, `Vol Life EE`, 2, false).String(),

		extractDollarsSum(taxes, `Federal Income Tax`, 1, false).String(),
		extractDollarsSum(header, `Additional Federal Income Tax`, 0, false).String(),
		extractDollarsSum(taxes, `Employee Medicare`, 1, false).String(),
		extractDollarsSum(taxes, `Social Security Employee Tax`, 1, false).String(),
		extractDollarsSum(taxes, `CA State Income Tax`, 1, false).String(),
		extractDollarsSum(taxes, `CA Private Disability Employee`, 1, false).String(),
	)

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

		`401K After-Tax`,
		`401K Pretax`,
		`Class C Offset`,
		`Dental`,
		`Group Term Life`,
		`GSU C Refund`,
		`Internet Reim`,
		`LegalAccess`,
		`LongTerm Dis`,
		`Medical`,
		`Vision`,
		`Vol Life EE`,

		`Federal Income Tax`,
		`Additional Federal Income Tax`,
		`Employee Medicare`,
		`Social Security Employee Tax`,
		`CA State Income Tax`,
		`CA Private Disability Employee`,
	)
}

func realMain() int {
	printHeader()
	for _, file := range os.Args[1:] {
		dumpPdf(file)
	}

	return 0
}
