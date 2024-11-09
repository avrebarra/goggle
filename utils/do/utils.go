package do

import "fmt"

type JointError struct {
	Errors []error
}

func (e JointError) Error() string {
	msg := ""
	for _, err := range e.Errors {
		if err == nil {
			continue
		}
		msg += err.Error() + "; "
	}
	if msg == "" {
		return ""
	}
	msg = msg[:len(msg)-2]
	msg = fmt.Sprintf("%d joint error: %s", len(e.Errors), msg)
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
