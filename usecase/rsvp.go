package usecase

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
	"time"

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

func (ru *rsvpUsecase) WriteRsvpsCsv(ctx context.Context, p *rsvp.Parameter) (*rsvp.File, error) {
	rsvpResult, err := ru.GetRsvps(ctx, p)
	if err != nil {
		return nil, err
	}

	file := new(rsvp.File)
	buffer := &bytes.Buffer{}
	writer := csv.NewWriter(buffer)

	records := [][]string{}

	//Set Header
	records = append(records, []string{"Number", "Name", "Address", "Attend", "Message", "Created Date"})

	for i, item := range rsvpResult.Data {
		var record = []string{}
		record = append(record, strconv.Itoa(i+1))
		record = append(record, item.Name)
		record = append(record, item.Address)
		record = append(record, item.Attend.String())
		record = append(record, item.Message)
		record = append(record, item.CreatedAt.Format("2006-01-02 15-04-05"))

		records = append(records, record)
	}

	if err := writer.WriteAll(records); err != nil {
		return nil, err
	}

	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("rsvp-%s.csv", timestamp)

	file.Content = buffer.Bytes()
	file.Name = filename

	return file, nil
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
