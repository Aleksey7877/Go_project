package ledger

import (
	"errors"
	"testing"
	"time"
)

func TestTransactionValidate(t *testing.T) {
	t.Parallel()

	validDate := time.Date(2026, 9, 10, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		tx      Transaction
		wantErr bool
	}{
		{
			name: "valid transaction",
			tx: Transaction{
				Amount:   100.35,
				Category: "food",
				Date:     validDate,
			},
			wantErr: false,
		},
		{
			name: "zero amount",
			tx: Transaction{
				Amount:   0,
				Category: "food",
				Date:     validDate,
			},
			wantErr: true,
		},
		{
			name: "negative amount",
			tx: Transaction{
				Amount:   -12.7,
				Category: "food",
				Date:     validDate,
			},
			wantErr: true,
		},
		{
			name: "empty category",
			tx: Transaction{
				Amount:   100,
				Category: "",
				Date:     validDate,
			},
			wantErr: true,
		},
		{
			name: "zero date",
			tx: Transaction{
				Amount:   100,
				Category: "food",
				Date: time.Time{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.tx.Validate()

			if tt.wantErr && err == nil {
				t.Errorf("expected error, got nil")
			}

			if !tt.wantErr && err != nil {
				t.Errorf("expected nil, got error: %v", err)
			}
		})
	}

}

func TestBudgetValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		budget  Budget
		wantErr bool
	}{
		{
			name: "valid budget",
			budget: Budget{
				Category: "food",
				Limit:    1000.5,
				Period:   "2026",
			},
			wantErr: false,
		},
		{
			name: "zero limit",
			budget: Budget{
				Category: "food",
				Limit:    0,
				Period:   "2026",
			},
			wantErr: true,
		},
		{
			name: "negative limit",
			budget: Budget{
				Category: "food",
				Limit:    -114.5,
				Period:   "2026",
			},
			wantErr: true,
		},
		{
			name: "empty category",
			budget: Budget{
				Category: "",
				Limit:    1000.5,
				Period:   "2026",
			},
			wantErr: true,
		},
		{
			name: "empty period",
			budget: Budget{
				Category: "food",
				Limit:    1000.5,
				Period:   "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.budget.Validate()

			if tt.wantErr && err == nil {
				t.Errorf("expected error, got nil")
			}

			if !tt.wantErr && err != nil {
				t.Errorf("expected nil, got error: %v", err)
			}
		})
	}
}

func TestAddTransactionBudgetRules(t *testing.T) {
	Reset()
	t.Cleanup(Reset)

	budget := Budget{
		Category: "food",
		Limit:    5000,
		Period:   "2026",
	}

	validDate := time.Date(2026, time.September, 10, 0, 0, 0, 0, time.UTC)

	createdBudget, err := SetBudget(budget)
	if err != nil {
		t.Fatalf("SetBudget() unexpected error: %v", err)
	}
	if createdBudget.Category != budget.Category {
		t.Errorf("expected category %q, got %q", budget.Category, createdBudget.Category)
	}

	_, err = AddTransaction(Transaction{Amount: 100.35, Category: "food", Date: validDate})
	if err != nil {
		t.Fatalf("AddTransaction() unexpected error: %v", err)
	}

	if len(ListTransactions()) != 1 {
		t.Errorf("expected len(ListTransactions()) 1, got %d", len(ListTransactions()))
	}

	_, err = AddTransaction(Transaction{Amount: 5000.1, Category: "food", Date: validDate})
	if !errors.Is(err, ErrBudgetExceeded) {
		t.Errorf("expected ErrBudgetExceeded, got %v", err)
	}

	if len(ListTransactions()) != 1 {
		t.Errorf("expected len(ListTransactions()) 1, got %d", len(ListTransactions()))
	}
}
