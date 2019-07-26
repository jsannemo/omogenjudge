// Problem testdata verification.
package testdata

import (
	toolspb "github.com/jsannemo/omogenjudge/problemtools/api"
	"github.com/jsannemo/omogenjudge/problemtools/util"
)

// VerifyTestdata verifies the test data of the given problem.
func VerifyTestdata(problem *toolspb.Problem, reporter util.Reporter) error {
	// TODO verify correctness of testdata
	// TODO check groups are not empty
	return nil
}
