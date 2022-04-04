package event

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/notification"
	"github.com/techlabs/swabhav/tsam/repository"
)

// EventPool consists of the pool of events.
var EventPool *Pool

// Pool consists of events (only notification for now).
type Pool struct {
	Events chan Event
	DB     *gorm.DB
	Repo   repository.Repository
}

// InitializePool will initialize the EventPool so it can be used by other modules to push events.
func InitializePool(db *gorm.DB, repo repository.Repository) {
	EventPool = &Pool{
		Events: make(chan Event, 100),
		Repo:   repo,
		DB:     db,
	}
	fmt.Println("EVENT POOL -----------------", EventPool)
	go EventPool.listen()
}

// ?Niranjan
// func (p *Pool) RegisterEvent(listener Listener, names ...Name) error {
// 	for _, name := range names {
// 		if _, ok := p.events[name]; ok {
// 			return fmt.Errorf("The '%s' event is already registered", name)
// 		}

// 		p.events[name] = listener
// 	}

// 	return nil
// }

// FireEvent adds an event to event pool.
func FireEvent(eventName Name, notify notification.Notifier) error {
	if EventPool == nil {
		return errors.NewValidationError("EventPool is not initialized. Call InitializePool")
	}
	event := Event{
		Name:     eventName,
		DB:       EventPool.DB,
		Notifier: notify,
	}
	EventPool.Events <- event
	fmt.Println("Event channel -----------------", EventPool.Events)
	return nil
}

func (p *Pool) listen() {
	for event := range p.Events {
		fmt.Println("In listen in range ------------")
		event.Notifier.Notify(p.DB, p.Repo)
	}
	fmt.Println("*****************************Ending pool listen*************************************")
}
