package vaultclient

import (
	"fmt"
)

type ErrParamMinLen struct {
	field   string
	context string
	Min     int
}

func (e *ErrParamMinLen) Error() string {
	return fmt.Sprintf("parameter validation failed: %s", e.context)
}

func (e *ErrParamMinLen) Field() string {
	return e.field
}

func (e *ErrParamMinLen) SetContext(ctx string) {
	e.context = ctx
}

func (e *ErrParamMinLen) AddNestedContext(ctx string) {
	if e.context != "" {
		e.context = ctx + "." + e.context
	} else {
		e.context = ctx
	}
}

type ErrParamMinValue struct {
	field   string
	context string
	Min     int64
}

func (e *ErrParamMinValue) Error() string {
	return fmt.Sprintf("parameter validation failed: %s", e.context)
}

func (e *ErrParamMinValue) Field() string {
	return e.field
}

func (e *ErrParamMinValue) SetContext(ctx string) {
	e.context = ctx
}

func (e *ErrParamMinValue) AddNestedContext(ctx string) {
	if e.context != "" {
		e.context = ctx + "." + e.context
	} else {
		e.context = ctx
	}
}

type ErrParamMaxValue struct {
	field   string
	context string
	Max     int64
}

func (e *ErrParamMaxValue) Error() string {
	return fmt.Sprintf("parameter validation failed: %s", e.context)
}

func (e *ErrParamMaxValue) Field() string {
	return e.field
}

func (e *ErrParamMaxValue) SetContext(ctx string) {
	e.context = ctx
}

func (e *ErrParamMaxValue) AddNestedContext(ctx string) {
	if e.context != "" {
		e.context = ctx + "." + e.context
	} else {
		e.context = ctx
	}
}

func NewErrParamMinLen(param string, minLen int) *ErrParamMinLen {
	return &ErrParamMinLen{
		field:   param,
		context: fmt.Sprintf("%s must be at least %d character(s)", param, minLen),
		Min:     minLen,
	}
}

func NewErrParamMinValue(param string, minVal int64) *ErrParamMinValue {
	return &ErrParamMinValue{
		field:   param,
		context: fmt.Sprintf("%s must be at least %d", param, minVal),
		Min:     minVal,
	}
}

func NewErrParamMaxValue(param string, maxVal int64) *ErrParamMaxValue {
	return &ErrParamMaxValue{
		field:   param,
		context: fmt.Sprintf("%s must not exceed %d", param, maxVal),
		Max:     maxVal,
	}
}
