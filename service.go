package rsvp

import (
	"context"
	"time"

	"github.com/globalsign/mgo/bson"
)

// Parameter is a struct to simplify passing parameter into function
type Parameter struct {
	Limit  int
	Offset int
}

// RsvpResult is a struct container to put result
type RsvpResult struct {
	Data  []*Rsvp
	Total int64
}

// Rsvp Entity
type Rsvp struct {
	ID        bson.ObjectId `json:"-" bson:"_id,omitempty"`
	Name      string        `json:"name" bson:"name"`
	Address   string        `json:"address" bson:"address"`
	Attend    uint8         `json:"attend" bson:"attend"`
	Message   string        `json:"message" bson:"message"`
	CreatedAt time.Time     `json:"created_at" bson:"created_at"`
}

// RsvpRepo provides data interchange between
// application and data provider.
type RsvpRepo interface {
	CreateRsvp(ctx context.Context, rp *Rsvp) error
	GetRsvps(ctx context.Context, p *Parameter) (*RsvpResult, error)
}

type Usecase interface {
	CreateRsvp(ctx context.Context, rp *Rsvp) error
	GetRsvps(ctx context.Context, p *Parameter) (*RsvpResult, error)
}
