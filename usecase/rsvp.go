package usecase

import (
	"context"
	"strings"

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
	p.Sort = ru.GetValidSortField(p.Sort)

	return ru.RsvpRepo.GetRsvps(ctx, p)
}

func (ru *rsvpUsecase) GetValidSortField(sf string) string {
	sortFields := map[string]struct{}{
		"created_at":  {},
		"-created_at": {},
		"name":        {},
		"-name":       {},
	}

	if _, valid := sortFields[strings.ToLower(sf)]; valid {
		return sf
	}

	return "-created_at"
}
