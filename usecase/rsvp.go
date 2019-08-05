package usecase

import (
	"context"

	rsvp "github.com/faris-arifiansyah/fws-rsvp"
)

// AccessProvider are collections of provider that used by usecase
type AccessProvider struct {
	RsvpRepo rsvp.RsvpRepo
}

type rsvpUsecase struct {
	*AccessProvider
}

func NewRsvpUsecase(pvd *AccessProvider) rsvp.Usecase {
	return &rsvpUsecase{pvd}
}

func (ru *rsvpUsecase) CreateRsvp(ctx context.Context, rp rsvp.Rsvp) (rsvp.Rsvp, error) {
	return ru.RsvpRepo.CreateRsvp(ctx, rp)
}

func (ru *rsvpUsecase) GetRsvps(ctx context.Context, p *rsvp.Parameter) (*rsvp.RsvpResult, error) {
	return ru.RsvpRepo.GetRsvps(ctx, p)
}
