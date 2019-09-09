package healthcheck

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/mitchellh/go-homedir"
)

type ResultItem struct {
	Address string
	Error   error
}

func (r *ResultItem) PrintOK() {
	fmt.Printf("%s\n", r.Address)
}

func (r *ResultItem) PrintError() {
	fmt.Printf("%s, error %v\n", r.Address, r.Error)
}

type ResultChan struct {
	OKChan    chan ResultItem
	ErrorChan chan ResultItem
}

func NewResult() ResultChan {
	return ResultChan{
		OKChan:    make(chan ResultItem),
		ErrorChan: make(chan ResultItem),
	}
}

type Result struct {
	TotalItems int
	OKItems    []ResultItem
	ErrorItems []ResultItem
}

func (r *ResultChan) WaitResults(totalItems int) Result {
	oks := make([]ResultItem, 0)
	errs := make([]ResultItem, 0)
	cnt := totalItems
	for cnt > 0 {
		select {
		case ok := <-r.OKChan:
			oks = append(oks, ok)
			cnt--
		case err := <-r.ErrorChan:
			errs = append(errs, err)
			cnt--
		}
	}

	return Result{TotalItems: totalItems, OKItems: oks, ErrorItems: errs}
}

func (a *Result) PrintResult() {
	hasErrors := len(a.ErrorItems) > 0
	if hasErrors {
		fmt.Printf("Failed %d/%d:\n", len(a.ErrorItems), a.TotalItems)
		for _, err := range a.ErrorItems {
			err.PrintError()
		}
	}

	if len(a.OKItems) > 0 {
		if hasErrors {
			fmt.Printf("\n")
		}

		fmt.Printf("OK %d/%d:\n", len(a.OKItems), a.TotalItems)
		for _, ok := range a.OKItems {
			ok.PrintOK()
		}
	}
}

type Addresses []string

func NewAddresses() Addresses {
	return make([]string, 0)
}

func (a Addresses) Len() int {
	return len(a)
}

func (a *Addresses) PrepareAddress(addrListFileName string, directAddrString string) {
	sep := regexp.MustCompile(`[^\d.:]+`)

	if addrListFileName != "" {
		expanded, _ := homedir.Expand(addrListFileName)
		if b, err := ioutil.ReadFile(expanded); err != nil {
			fmt.Printf("failed to read file %s, error: %v\n", addrListFileName, err)
		} else {
			lines := sep.Split(string(b), -1)
			a.MergeAddresses(lines)
		}
	}

	if directAddrString != "" {
		lines := sep.Split(directAddrString, -1)
		a.MergeAddresses(lines)
	}
}

func (a *Addresses) MergeAddresses(mergedAddresses Addresses) {
	for _, line := range mergedAddresses {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if StringSliceContains(*a, line) {
			fmt.Printf("duplicate address %s\n", line)
		} else {
			*a = append(*a, line)
		}
	}
}
