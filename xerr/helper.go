package xerr

import (
	"fmt"
	"strings"
)

// JointErrors combines multiple errors into a single error while preserving traceability
func JointErrors(errs ...error) error {
	if len(errs) == 0 {
		return nil
	}

	var nonNilErrs []error
	for _, err := range errs {
		if err != nil {
			nonNilErrs = append(nonNilErrs, err)
		}
	}

	if len(nonNilErrs) == 0 {
		return nil
	}

	if len(nonNilErrs) == 1 {
		return nonNilErrs[0]
	}

	var sb strings.Builder
	sb.WriteString("multiple errors occurred:\n")
	for i, err := range nonNilErrs {
		sb.WriteString(fmt.Sprintf("[%d] %v\n", i+1, err))
	}
	return fmt.Errorf("%s", sb.String())
}
