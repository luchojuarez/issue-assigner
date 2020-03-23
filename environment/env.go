package localEnvironment

import (
	"sync"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
)

var (
	instance *LocalEnvironment
	lock     = &sync.Mutex{}
)

type LocalEnvironment struct {
	restClients map[string]*resty.Client
	IsTestEnv   bool
}

// thread safe
func GetEnv() *LocalEnvironment {
	lock.Lock()
	defer lock.Unlock()

	if instance == nil {
		instance = &LocalEnvironment{
			restClients: make(map[string]*resty.Client),
			IsTestEnv:   true,
		}
	}
	return instance
}

func (this LocalEnvironment) GetResty(name string) *resty.Client {
	if this.restClients[name] == nil {
		newResty := resty.New()

		// TODO don't do this in production
		if this.IsTestEnv {
			httpmock.ActivateNonDefault(newResty.GetClient())
			httpmock.Activate()
		}
		this.restClients[name] = newResty
	}
	return this.restClients[name]
}
