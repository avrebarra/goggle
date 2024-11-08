package do

import "fmt"

func JoinErrors(errs []error) (err error) {
	if len(errs) == 0 {
		return
	}
	msg := ""
	for _, err := range errs {
		if err == nil {
			continue
		}
		msg += err.Error() + "; "
	}
	if msg == "" {
		return nil
	}
	msg = msg[:len(msg)-2]
	return fmt.Errorf(msg)
}
