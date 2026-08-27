package service

import (
	// "context"
	// "errors"
	"testing"
	"time"

	"github.com/Mr-Dryg/car-service-crm/internal/domain"
)

func TestCheckTimeAvailability(t *testing.T) {
	type testCase struct {
		name       string
		dateStr    string
		parsedDate time.Time
		err        error
	}
	tests := []testCase{
		{
			name:       "invalid format",
			dateStr:    "2006-01-02",
			parsedDate: time.Time{},
			err:        ErrInvalidDateFormat,
		},
		{
			name:       "date in past",
			dateStr:    "12-05-2006",
			parsedDate: time.Date(2006, 5, 12, 0, 0, 0, 0, time.UTC),
			err:        ErrDateInPast,
		},
		{
			name:       "good date",
			dateStr:    "31-12-2099",
			parsedDate: time.Date(2099, 12, 31, 0, 0, 0, 0, time.UTC),
			err:        nil,
		},
		{
			name:       "today test",
			dateStr:    time.Now().Truncate(24 * time.Hour).Format(domain.DateLayout),
			parsedDate: time.Now().Truncate(24 * time.Hour),
			err:        nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			get, err := checkTimeAvailability(0, test.dateStr)
			if err != test.err {
				t.Fatalf("excpect err: %v, get: %v", test.err, err)
			} else if err == nil && !get.Equal(test.parsedDate) {
				t.Fatalf("expect time: %v, get: %v\n", test.parsedDate, get)
			}
		})
	}
}

// type fakeOrderRepo struct {
// 	created *domain.Order
// }

// func (f *fakeOrderRepo) Create(ctx context.Context, order *domain.Order, parsedDate time.Time) error {
// 	f.created = order
// 	return nil
// }

// func (f *fakeOrderRepo) GetByOrderID(ctx context.Context, orderID int64) (*domain.Order, error) {
// 	if f.created.ID != orderID {
// 		return nil, errors.New("missing orderID")
// 	}
// 	return f.created, nil
// }

// func (f *fakeOrderRepo) GetByBranchID(ctx context.Context, branchID int64) ([]domain.Order, error) {
// 	return nil, nil
// }

// func (f *fakeOrderRepo) GetByCarID(ctx context.Context, carID int64) ([]domain.Order, error) {
// 	return nil, nil
// }

// func (f *fakeOrderRepo) UpdateStatus(ctx context.Context, orderID int64, status string) error {
// 	f.created.Status = status
// 	return nil
// }

// func (f *fakeOrderRepo) UpdateShedule(ctx context.Context, orderID int64, prefDate time.Time, prefTime string) error {
// 	f.created.PreferredDate = prefDate.Format(domain.DateLayout)
// 	f.created.PreferredTime = prefTime
// 	return nil
// }

// func (f *fakeOrderRepo) UpdateCost(ctx context.Context, orderID int64, cost float64) error {
// 	f.created.Cost = cost
// 	return nil
// }

// func (f *fakeOrderRepo) UpdateClientConfirmed(ctx context.Context, orderID int64, flag bool) error {
// 	f.created.ClientConfirmed = flag
// 	return nil
// }

// func (f *fakeOrderRepo) UpdateNotes(ctx context.Context, orderID int64, notes string) error {
// 	f.created.Notes = notes
// 	return nil
// }

// func TestCreate(t * testing.T)  {
// 	repo := &fakeOrderRepo{}
// 	svc := NewOrderService(repo)
// }
