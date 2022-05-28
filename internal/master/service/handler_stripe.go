package service

import (
	"context"
	"errors"
	"net/http"

	"github.com/ocean-rw/ocean/internal/lib/httputil"
	"github.com/ocean-rw/ocean/pkg/proto"
)

func (s *Service) AllocStripes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	stripes, err := s.allocStripes(ctx, proto.DefaultMode, 100)
	if err != nil {
		s.logger.Errorf("failed to alloc stripes, err: %s", err)
		httputil.ReplyErr(w, http.StatusInternalServerError, err)
		return
	}
	err = httputil.ReplyJSON(w, http.StatusOK, stripes)
	if err != nil {
		s.logger.Errorf("failed to reply data, err: %s", err)
	}
}

func (s *Service) allocStripes(ctx context.Context, mode proto.Mode, count int) ([]*proto.Stripe, error) {
	sids, err := s.db.IDTable.AllocStripeID(ctx, count)
	if err != nil {
		return nil, err
	}
	stripes := make([]*proto.Stripe, 0, len(sids))
	for _, id := range sids {
		disks, err := s.db.DiskTable.AllocDisks(ctx, mode, nil)
		if err != nil {
			s.logger.Errorf("failed to alloc disks, err: %s", err)
			continue
		}
		stripes = append(stripes, &proto.Stripe{
			ID:    id,
			Mode:  mode,
			Disks: disks,
		})
	}

	if len(stripes) == 0 {
		return nil, errors.New("empty stripe")
	}

	err = s.db.StripeTable.Insert(ctx, stripes)
	if err != nil {
		return nil, err
	}
	return stripes, nil
}
