package config

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type InitContext struct {
	CurrentFileName      string
	SourceHashesRequired []string
	SourceHashesActual   []string
	Data                 []byte
	Logger               *zerolog.Logger
	ErrPtr               *error
	OkPtr                *bool
}

func (ic *InitContext) FromFile(fileName string, hashes ...string) *InitContext {
	ic.CurrentFileName = fileName
	if len(hashes) > 0 {
		ic.SourceHashesRequired = hashes
	}
	return ic
}

func (ic *InitContext) FromBytes(data []byte, hashes ...string) *InitContext {
	ic.Data = data
	if len(hashes) > 0 {
		ic.SourceHashesRequired = hashes
	}
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
	var hash string

	func() {
		hasher := sha256.New()
		switch {
		case len(ic.Data) > 0:
			hasher.Write(ic.Data)
			hash = base64.URLEncoding.EncodeToString(hasher.Sum(nil))
			if !ic.sourceHashIsCorrect(hash) {
				err = fmt.Errorf("incorrect integrity hash '%v' at data", hash)
				return
			}

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

		case len(ic.CurrentFileName) > 0:
			var f *os.File
			f, err = os.Open(ic.CurrentFileName)
			if err != nil {
				return
			}
			defer f.Close()
			if _, err = io.Copy(hasher, f); err != nil {
				return
			}
			hash = base64.URLEncoding.EncodeToString(hasher.Sum(nil))
			if !ic.sourceHashIsCorrect(hash) {
				err = fmt.Errorf("incorrect integrity hash '%v' at file '%v'", hash, ic.CurrentFileName)
				return
			}

			switch {
			case reSuffixYaml.MatchString(ic.CurrentFileName) == true:
				c, err = parseYamlFile(ic.CurrentFileName)
				return
			case reSuffixJson.MatchString(ic.CurrentFileName) == true:
				c, err = parseJsonFile(ic.CurrentFileName)
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

	ic.SourceHashesActual = append(ic.SourceHashesActual, hash)

	// this does inherit these..
	c.ErrPtr = ic.ErrPtr
	c.OkPtr = ic.OkPtr
	c.InitContext = ic

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
		ic.CurrentFileName = currConfigFileName
		conf := ic.Err(&err).Load()
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
	result = readParent(filepath.Dir(ic.CurrentFileName), ic.CurrentFileName)
	result.Set([]string{"parent"}, nil)
	result.Set([]string{"parents"}, nil)

	parentsInherited := make([]string, 0, len(filesAlreadyRead))
	for k := range filesAlreadyRead {
		parentsInherited = append(parentsInherited, k)
	}
	result.Set([]string{"parents-inherited"}, parentsInherited)

	ic.Logger.Info().Msg("K2aUDgz: reading the config file(s) OK")
	return
}

func (ic *InitContext) sourceHashIsCorrect(hash string) (ok bool) {
	if len(ic.SourceHashesRequired) == 0 {
		return true
	}
	i1 := len(ic.SourceHashesActual)
	if i1 >= len(ic.SourceHashesRequired) {
		return false
	}
	if ic.SourceHashesRequired[i1] != hash {
		return false
	}
	return true
}
