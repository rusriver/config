package config

import (
	"bytes"
)

// Config represents a configuration with convenient access methods.
type Config struct {
	DataSubTree            any
	lastError              error
	ok                     *bool
	err                    *error
	dontPanicFlag          bool
	Source                 *Source `json:"-"`
	relativePathFromParent []string
	parent                 *Config
}

func (c *Config) ChildCopy() *Config {
	return &Config{
		DataSubTree:            c.DataSubTree,
		lastError:              c.lastError,
		ok:                     c.ok,
		err:                    c.err,
		dontPanicFlag:          c.dontPanicFlag,
		Source:                 c.Source,
		relativePathFromParent: c.relativePathFromParent,
		parent:                 c,
	}
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

func (c *Config) LastError() error {
	return c.lastError
}

// P() does not create a path, if it didn't exist. So, if used with Set(),
// it will take you as far as there is something, not farther.
func (c *Config) P(pathParts ...string) *Config {
	c.resetErrOkState()
	c2 := c.ChildCopy()
	var err error
	c2.DataSubTree, err = goByPath(c2.DataSubTree, pathParts)
	if err != nil {
		c2.handleError(err)
	}
	c2.relativePathFromParent = pathParts
	return c2
}

func (c *Config) resetErrOkState() {
	c.lastError = nil
	if c.ok != nil {
		*c.ok = true
	}
	if c.err != nil {
		*c.err = nil
	}

}

func (c *Config) handleError(err error) {
	if err == nil {
		return
	} else {
		c.lastError = err
		if c.err != nil {
			*c.err = err
		}
		if c.ok != nil {
			*c.ok = false
		}
		if c.err == nil && c.ok == nil && !c.dontPanicFlag {
			panic(err)
		}
	}
}

// Sets a nested config according to a path, relative from current location.
// Without a Source object set, acts a by-pass to the NonThreadSafe_Set().
// If you don't want to specify path, and just want to use it from current location,
// then invoke with nil path, c.Set(nil, value)
// It is totally asynchronous, and it's effect is somewhat delayed. Use explicit
// flush signal if you want to synchronize explicitly.
// P() does not create a path, if it didn't exist. So, if used with Set(),
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
	c.resetErrOkState()
	err := set(c.DataSubTree, pathParts, v)
	if err != nil {
		c.handleError(err)
	}
}

func (c *Config) U() (c2 *Config) {
	c.resetErrOkState()
	c2 = c.ChildCopy()
	c2.dontPanicFlag = true
	return c2
}

func (c *Config) NotU() (c2 *Config) {
	c.resetErrOkState()
	c2 = c.ChildCopy()
	c2.dontPanicFlag = false
	return c2
}

func (c *Config) Ok(okRef *bool) (c2 *Config) {
	c.resetErrOkState()
	c2 = c.ChildCopy()
	c2.ok = okRef
	return c2
}

func (c *Config) E(err *error) (c2 *Config) {
	c.resetErrOkState()
	c2 = c.ChildCopy()
	c2.err = err
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
