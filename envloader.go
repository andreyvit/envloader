// Package envloader assists in writing executables that load their config from the environment.
package envloader

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

// Required is a convenient value to pass to VarSet.Add for variables that are always required.
func Required() bool {
	return true
}

// Optional is a convenient value to pass to VarSet.Add for variables that are always optional.
func Optional() bool {
	return false
}

// WhenTrue returns a value to pass to VarSet.Add for variables that are required when the given variable is true.
func WhenTrue(boolVar *bool) func() bool {
	return func() bool {
		return *boolVar
	}
}

// WhenFalse returns a value to pass to VarSet.Add for variables that are required when the given variable is false.
func WhenFalse(boolVar *bool) func() bool {
	return func() bool {
		return !*boolVar
	}
}

// Var defines a single environment variable.
type Var struct {
	EnvKey   string
	Required func() bool
	Value    flag.Value
	Desc     string

	IsSpecified bool
}

// VarSet is a slice of environment variable definitions. The ordering matters,
// both when printing the values (obviously), and also when parsing, because
// later variables can refer to the values of prior ones.
type VarSet []*Var

// Var adds a given value to the set of environment variable definitions.
//
// Required func specifies the conditions when the value is required. Very often,
// a feature is enabled by a master on/off variable, and a bunch of configuration
// parameters are only required when it is on. Use Required, Optional, WhenTrue
// and WhenFalse helpers, or define your own custom function. These functions
// can refer to the values of the previously defined variables.
//
// Use StringVar, BoolVar, IntVar & similar helpers defined in this package
// to make flag.Value for your variables.
func (vars *VarSet) Var(envKey string, required func() bool, value flag.Value, desc string) *Var {
	v := &Var{
		EnvKey:   envKey,
		Required: required,
		Value:    value,
		Desc:     desc,
	}
	*vars = append(*vars, v)
	return v
}

// String returns a shell script that defines all variables in the set.
// Variable descriptions are added as comments.
func (vars VarSet) String() string {
	var buf strings.Builder
	vars.PrintTo(&buf)
	return buf.String()
}

// Print prints a shell script that defines all variables in the set to os.Stdout.
// Variable descriptions are added as comments.
func (vars VarSet) Print() {
	vars.PrintTo(os.Stdout)
}

// PrintTo prints a shell script that defines all variables in the set.
// Variable descriptions are added as comments.
func (vars VarSet) PrintTo(out io.Writer) {
	for _, vr := range vars {
		usage := vr.Desc
		if usage != "" {
			usage = "# " + strings.ReplaceAll(usage, "\n", "\n# ") + "\n"
		}

		valueStr := vr.Value.String()
		if valueStr == "" {
			valueStr = "..."
		}

		fmt.Fprintf(out, "%s%s=%s\n", usage, vr.EnvKey, valueStr)
		// fmt.Fprintf(out, "  %s\n    \t%s\n", vr.EnvKey, strings.ReplaceAll(usage.String(), "\n", "\n    \t"))
	}
}

// Parse parses the current environment variable values. If parsing fails,
// prints an error message and exits the program with error code 2.
func (vars VarSet) Parse() {
	e := vars.TryParse()
	if e != nil {
		PrintError(e, os.Stderr)
		os.Exit(2)
	}
}

// TryParse parses the current environment variable values.
// Returns nil when successful, a pointer to Error when not.
func (vars VarSet) TryParse() *Error {
	return vars.TryParseFrom(os.Getenv)
}

// TryParseFrom parses environment variable values returned by the given function.
// Returns nil when successful, a pointer to Error when not.
func (vars VarSet) TryParseFrom(getenv func(string) string) *Error {
	var e *Error

	for _, vr := range vars {
		raw := getenv(vr.EnvKey)
		if raw != "" {
			err := vr.Value.Set(raw)
			if err != nil {
				if e == nil {
					e = &Error{}
				}
				e.InvalidValues = append(e.InvalidValues, &InvalidValue{vr.EnvKey, err})
				continue
			}
			vr.IsSpecified = true
		}
	}

	for _, vr := range vars {
		if !vr.IsSpecified && vr.Required() {
			if e == nil {
				e = &Error{}
			}
			e.MissingVars = append(e.MissingVars, vr)
		}
	}

	return e
}

// Error describes environment variable problems encountered by TryParse.
type Error struct {
	InvalidValues []*InvalidValue
	MissingVars   VarSet
}

// PrintError performs default printing of the given error returned by TryParse.
func PrintError(e *Error, w io.Writer) {
	for _, iv := range e.InvalidValues {
		fmt.Fprintf(w, "** %s\n", iv.Error())
	}
	if len(e.MissingVars) > 1 {
		fmt.Fprintf(w, "** missing values for the following %d environment variables:\n%s\n", len(e.MissingVars), e.MissingVars.String())
	} else if len(e.MissingVars) == 1 {
		fmt.Fprintf(w, "** missing value for the following environment variable:\n%s\n", e.MissingVars.String())
	}
}

// InvalidValue is an error returned as part of Error struct for environment variable values that failed to parse.
type InvalidValue struct {
	EnvKey string
	Cause  error
}

func (e *InvalidValue) Unwrap() error {
	return e.Cause
}

func (e *InvalidValue) Error() string {
	return fmt.Sprintf("invalid value of environment variable %s: %v", e.EnvKey, e.Cause)
}

// PrintAction returns flag.Value that can be used with flag.Var to print all environment variables in shell format.
//
// Use like this:
//
//	flag.Var(vars.PrintAction(), "print-env", "print all supported environment variables in shell format")
func (vars VarSet) PrintAction() flag.Value {
	return printAction(vars)
}

type printAction VarSet

func (_ printAction) String() string {
	return ""
}

func (_ printAction) IsBoolFlag() bool {
	return true
}

func (a printAction) Set(string) error {
	VarSet(a).PrintTo(os.Stdout)
	return flag.ErrHelp
}
