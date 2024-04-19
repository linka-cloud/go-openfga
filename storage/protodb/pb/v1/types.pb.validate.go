// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: v1/types.proto

package pbv1

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort
)

// Validate checks the field values on Tuple with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Tuple) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Tuple with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in TupleMultiError, or nil if none found.
func (m *Tuple) ValidateAll() error {
	return m.validate(true)
}

func (m *Tuple) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Key

	// no validation rules for StoreId

	if all {
		switch v := interface{}(m.GetTupleKey()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, TupleValidationError{
					field:  "TupleKey",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, TupleValidationError{
					field:  "TupleKey",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetTupleKey()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return TupleValidationError{
				field:  "TupleKey",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetCreatedAt()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, TupleValidationError{
					field:  "CreatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, TupleValidationError{
					field:  "CreatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetCreatedAt()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return TupleValidationError{
				field:  "CreatedAt",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return TupleMultiError(errors)
	}

	return nil
}

// TupleMultiError is an error wrapping multiple validation errors returned by
// Tuple.ValidateAll() if the designated constraints aren't met.
type TupleMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m TupleMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m TupleMultiError) AllErrors() []error { return m }

// TupleValidationError is the validation error returned by Tuple.Validate if
// the designated constraints aren't met.
type TupleValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e TupleValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e TupleValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e TupleValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e TupleValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e TupleValidationError) ErrorName() string { return "TupleValidationError" }

// Error satisfies the builtin error interface
func (e TupleValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sTuple.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = TupleValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = TupleValidationError{}

// Validate checks the field values on Assertions with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Assertions) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Assertions with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in AssertionsMultiError, or
// nil if none found.
func (m *Assertions) ValidateAll() error {
	return m.validate(true)
}

func (m *Assertions) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Key

	// no validation rules for StoreId

	// no validation rules for ModelId

	for idx, item := range m.GetAssertions() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, AssertionsValidationError{
						field:  fmt.Sprintf("Assertions[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, AssertionsValidationError{
						field:  fmt.Sprintf("Assertions[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return AssertionsValidationError{
					field:  fmt.Sprintf("Assertions[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return AssertionsMultiError(errors)
	}

	return nil
}

// AssertionsMultiError is an error wrapping multiple validation errors
// returned by Assertions.ValidateAll() if the designated constraints aren't met.
type AssertionsMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m AssertionsMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m AssertionsMultiError) AllErrors() []error { return m }

// AssertionsValidationError is the validation error returned by
// Assertions.Validate if the designated constraints aren't met.
type AssertionsValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e AssertionsValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e AssertionsValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e AssertionsValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e AssertionsValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e AssertionsValidationError) ErrorName() string { return "AssertionsValidationError" }

// Error satisfies the builtin error interface
func (e AssertionsValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sAssertions.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = AssertionsValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = AssertionsValidationError{}

// Validate checks the field values on Model with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Model) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Model with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in ModelMultiError, or nil if none found.
func (m *Model) ValidateAll() error {
	return m.validate(true)
}

func (m *Model) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Key

	// no validation rules for StoreId

	if all {
		switch v := interface{}(m.GetModel()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, ModelValidationError{
					field:  "Model",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, ModelValidationError{
					field:  "Model",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetModel()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return ModelValidationError{
				field:  "Model",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return ModelMultiError(errors)
	}

	return nil
}

// ModelMultiError is an error wrapping multiple validation errors returned by
// Model.ValidateAll() if the designated constraints aren't met.
type ModelMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ModelMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ModelMultiError) AllErrors() []error { return m }

// ModelValidationError is the validation error returned by Model.Validate if
// the designated constraints aren't met.
type ModelValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ModelValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ModelValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ModelValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ModelValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ModelValidationError) ErrorName() string { return "ModelValidationError" }

// Error satisfies the builtin error interface
func (e ModelValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sModel.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ModelValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ModelValidationError{}

// Validate checks the field values on Change with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Change) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Change with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in ChangeMultiError, or nil if none found.
func (m *Change) ValidateAll() error {
	return m.validate(true)
}

func (m *Change) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Key

	// no validation rules for StoreId

	if all {
		switch v := interface{}(m.GetChange()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, ChangeValidationError{
					field:  "Change",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, ChangeValidationError{
					field:  "Change",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetChange()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return ChangeValidationError{
				field:  "Change",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return ChangeMultiError(errors)
	}

	return nil
}

// ChangeMultiError is an error wrapping multiple validation errors returned by
// Change.ValidateAll() if the designated constraints aren't met.
type ChangeMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ChangeMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ChangeMultiError) AllErrors() []error { return m }

// ChangeValidationError is the validation error returned by Change.Validate if
// the designated constraints aren't met.
type ChangeValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ChangeValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ChangeValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ChangeValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ChangeValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ChangeValidationError) ErrorName() string { return "ChangeValidationError" }

// Error satisfies the builtin error interface
func (e ChangeValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sChange.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ChangeValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ChangeValidationError{}
