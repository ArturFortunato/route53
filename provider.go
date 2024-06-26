package route53

import (
	"context"
	"fmt"
	"time"

	r53 "github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/libdns/libdns"
)

// Provider implements the libdns interfaces for Route53
type Provider struct {
	MaxRetries         int           `json:"max_retries,omitempty"`
	MaxWaitDur         time.Duration `json:"max_wait_dur,omitempty"`
	WaitForPropagation bool          `json:"wait_for_propagation,omitempty"`
	Region             string        `json:"region,omitempty"`
	AWSProfile         string        `json:"aws_profile,omitempty"`
	AccessKeyId        string        `json:"access_key_id,omitempty"`
	SecretAccessKey    string        `json:"secret_access_key,omitempty"`
	Token              string        `json:"token,omitempty"`
	client             *r53.Client
}

// GetRecords lists all the records in the zone.
func (p *Provider) GetRecords(ctx context.Context, zone string) ([]libdns.Record, error) {
	p.init(ctx)

	zoneID, err := p.getZoneID(ctx, zone)
	if err != nil {
		return nil, err
	}

	records, err := p.getRecords(ctx, zoneID, zone)
	if err != nil {
		return nil, err
	}

	return records, nil
}

// AppendRecords adds records to the zone. It returns the records that were added.
func (p *Provider) AppendRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	p.init(ctx)
	fmt.Println("Appending records for zone: ", zone)

	zoneID, err := p.getZoneID(ctx, zone)
	if err != nil {
		fmt.Println("Failed to get zone ID", err)
		return nil, err
	}

	var createdRecords []libdns.Record

	for _, record := range records {
		fmt.Println("Creating record: ", record.ID, record.Name, record.Value)
		newRecord, err := p.createRecord(ctx, zoneID, record, zone)
		if err != nil {
			fmt.Println("Failed to create record", err)
			return nil, err
		}
		createdRecords = append(createdRecords, newRecord)
	}

	return createdRecords, nil
}

// DeleteRecords deletes the records from the zone. If a record does not have an ID,
// it will be looked up. It returns the records that were deleted.
func (p *Provider) DeleteRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	p.init(ctx)

	zoneID, err := p.getZoneID(ctx, zone)
	if err != nil {
		return nil, err
	}

	var deletedRecords []libdns.Record

	for _, record := range records {
		deletedRecord, err := p.deleteRecord(ctx, zoneID, record, zone)
		if err != nil {
			return nil, err
		}
		deletedRecords = append(deletedRecords, deletedRecord)
	}

	return deletedRecords, nil
}

// SetRecords sets the records in the zone, either by updating existing records
// or creating new ones. It returns the updated records.
func (p *Provider) SetRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	p.init(ctx)

	zoneID, err := p.getZoneID(ctx, zone)
	if err != nil {
		return nil, err
	}

	var updatedRecords []libdns.Record

	for _, record := range records {
		updatedRecord, err := p.updateRecord(ctx, zoneID, record, zone)
		if err != nil {
			return nil, err
		}
		updatedRecords = append(updatedRecords, updatedRecord)
	}

	return updatedRecords, nil
}

// Interface guards
var (
	_ libdns.RecordGetter   = (*Provider)(nil)
	_ libdns.RecordAppender = (*Provider)(nil)
	_ libdns.RecordSetter   = (*Provider)(nil)
	_ libdns.RecordDeleter  = (*Provider)(nil)
)
