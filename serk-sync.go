package config

// thread-unsafe
type SerkAdapterSync struct {
	C *Config
}

var _ ISerkDataAPI = (*SerkAdapterSync)(nil)

func (a *SerkAdapterSync) Test() {

}
