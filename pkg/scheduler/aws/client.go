package aws

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

var (
	sess          *session.Session
	route53client *route53.Route53
)

func Init() error {
	s, err := session.NewSession()
	if err != nil {
		return err
	}
	sess = s
	route53client = route53.New(sess)
	return nil
}

func Session() *session.Session {
	return sess
}
