package config

// thread-safe
type SerkAdapterAsync struct {
	CS *Source
}

var _ ISerkDataAPI = (*SerkAdapterAsync)(nil)

func (a *SerkAdapterAsync) Test() {

}
