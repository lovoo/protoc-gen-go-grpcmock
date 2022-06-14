package routeguide

import (
	"context"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	DresdenLatitude  = 51.050409
	DresdenLongitude = 13.737262

	GermanyLatitudeMin  = 47.2701114
	GermanyLongitudeMin = 5.8663153
	GermanyLatitudeMax  = 55.099161
	GermanyLongitudeMax = 15.0419319
)

var (
	DresdenCenter = &Point{Latitude: e7(DresdenLatitude), Longitude: e7(DresdenLongitude)}

	DresdenNote = &RouteNote{Location: DresdenCenter, Message: "Dresden"}

	GermanyBoundingBox = &Rectangle{
		Lo: &Point{Latitude: e7(GermanyLatitudeMin), Longitude: e7(GermanyLongitudeMin)},
		Hi: &Point{Latitude: e7(GermanyLatitudeMin), Longitude: e7(GermanyLongitudeMin)},
	}
)

func TestGetFeature(t *testing.T) {
	// Create a new mock client for the RouteGuide service.
	m := NewMockRouteGuideClient()
	defer m.AssertExpectations(t)

	// Create the request and response.
	ctx := context.Background()
	req := DresdenCenter
	res := &Feature{Name: "Dresden", Location: req}

	// Set up the expectation.
	m.OnGetFeature(ctx, req).Return(res, nil)

	// Call the client.
	r, err := m.GetFeature(ctx, req)

	// Check that the response is as expected.
	assert.NoError(t, err)
	assert.Equal(t, res, r)
}

func TestListFeatures(t *testing.T) {
	// Create a new mock client for the RouteGuide service.
	m := NewMockRouteGuideClient()
	defer m.AssertExpectations(t)

	// Create the request and response.
	ctx := context.Background()
	req := GermanyBoundingBox
	res := NewMockRouteGuide_ListFeaturesClient()
	defer res.AssertExpectations(t)
	feat := &Feature{Name: "Dresden", Location: DresdenCenter}

	// Set up the expectations.
	res.OnRecv().Return(feat, nil)
	m.OnListFeatures(ctx, req).Return(res, nil)

	// Call the client.
	r, err := m.ListFeatures(ctx, req)

	// Check that the response is as expected.
	assert.NoError(t, err)
	assert.Equal(t, res, r)

	// Use the client streaming handler.
	f, err := r.Recv()

	// Check that the streamed response is as expected.
	assert.NoError(t, err)
	assert.NotNil(t, f)
	assert.Equal(t, feat, f)
}

func TestRecordRoute(t *testing.T) {
	// Create a new mock client for the RouteGuide service.
	m := NewMockRouteGuideClient()
	defer m.AssertExpectations(t)

	// Create the request and response.
	ctx := context.Background()
	res := NewMockRouteGuide_RecordRouteClient()
	defer res.AssertExpectations(t)
	routs := &RouteSummary{PointCount: 1, FeatureCount: 1, Distance: 1, ElapsedTime: 1}

	// Set up the expectations.
	res.OnSend(AnyPoint()).Return(nil)
	res.OnCloseAndRecv().Return(routs, nil)
	m.OnRecordRoute(ctx).Return(res, nil)

	// Call the client.
	r, err := m.RecordRoute(ctx)

	// Check that the response is as expected.
	assert.NoError(t, err)
	assert.Equal(t, res, r)

	// Use the client streaming handler.
	err = r.Send(DresdenCenter)

	// Check that the response is as expected.
	assert.NoError(t, err)

	// Use the client streaming handler.
	rs, err := r.CloseAndRecv()

	// Check that the response is as expected.
	assert.NoError(t, err)
	assert.NotNil(t, rs)
	assert.Equal(t, routs, rs)
}

func TestRouteChat(t *testing.T) {
	m := NewMockRouteGuideClient()
	defer m.AssertExpectations(t)

	// Create the request and response.
	ctx := context.Background()
	res := NewMockRouteGuide_RouteChatClient()
	defer res.AssertExpectations(t)
	routn := &RouteNote{Location: DresdenCenter, Message: "Dresden, Germany"}

	// Set up the expectations.
	res.OnSend(AnyRouteNote()).Return(nil)
	res.OnRecv().Return(routn, nil)
	m.OnRouteChat(ctx).Return(res, nil)

	// Call the client.
	r, err := m.RouteChat(ctx)

	// Check that the response is as expected.
	assert.NoError(t, err)
	assert.Equal(t, res, r)

	// Use the client streaming handler.
	err = r.Send(DresdenNote)

	// Check that the response is as expected.
	assert.NoError(t, err)

	// Use the client streaming handler.
	rn, err := r.Recv()

	// Check that the streamed response is as expected.
	assert.NoError(t, err)
	assert.NotNil(t, rn)
	assert.Equal(t, routn, rn)
}

// e7 converts a coordinate given in degrees into the E7 representation
// (degrees multiplied by 10**7 and rounded to the nearest integer).
func e7(coord float64) int32 {
	return int32(coord * math.Pow10(7))
}
