package localEnvironment

import (
	"os"
	"strings"
	"sync"

	"github.com/luchojuarez/issue-assigner/models"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
)

var (
	instance *LocalEnvironment
	lock     = &sync.Mutex{}
)

type LocalEnvironment struct {
	restClients       map[string]*resty.Client
	UserStorage       map[string]*models.User
	PrStorage         map[string][]*models.PR
	EventTraceStorage []*models.Event
	IsTestEnv         bool
}

// thread safe
func GetEnv() *LocalEnvironment {
	lock.Lock()
	defer lock.Unlock()

	if instance == nil {
		instance = &LocalEnvironment{
			restClients:       make(map[string]*resty.Client),
			UserStorage:       make(map[string]*models.User),
			PrStorage:         make(map[string][]*models.PR),
			EventTraceStorage: make([]*models.Event, 0),
			IsTestEnv:         false,
		}
	}
	return instance
}

func (this *LocalEnvironment) GetResty(name string) *resty.Client {
	if this.restClients[name] == nil {
		newResty := resty.New()

		if IsTestEnv() {
			httpmock.ActivateNonDefault(newResty.GetClient())
			httpmock.Activate()
		}
		this.restClients[name] = newResty
	}
	return this.restClients[name]
}

func (this *LocalEnvironment) GetUserStorage() *map[string]*models.User {
	return &this.UserStorage
}

func (this *LocalEnvironment) ClearUserStorage() {
	this.UserStorage = make(map[string]*models.User)
}

func (this *LocalEnvironment) GetPrStorage() *map[string][]*models.PR {
	return &this.PrStorage
}

func (this *LocalEnvironment) ClearPrStorage() {
	this.PrStorage = make(map[string][]*models.PR)
}

//thread safe
func (this *LocalEnvironment) AddEventSafe(elem *models.Event) {
	lock.Lock()
	defer lock.Unlock()
	this.EventTraceStorage = append(this.EventTraceStorage, elem)
}

func (this *LocalEnvironment) AddEvent(elem *models.Event) {
	this.EventTraceStorage = append(this.EventTraceStorage, elem)
}

func (this *LocalEnvironment) GetAllEvents() *[]*models.Event {
	return &this.EventTraceStorage
}

func (this *LocalEnvironment) CleanAll() {
	instance = nil
}

func (this *LocalEnvironment) ClearEventTracer() {
	lock.Lock()
	defer lock.Unlock()
	this.EventTraceStorage = make([]*models.Event, 0)
}

func IsTestEnv() bool {
	return strings.Contains(os.Args[0], ".test")
}
