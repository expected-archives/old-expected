package aws

import "github.com/aws/aws-sdk-go/aws/session"

var sess *session.Session

func Init() error {
	s, err := session.NewSession()
	if err != nil {
		return err
	}
	sess = s
	return nil
}

func Session() *session.Session {
	return sess
}
