package arc

import "errors"

type ErrorResult struct {
	err  error
	data interface{}
}

func NewErrorResult(err error, data interface{}) *ErrorResult {
	return &ErrorResult{err: err, data: data}
}

func (e *ErrorResult) Unwrap() interface{} {
	if e.err != nil {
		panic(e.err)
	}
	return e.data
}

func (e *ErrorResult) UnwrapOr(v interface{}) interface{} {
	if e.err != nil {
		return v
	}
	return e.data
}

func (e *ErrorResult) UnwrapOrElse(f func() interface{}) interface{} {
	if e.err != nil {
		return f()
	}
	return e.data
}

// 0 err  1 data
func Result(vs ...interface{}) *ErrorResult {
	if len(vs) == 1 {
		if vs[0] == nil {
			return &ErrorResult{
				err:  nil,
				data: nil,
			}
		}
		if e, ok := vs[0].(error); ok {
			return &ErrorResult{
				err:  e,
				data: nil,
			}
		}
	}
	if len(vs) == 2 {
		if vs[0] == nil {
			return &ErrorResult{
				err:  nil,
				data: vs[1],
			}
		}
		if e, ok := vs[0].(error); ok {
			return &ErrorResult{
				err:  e,
				data: vs[1],
			}
		}
	}

	return &ErrorResult{
		err:  errors.New("error result format"),
		data: nil,
	}

}

type BindFunc func(v interface{}) error

func Exec(f BindFunc, value interface{}) *ErrorResult {
	err := f(value)
	return NewErrorResult(err, value)
}
