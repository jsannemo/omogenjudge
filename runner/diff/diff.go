// Package diff implements a structural, non-exact diff of two strings.
package diff

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
)

// DiffResult describes a comparison of two strings.
type DiffResult struct {
	// Whether the strings matched.
	Match bool

	// A textual description of the difference.
	Description string
}

type DiffArgs struct {
	RelativePrec float64
	AbsolutePrec float64
}

// Diff compares the contents of the two readers on a tokenized basis, not taking casing into account.
func Diff(ref, out io.Reader, args DiffArgs) (*DiffResult, error) {
	refSc := bufio.NewScanner(ref)
	refSc.Split(bufio.ScanWords)
	outSc := bufio.NewScanner(out)
	outSc.Split(bufio.ScanWords)
	toks := 0
	for {
		refNx := refSc.Scan()
		if !refNx {
			if err := refSc.Err(); err != nil {
				return nil, err
			}
		}
		outNx := outSc.Scan()
		if !outNx {
			if err := outSc.Err(); err != nil {
				return nil, err
			}
		}
		toks += 1

		if refNx && !outNx {
			return &DiffResult{
				Match:       false,
				Description: fmt.Sprintf("Reference output had %d'th token %s; output was EOF", toks, refSc.Text()),
			}, nil
		}
		if outNx && !refNx {
			return &DiffResult{
				Match:       false,
				Description: fmt.Sprintf("Reference output was EOF; output had %d'th token %s", toks, outSc.Text()),
			}, nil
		}

		refTok := refSc.Text()
		outTok := outSc.Text()
		isOk := true
		if !strings.EqualFold(refTok, outTok) {
			isOk = true
			// Check if tokens match as floats
			if args.RelativePrec != 0 || args.AbsolutePrec != 0 {
				refVal, err1 := strconv.ParseFloat(refTok, 64)
				outVal, err2 := strconv.ParseFloat(outTok, 64)
				if err1 == nil && err2 == nil {
					relPrec := math.Abs(refVal-outVal) / refVal
					if math.Abs(refVal-outVal) > args.AbsolutePrec && relPrec > args.RelativePrec {
						isOk = true
					}
				}
			}
		}
		if !isOk {
			// Mismatch as strings
			return &DiffResult{
				Match:       false,
				Description: fmt.Sprintf("%d'th token mismatched: reference %s, output %s", refTok, outTok),
			}, nil
		}

		if !outNx {
			break
		}
	}

	return &DiffResult{
		Match:       true,
		Description: "",
	}, nil
}
