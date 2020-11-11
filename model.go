package healthcheck

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/mitchellh/go-homedir"
)

// ResultItem ...
type ResultItem struct {
	Address Address
	Error   error
}

// PrintOK prints check OK
func (r *ResultItem) PrintOK() {
	fmt.Printf("%s %s\n", r.Address.Addr, r.Address.Desc)
}

// PrintError prints check error
func (r *ResultItem) PrintError() {
	fmt.Printf("%s %s, error %v\n", r.Address.Addr, r.Address.Desc, r.Error)
}

// ResultChan ...
type ResultChan struct {
	OKChan    chan ResultItem
	ErrorChan chan ResultItem
}

// NewResult makes a ResultChan
func NewResult() ResultChan {
	return ResultChan{
		OKChan:    make(chan ResultItem),
		ErrorChan: make(chan ResultItem),
	}
}

// Result is a structure of check results
type Result struct {
	TotalItems int
	OKItems    []ResultItem
	ErrorItems []ResultItem
}

// WaitResults waits all results returned.
func (r *ResultChan) WaitResults(totalItems int) Result {
	oks := make([]ResultItem, 0, totalItems)
	errs := make([]ResultItem, 0, totalItems)

	for i := 0; i < totalItems; i++ {
		select {
		case ok := <-r.OKChan:
			oks = append(oks, ok)
		case err := <-r.ErrorChan:
			errs = append(errs, err)
		}
	}

	return Result{TotalItems: totalItems, OKItems: oks, ErrorItems: errs}
}

// PrintResult prints result.
func (a *Result) PrintResult() {
	fmt.Printf("OK %d/%d:\n", len(a.OKItems), a.TotalItems)
	fmt.Printf("Failed %d/%d:\n", len(a.ErrorItems), a.TotalItems)

	if len(a.ErrorItems) > 0 {
		fmt.Printf("Failed Items:\n")
		for _, err := range a.ErrorItems {
			err.PrintError()
		}
	}

	if len(a.OKItems) > 0 {
		fmt.Printf("OK Items:\n")

		for _, ok := range a.OKItems {
			ok.PrintOK()
		}
	}
}

type Address struct {
	Addr string // 地址 ip:port
	Desc string // 描述
}

// Addresses represents an array of addresses
type Addresses []Address

// NewAddresses make a new Addresses
func NewAddresses() Addresses { return make([]Address, 0) }

// Len returns the length of address
func (a Addresses) Len() int { return len(a) }

// ParseAddress prepares the address
func (a *Addresses) ParseAddress(addrListFileName string, directAddrString string) {
	context := ""

	if addrListFileName != "" {
		expanded, _ := homedir.Expand(addrListFileName)
		if b, err := ioutil.ReadFile(expanded); err != nil {
			fmt.Printf("failed to read file %s, error: %v\n", addrListFileName, err)
		} else {
			context = string(b)
		}
	}

	a.parseLines(context + "\n" + directAddrString)
}

func (a *Addresses) parseLines(content string) {
	sep := regexp.MustCompile(`\d+(\.\d+)+:\d+`)

	idx := 0
	addr := ""

	for {
		found := sep.FindStringSubmatchIndex(content[idx:])
		if len(found) == 0 {
			break
		}

		if addr != "" {
			a.MergeAddresses(Address{
				Addr: addr,
				Desc: strings.TrimSpace(content[idx : idx+found[0]]),
			})
		}

		addr = strings.TrimSpace(content[idx+found[0] : idx+found[1]])
		idx += found[1]
	}

	if addr != "" {
		a.MergeAddresses(Address{
			Addr: addr,
			Desc: strings.TrimSpace(content[idx:]),
		})
	}
}

// MergeAddresses merges addresses.
func (a *Addresses) MergeAddresses(address Address) {
	for _, line := range *a {
		if line.Addr == address.Addr {
			fmt.Printf("duplicate address %s\n", line.Addr)
			return
		}
	}

	*a = append(*a, address)
}
