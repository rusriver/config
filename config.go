package config

import (
	"bytes"
	"strings"
)

// Config represents a configuration with convenient access methods.
type Config struct {
	DataSubTree            any
	OkPtr                  *bool
	ErrPtr                 *error
	ExpressionStatus       ExpressionFailure
	dontPanicFlag          bool
	Source                 *Source `json:"-"`
	relativePathFromParent []string
	parent                 *Config
}

type ExpressionFailure int

const (
	ExpressionStatus_0_Norm = iota
	ExpressionStatus_1_Failed
	ExpressionStatus_2_DefaultCallbackAlreadyUsedOnce
)

// Makes a copy of Config struct. The actual data is linked by reference.
func (c *Config) ChildCopy() (c2 *Config) {
	if c != nil {
		c2 = &Config{
			DataSubTree:            c.DataSubTree,
			OkPtr:                  c.OkPtr,
			ErrPtr:                 c.ErrPtr,
			ExpressionStatus:       c.ExpressionStatus,
			dontPanicFlag:          c.dontPanicFlag,
			Source:                 c.Source,
			relativePathFromParent: nil,
			parent:                 c,
		}
	} else {
		c2 = &Config{}
	}
	return
}

func (c *Config) GetCurrentLocationPlusPath(pathParts ...string) (path []string) {
	path = make([]string, 0, 10)
	var f func(*Config)
	f = func(c *Config) {
		if c.parent != nil {
			f(c.parent)
		}
		path = append(path, c.relativePathFromParent...)
	}
	f(c)
	path = append(path, pathParts...)
	return
}

// Traverses the struct down the path.
// P() does not create a path, if it didn't exist. So, if used before Set(),
// it will take you as far as there is something, not farther.
func (c *Config) P(pathParts ...string) *Config {
	c2 := c.ChildCopy()
	var err error
	c2.DataSubTree, err = goByPath(c2.DataSubTree, pathParts)
	if err != nil {
		c2.handleError(err)
	}
	c2.relativePathFromParent = pathParts
	return c2
}

func (c *Config) DotP(path string) *Config {
	return c.P(strings.Split(path, ".")...)
}

func (c *Config) SlashP(path string) *Config {
	return c.P(strings.Split(path, "/")...)
}

// Resets any errors, accumulated in previous expressions on this Config object.
// Sets Ok=true, Err=nil, ExpressionStatus=0_Norm, if any
func (c *Config) ErrOk() *Config {
	if c.OkPtr != nil {
		*c.OkPtr = true
	}
	if c.ErrPtr != nil {
		*c.ErrPtr = nil
	}
	c.ExpressionStatus = ExpressionStatus_0_Norm
	return c
}

// Sets ExpressionStatus=failed if it was OK; sets Err, if it wasn't already; sets Ok=false, if it's present;
// Then panics, unless there is present either of Ok, Err, or dontPanicFlag (it can be set with U()).
func (c *Config) handleError(err error) {
	if err == nil {
		return
	} else {
		if c.ExpressionStatus < ExpressionStatus_1_Failed {
			c.ExpressionStatus = ExpressionStatus_1_Failed
		}
		if c.ErrPtr != nil {
			if *c.ErrPtr == nil {
				// only set error, if there wasn't error already set;
				// this way, we keep the very first occurred error in there,
				// which is what we want (we don't want meaningless error-because-of-error errors)
				*c.ErrPtr = err
			}
		}
		if c.OkPtr != nil {
			*c.OkPtr = false
		}
		if c.ErrPtr == nil && c.OkPtr == nil && !c.dontPanicFlag {
			panic(err)
		}
	}
}

func (c *Config) isExpressionOk() (ok bool) {
	if c.ErrPtr != nil {
		return *c.ErrPtr == nil
	}
	if c.OkPtr != nil {
		return *c.OkPtr == true
	}
	return c.ExpressionStatus == ExpressionStatus_0_Norm
}

// Sets a nested config according to a path, relative from current location.
// Without a Source object set, acts as a by-pass to the NonThreadSafe_Set().
// If you don't want to specify path, and just want to use it from current location,
// then invoke with nil path, c.Set(nil, value)
// It is totally asynchronous, and it's effect is somewhat delayed. Use explicit
// flush signal if you want to synchronize explicitly.
// P() does not create a path, if it didn't exist. So, if used before Set(),
// it will take you as far as there is something, not farther.
func (c *Config) Set(pathParts []string, v interface{}) {
	if c.Source == nil {
		// Current location is implicit, by Config object
		c.NonThreadSafe_Set(pathParts, v)
		return
	} else {
		// We need an absolute full path to current location, in this case.
		// Current location is assembled by traversing the up-links to parents, and getting
		// all paths you ever went down with P(). This is thread-safe operation.

		loc := c.GetCurrentLocationPlusPath(pathParts...)

		msg := &MsgCmd{
			Command:  Command_Set,
			FullPath: loc,
			V:        v,
		}
		if len(c.Source.ChCmd) >= cap(c.Source.ChCmd)/10*7 {
			// Please look at 20230618-go-tests/3 for explanation.
			// Also, we signal on 70%, so while the WBUG does deep copy, there's still a room
			// for more commands.
			select {
			case c.Source.ChFlushSignal <- &MsgFlushSignal{}:
			default:
			}
		}
		c.Source.ChCmd <- msg
	}
}

// Sets a nested config according to a path, relative from current location.
// If you don't want to specify path, and just want to use it from current location,
// then invoke with nil path.
func (c *Config) NonThreadSafe_Set(pathParts []string, v interface{}) {
	err := set(c.DataSubTree, pathParts, v)
	if err != nil {
		c.handleError(err)
	}
}

// Sets dontPanicFlag=true, so that failing operations won't panic, if there's no Err or Ok set.
func (c *Config) U() (c2 *Config) {
	c2 = c.ChildCopy()
	c2.dontPanicFlag = true
	return c2
}

// Reverse of U()
func (c *Config) UnU() (c2 *Config) {
	c2 = c.ChildCopy()
	c2.dontPanicFlag = false
	return c2
}

// Attaches a ok bool variable to the expression, by reference, with it's current value.
// Failures in subsequent operations may set it to false only, so make sure its current state is true.
func (c *Config) Ok(okRef *bool) (c2 *Config) {
	c2 = c.ChildCopy()
	c2.OkPtr = okRef
	return c2
}

// Attaches an err error variable to the expression, by reference. Make sure to reset it to nil
// when reusing between expressions.
func (c *Config) Err(err *error) (c2 *Config) {
	c2 = c.ChildCopy()
	c2.ErrPtr = err
	return c2
}

func SplitPathToParts(key string) []string {
	parts := []string{}

	bracketOpened := false
	var buffer bytes.Buffer
	for _, char := range key {
		if char == 91 || char == 93 { // [ ]
			bracketOpened = char == 91
			continue
		}
		if char == 46 && !bracketOpened { // point
			parts = append(parts, buffer.String())
			buffer.Reset()
			continue
		}

		buffer.WriteRune(char)
	}

	if buffer.String() != "" {
		parts = append(parts, buffer.String())
		buffer.Reset()
	}

	return parts
}
