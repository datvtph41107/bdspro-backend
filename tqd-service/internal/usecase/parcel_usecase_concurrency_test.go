package usecase

import (
	"context"
	"testing"
	"time"

	"tqd/internal/dto"
	"tqd/internal/interface/repo"
)

type parcelForkJoinProbe struct {
	repo.IParcelRepo
	analysisStarted chan struct{}
	countStarted    chan struct{}
	releaseAnalysis chan struct{}
	releaseCount    chan struct{}
}

func (p *parcelForkJoinProbe) SummarizePolygonAnalysis(
	context.Context,
	[]repo.Point,
	bool,
	uint32,
	*dto.PolygonTargetFilter,
	bool,
) (*dto.PolygonAnalysisSummary, []dto.ParcelLayerInfo, error) {
	close(p.analysisStarted)
	<-p.releaseAnalysis
	return &dto.PolygonAnalysisSummary{}, nil, nil
}

func (p *parcelForkJoinProbe) CountParcelsByPolygon(
	context.Context,
	[]repo.Point,
	bool,
	*dto.PolygonTargetFilter,
) (int64, error) {
	close(p.countStarted)
	<-p.releaseCount
	return 7, nil
}

func TestFindMapTargetsByPolygonJoinsAllAnalysisBranches(t *testing.T) {
	probe := &parcelForkJoinProbe{
		analysisStarted: make(chan struct{}),
		countStarted:    make(chan struct{}),
		releaseAnalysis: make(chan struct{}),
		releaseCount:    make(chan struct{}),
	}
	u := &parcelUsecase{parcelRepo: probe}
	zoom := uint32(15)
	points := []repo.Point{
		{Lat: 0, Lng: 0},
		{Lat: 0, Lng: 0.001},
		{Lat: 0.001, Lng: 0},
	}

	done := make(chan error, 1)
	go func() {
		_, err := u.FindMapTargetsByPolygon(
			context.Background(),
			points,
			true,
			1,
			1,
			20,
			&zoom,
			nil,
			true,
		)
		done <- err
	}()

	select {
	case <-probe.analysisStarted:
	case <-time.After(time.Second):
		t.Fatal("analysis branch did not start")
	}
	select {
	case <-probe.countStarted:
	case <-time.After(time.Second):
		t.Fatal("count branch did not start")
	}

	close(probe.releaseAnalysis)
	select {
	case err := <-done:
		t.Fatalf("method returned before count branch joined: %v", err)
	case <-time.After(20 * time.Millisecond):
		// Expected: the second request-scoped branch is still owned by this call.
	}

	close(probe.releaseCount)
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("FindMapTargetsByPolygon() error = %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("method did not return after both analysis branches completed")
	}
}
