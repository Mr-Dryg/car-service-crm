package domain

import "testing"

func TestChangeStatus(t *testing.T) {
	type testCase struct {
		name         string
		currentState string
		newState     string
		expectError  error
	}

	tests := []testCase{
		{
			name:         "new -> confirmed",
			currentState: StatusNew,
			newState:     StatusConfirmed,
			expectError:  nil,
		},
		{
			name:         "new -> completed",
			currentState: StatusNew,
			newState:     StatusCompleted,
			expectError:  ErrInvalidStatusTransition(StatusNew, StatusCompleted),
		},
		{
			name:         "cancel request -> ready",
			currentState: StatusCancelRequested,
			newState:     StatusReady,
			expectError:  ErrInvalidStatusTransition(StatusCancelRequested, StatusReady),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			order := Order{Status: tc.currentState}
			err := order.ChangeStatus(tc.newState)

			if tc.expectError == nil {
				if err != nil {
					t.Errorf("expect: nil, but err: %v", err)
				} else if tc.newState != order.Status {
					t.Errorf("status didn't change: expect %q, got %q", tc.newState, order.Status)
				}
			} else if err == nil || err.Error() != tc.expectError.Error() {
				t.Errorf("expect: %v, but err: %v", tc.expectError, err)
			}
		})
	}
}
