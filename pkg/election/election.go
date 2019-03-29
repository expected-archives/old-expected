package election

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
	"time"
)

const TTL = time.Second * 10

// Election represent a way to elect 1 master over multiple nodes.
type Election struct {
	ServiceName string

	session *concurrency.Session
	client  *clientv3.Client
	key     string
}

// NewElection create an instance of Election.
func NewElection(serviceName, etcdAddress string) (*Election, error) {
	client, err := clientv3.New(clientv3.Config{Endpoints: []string{etcdAddress}})
	if err != nil {
		return nil, err
	}
	return &Election{
		ServiceName: serviceName,
		session:     nil,
		client:      client,
		key:         fmt.Sprintf("service/%s/leader", serviceName),
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
				return isLeader
			}
			time.Sleep(1 * time.Second)
		}
	}
}

// newSession creates a newSession in consul with a TTL and behavior set to delete
func (e *Election) newSession() error {
	session, err := concurrency.NewSession(e.client, concurrency.WithTTL(int(TTL.Seconds())))
	if err != nil {
		return err
	}
	e.session = session

	return nil
}

// acquireSession try to acquire a session.
// Return a bool which is the representation if it is the leader.
func (e *Election) acquireSession() (bool, error) {
	elect := concurrency.NewElection(e.session, e.key)
	if err := elect.Campaign(context.Background(), "e1"); err != nil {
		return false, err
	}
	return true, nil
}

// Close gracefully destroy the session.
func (e *Election) Close() error {
	return e.session.Close()
}
