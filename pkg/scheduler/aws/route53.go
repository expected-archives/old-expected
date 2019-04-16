package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
)

func Route53AddRecord(hostedZoneId, recordType, recordName string, values []string) error {
	var records []*route53.ResourceRecord

	for _, value := range values {
		records = append(records, &route53.ResourceRecord{
			Value: aws.String(value),
		})
	}
	_, err := route53client.ChangeResourceRecordSets(&route53.ChangeResourceRecordSetsInput{
		HostedZoneId: aws.String(hostedZoneId),
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				{
					Action: aws.String("CREATE"),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Type:            aws.String(recordType),
						Name:            aws.String(recordName),
						TTL:             aws.Int64(60),
						ResourceRecords: records,
					},
				},
			},
		},
	})
	return err
}

func Route53DelRecord(hostedZoneId, recordType, recordName string) error {
	_, err := route53client.ChangeResourceRecordSets(&route53.ChangeResourceRecordSetsInput{
		HostedZoneId: aws.String(hostedZoneId),
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				{
					Action: aws.String("DELETE"),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Type: aws.String(recordType),
						Name: aws.String(recordName),
					},
				},
			},
		},
	})
	return err
}
