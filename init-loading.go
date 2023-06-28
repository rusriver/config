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
	ErrPtr   *error
	OkPtr    *bool
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

func (ic *InitContext) Err(err *error) *InitContext {
	ic.ErrPtr = err
	return ic
}

func (ic *InitContext) Ok(ok *bool) *InitContext {
	ic.OkPtr = ok
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

	if err != nil {
		if ic.ErrPtr != nil {
			*ic.ErrPtr = err
		}
		if ic.OkPtr != nil {
			*ic.OkPtr = false
		}
		if ic.ErrPtr == nil && ic.OkPtr == nil {
			panic(err)
		}
		return nil
	}

	// this does inherit these..
	c.ErrPtr = ic.ErrPtr
	c.OkPtr = ic.OkPtr

	return c
}

func (ic *InitContext) LoadWithParenting() (result *Config) {
	if ic.Logger == nil {
		ic.Logger = &log.Logger
	}
	ic.Logger.Info().Msgf("ziPdTJw: reading the config file(s)...")
	filesAlreadyRead := map[string]bool{}
	isRoot := true
	depth := 0
	var readParent func(baseDir, configFileName string) *Config
	readParent = func(baseDir, currConfigFileName string) *Config {
		depth--
		defer func() { depth++ }()
		logger := ic.Logger.With().Int("depth", depth).Logger()
		logger.Info().Msgf("EZWLkX: reading the config file '%v'...", currConfigFileName)
		filesAlreadyRead[currConfigFileName] = true
		var err error
		conf := (&InitContext{FileName: currConfigFileName}).Err(&err).Load()
		if err != nil {
			logger.Err(err).Msgf("fYmNdkUt: config.ParseYamlFile('%v') failed", currConfigFileName)
			panic(err)
		}
		if isRoot {
			isRoot = false
			id := conf.ErrOk().P("id").String()
			logger.Info().Msgf("KPPEY7ZW: config file '%v' id='%v' err='%v'", currConfigFileName, id, err)
		}
		parents := []string{}
		ok := true
		p1 := conf.Ok(&ok).P("parent").String()
		if ok {
			parents = append(parents, p1)
		}
		list := conf.P("parents").ListString()
		parents = append(parents, list...)
		var aggregatedParentConf *Config
		for _, parentConfigFileName := range parents {
			parentFullPath := baseDir + "/" + parentConfigFileName
			if filesAlreadyRead[parentFullPath] {
				logger.Err(err).Msgf("AweL9D: config file loop: the file '%v' already read", parentFullPath)
				panic(err)
			}
			confParent := readParent(filepath.Dir(parentFullPath), parentFullPath)
			if aggregatedParentConf == nil {
				logger.Info().Msgf("KUY76-1: set aggregated parent from '%v'", parentFullPath)
				aggregatedParentConf = confParent
			} else {
				logger.Info().Msgf("KUY76-2: extend aggregated parent with '%v'", parentFullPath)
				aggregatedParentConf.ExtendBy_v2(confParent)
			}
		}
		if aggregatedParentConf != nil {
			logger.Info().Msgf("KUY76-3: extend aggregated parent with '%v' and return it", currConfigFileName)
			aggregatedParentConf.ExtendBy_v2(conf)
			conf = aggregatedParentConf
		}
		return conf
	}
	result = readParent(filepath.Dir(ic.FileName), ic.FileName)
	result.Set([]string{"parent"}, nil)
	result.Set([]string{"parents"}, nil)
	ic.Logger.Info().Msg("K2aUDgz: reading the config file(s) OK")
	return
}
