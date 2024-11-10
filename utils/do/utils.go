package do

import "fmt"

type JointError struct {
	Errors []error
}

func (e JointError) Error() string {
	count := 0
	msg := ""
	for _, err := range e.Errors {
		if err == nil {
			continue
		}
		msg += err.Error() + "; "
		count++
	}
	if msg == "" {
		return "empty error"
	}
	msg = msg[:len(msg)-2]
	msg = fmt.Sprintf("%d joint error: %s", count, msg)
	return msg
}

func JoinErrors(errs []error) (err error) {
	je := JointError{Errors: []error{}}
	for _, v := range errs {
		if v != nil {
			je.Errors = append(je.Errors, v)
		}
	}
	if len(je.Errors) == 0 {
		return nil
	}
	err = je
	return
}
