package election

import (
	"fmt"
	"github.com/expectedsh/expected/pkg/services/consul"
	"github.com/hashicorp/consul/api"
	"time"
)

const TTL = "10s"

// Election represent a way to elect 1 master over multiple nodes.
type Election struct {
	CheckInterval time.Duration
	ServiceName   string

	sessionID string
	client    *api.Client
	key       string
	doneChan  chan struct{}
}

// NewElection create an instance of Election.
func NewElection(service consul.Service) (*Election, error) {
	return &Election{
		CheckInterval: time.Second * 2,
		ServiceName:   service.Config().AppName,
		sessionID:     "",
		client:        service.Client(),
		key:           fmt.Sprintf("service/%s/leader", service.Config().AppName),
		doneChan:      make(chan struct{}),
	}, err
}

// ElectLeader check with consul if the current session (with the serviceName)
// has a leader, if not, the leader is the runner of this function.
// If lock is set to true, when the leader will crash another will pick
// his role.
func (e *Election) ElectLeader(lock bool) bool {
	err := e.newSession()
	if err != nil {
		return false
	}
	isLeader, err := e.acquireSession()
	if isLeader || !lock {
		go e.renewSession()
		return isLeader
	} else {
		for {
			isLeader, err := e.acquireSession()
			if err != nil {
				err := e.newSession()
				if err != nil {
					return false
				}
				continue
			}
			if isLeader {
				go e.renewSession()
				return isLeader
			}
			time.Sleep(e.CheckInterval)
		}
	}
}

// newSession creates a newSession in consul with a TTL and behavior set to delete
func (e *Election) newSession() error {
	sessionEntry := &api.SessionEntry{
		TTL:      TTL,
		Behavior: api.SessionBehaviorDelete,
	}
	sessionID, _, err := e.client.Session().Create(sessionEntry, nil)
	if err != nil {
		return err
	}
	e.sessionID = sessionID

	return nil
}

func (e *Election) renewSession() error {
	return e.client.Session().RenewPeriodic(TTL, e.sessionID, nil, e.doneChan)
}

// acquireSession try to acquire a session.
// Return a bool which is the representation if it is the leader.
func (e *Election) acquireSession() (bool, error) {
	pair := &api.KVPair{
		Key:     e.key,
		Value:   []byte(e.sessionID),
		Session: e.sessionID,
	}

	acquired, _, err := e.client.KV().Acquire(pair, nil)

	return acquired, err
}

// Close gracefully destroy the session.
func (e *Election) Close() error {
	close(e.doneChan)
	_, err := e.client.Session().Destroy(e.sessionID, nil)
	if err != nil {
		return err
	}

	return nil
}
