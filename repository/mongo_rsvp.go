package repository

import (
	"context"

	rsvp "github.com/faris-arifiansyah/fws-rsvp"
	"github.com/faris-arifiansyah/mgoi"
)

type mongoRsvp struct {
	db mgoi.DatabaseManager
}

func NewMongoRsvp(db mgoi.DatabaseManager) rsvp.RsvpRepo {
	return &mongoRsvp{db}
}

func (mr *mongoRsvp) CreateRsvp(ctx context.Context, rp *rsvp.Rsvp) error {
	return nil
}

func (mr *mongoRsvp) GetRsvps(ctx context.Context, p *rsvp.Parameter) (*rsvp.RsvpResult, error) {
	return nil, nil
}
