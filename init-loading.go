package config

import (
	"errors"
	"path/filepath"
	"regexp"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type InitContext struct {
	FileName string
	Data     []byte
	Logger   *zerolog.Logger
	err      *error
	ok       *bool
}

func (ic *InitContext) FromFile(fileName string) *InitContext {
	ic.FileName = fileName
	return ic
}

func (ic *InitContext) FromBytes(data []byte) *InitContext {
	ic.Data = data
	return ic
}

func (ic *InitContext) WithLogger(logger *zerolog.Logger) *InitContext {
	ic.Logger = logger
	return ic
}

func (ic *InitContext) E(err *error) *InitContext {
	ic.err = err
	return ic
}

func (ic *InitContext) Ok(ok *bool) *InitContext {
	ic.ok = ok
	return ic
}

var reSuffixYaml = regexp.MustCompile(`\.[Yy][Aa]?[Mm][Ll]\s*$`)
var reSuffixJson = regexp.MustCompile(`\.(JSON|json)\s*$`)

func (ic *InitContext) Load() *Config {
	var c *Config
	var err error

	func() {
		switch {
		case len(ic.Data) > 0:
			// c, err = parseSerk(ic.Data)
			// if err == nil {
			// 	return
			// }
			c, err = parseYaml(ic.Data)
			if err == nil {
				return
			}
			c, err = parseJson(ic.Data)
			return

		case len(ic.FileName) > 0:
			switch {
			case reSuffixYaml.MatchString(ic.FileName) == true:
				c, err = parseYamlFile(ic.FileName)
				return
			case reSuffixJson.MatchString(ic.FileName) == true:
				c, err = parseJsonFile(ic.FileName)
				return
			default:
				err = errors.New("unknown file suffix")
				return
			}
		default:
			err = errors.New("data or file not specified")
			return
		}
	}()

	// this doe inherit these..
	c.err = ic.err
	c.ok = ic.ok

	if err == nil {
		if ic.ok != nil {
			*ic.ok = true
		}
	} else {
		if ic.err != nil {
			*ic.err = err
		}
		if ic.ok != nil {
			*ic.ok = false
		}
		if ic.err == nil && ic.ok == nil {
			panic(err)
		}
	}
	return c
}

func (ic *InitContext) LoadWithParenting() (result *Config) {
	if ic.Logger == nil {
		ic.Logger = &log.Logger
	}
	ic.Logger.Info().Msgf("ziPdTJw: reading the config file(s)...")
	baseDir := filepath.Dir(ic.FileName)
	filesAlreadyRead := map[string]bool{}
	isRoot := true
	var readNext func(configFileName string) *Config
	readNext = func(configFileName string) *Config {
		ic.Logger.Info().Msgf("EZWLkX: reading the config file '%v'...", configFileName)
		filesAlreadyRead[configFileName] = true
		var err error
		c1 := (&InitContext{FileName: configFileName}).E(&err).Load()
		if err != nil {
			ic.Logger.Err(err).Msgf("fYmNdkUt: config.ParseYamlFile('%v') failed", configFileName)
			panic(err)
		}
		if isRoot {
			isRoot = false
			id := c1.E(&err).P("id").String()
			ic.Logger.Info().Msgf("KPPEY7ZW: config file '%v' id='%v' err='%v'", configFileName, id, err)
		}
		parents := []string{}
		p := c1.P("parent").String()
		if len(p) > 0 {
			parents = append(parents, p)
		}
		if list := c1.P("parents").ListString(); err == nil {
			parents = append(parents, list...)
		}
		for _, configFileName := range parents {
			path := baseDir + "/" + configFileName
			if filesAlreadyRead[path] {
				ic.Logger.Err(err).Msgf("AweL9D: config file loop: the file '%v' already read", path)
				panic(err)
			}
			cN := readNext(path)
			cN.ExtendBy_v2(c1)
			c1 = cN
		}
		return c1
	}
	result = readNext(ic.FileName)
	result.Set([]string{"parent"}, nil)
	result.Set([]string{"parents"}, nil)
	ic.Logger.Info().Msg("K2aUDgz: reading the config file(s) OK")
	return
}
