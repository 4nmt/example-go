// +build integration

package lendbook

import (
	"context"
	"reflect"
	"testing"

	"github.com/jinzhu/gorm"
	testutil "github.com/neverdiefc/example-go/config/database/pg/util"
	"github.com/neverdiefc/example-go/domain"
)

func TestPGService_Create(t *testing.T) {
	t.Parallel()
	testDB, _, cleanup := testutil.CreateTestDatabase(t)
	defer cleanup()
	err := testutil.MigrateTables(testDB)
	if err != nil {
		t.Fatalf("Failed to migrate table by error %v", err)
	}

	user := domain.User{}
	err = testDB.Create(&user).Error
	if err != nil {
		t.Fatalf("Failed to create user by error %v", err)
	}

	book := domain.Book{}
	err = testDB.Create(&book).Error
	if err != nil {
		t.Fatalf("Failed to create book by error %v", err)
	}

	fakeLendbookID := domain.MustGetUUIDFromString("1698bbd6-e0c8-4957-a5a9-8c536970994b")

	type args struct {
		p *domain.Lendbook
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "create successfully",
			args: args{
				&domain.Lendbook{
					BookID: book.ID,
					UserID: user.ID,
				},
			},
		},
		{
			name: "failed update (invalid bookID)",
			args: args{
				&domain.Lendbook{
					BookID: fakeLendbookID,
					UserID: user.ID,
				},
			},
			wantErr: true,
		},
		{
			name: "failed update (invalid lendbookID)",
			args: args{
				&domain.Lendbook{
					BookID: book.ID,
					UserID: fakeLendbookID,
				},
			},
			wantErr: true,
		},
		{
			name: "failed update ",
			args: args{
				&domain.Lendbook{
					BookID: book.ID,
					UserID: user.ID,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &pgService{
				db: testDB,
			}
			if err := s.Create(context.Background(), tt.args.p); (err != nil) != tt.wantErr {
				t.Errorf("pgService.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPGService_Update(t *testing.T) {
	t.Parallel()
	testDB, _, cleanup := testutil.CreateTestDatabase(t)
	defer cleanup()
	err := testutil.MigrateTables(testDB)
	if err != nil {
		t.Fatalf("Failed to migrate table by error %v", err)
	}

	user := domain.User{}
	err = testDB.Create(&user).Error
	if err != nil {
		t.Fatalf("Failed to create user by error %v", err)
	}

	book := domain.Book{}
	err = testDB.Create(&book).Error
	if err != nil {
		t.Fatalf("Failed to create book by error %v", err)
	}

	mathBook := domain.Book{}
	err = testDB.Create(&mathBook).Error
	if err != nil {
		t.Fatalf("Failed to create book by error %v", err)
	}

	lendbook := domain.Lendbook{
		BookID: book.ID,
		UserID: user.ID,
	}
	err = testDB.Create(&lendbook).Error
	if err != nil {
		t.Fatalf("Failed to create lendbook by error %v", err)
	}

	fakeLendbookID := domain.MustGetUUIDFromString("1698bbd6-e0c8-4957-a5a9-8c536970994b")

	type args struct {
		p *domain.Lendbook
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "success update",
			args: args{
				&domain.Lendbook{
					Model:  domain.Model{ID: lendbook.ID},
					BookID: mathBook.ID,
					UserID: user.ID,
				},
			},
		},
		{
			name: "failed update ",
			args: args{
				&domain.Lendbook{
					Model:  domain.Model{ID: fakeLendbookID},
					UserID: user.ID,
				},
			},
			wantErr: ErrNotFound,
		},
		{
			name: "failed update (invalid bookID)",
			args: args{
				&domain.Lendbook{
					Model:  domain.Model{ID: lendbook.ID},
					BookID: fakeLendbookID,
					UserID: user.ID,
				},
			},
			wantErr: ErrRecordBookNotFound,
		},
		{
			name: "failed update (invalid lendbookID)",
			args: args{
				&domain.Lendbook{
					Model:  domain.Model{ID: lendbook.ID},
					BookID: book.ID,
					UserID: fakeLendbookID,
				},
			},
			wantErr: ErrRecordUserNotFound,
		},
		{
			name: "failed update ",
			args: args{
				&domain.Lendbook{
					Model:  domain.Model{ID: lendbook.ID},
					BookID: mathBook.ID,
					UserID: user.ID,
				},
			},
			wantErr: ErrBookIsBusy,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &pgService{
				db: testDB,
			}
			_, err := s.Update(context.Background(), tt.args.p)
			if err != nil && err != tt.wantErr {
				t.Errorf("pgService.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && tt.wantErr != nil {
				t.Errorf("pgService.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPGService_Find(t *testing.T) {
	t.Parallel()
	testDB, _, cleanup := testutil.CreateTestDatabase(t)
	defer cleanup()
	err := testutil.MigrateTables(testDB)
	if err != nil {
		t.Fatalf("Failed to migrate table by error %v", err)
	}

	lendbook := domain.Lendbook{}
	err = testDB.Create(&lendbook).Error
	if err != nil {
		t.Fatalf("Failed to create lendbook by error %v", err)
	}

	fakeLendbookID := domain.MustGetUUIDFromString("1698bbd6-e0c8-4957-a5a9-8c536970994b")

	type args struct {
		p *domain.Lendbook
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.Lendbook
		wantErr error
	}{
		{
			name: "success find correct lendbook",
			args: args{
				&domain.Lendbook{
					Model: domain.Model{ID: lendbook.ID},
				},
			},
			want: &domain.Lendbook{
				Model: domain.Model{ID: lendbook.ID},
			},
		},
		{
			name: "failed find lendbook by not exist lendbook id",
			args: args{
				&domain.Lendbook{
					Model: domain.Model{ID: fakeLendbookID},
				},
			},
			wantErr: ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &pgService{
				db: testDB,
			}

			got, err := s.Find(context.Background(), tt.args.p)
			if err != nil && err != tt.wantErr {
				t.Errorf("pgService.Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && tt.wantErr != nil {
				t.Errorf("pgService.Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != nil && got.ID.String() != tt.want.ID.String() {
				t.Errorf("pgService.Find() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPGService_FindAll(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		in0 context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []domain.Lendbook
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &pgService{
				db: tt.fields.db,
			}
			got, err := s.FindAll(tt.args.in0)
			if (err != nil) != tt.wantErr {
				t.Errorf("pgService.FindAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("pgService.FindAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPGService_Delete(t *testing.T) {
	t.Parallel()
	testDB, _, cleanup := testutil.CreateTestDatabase(t)
	defer cleanup()
	err := testutil.MigrateTables(testDB)
	if err != nil {
		t.Fatalf("Failed to migrate table by error %v", err)
	}

	lendbook := domain.Lendbook{}
	err = testDB.Create(&lendbook).Error
	if err != nil {
		t.Fatalf("Failed to create lendbook by error %v", err)
	}

	fakeLendbookID := domain.MustGetUUIDFromString("1698bbd6-e0c8-4957-a5a9-8c536970994b")

	type args struct {
		p *domain.Lendbook
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "success delete",
			args: args{
				&domain.Lendbook{
					Model: domain.Model{ID: lendbook.ID},
				},
			},
		},
		{
			name: "failed delete by not exist lendbook id",
			args: args{
				&domain.Lendbook{
					Model: domain.Model{ID: fakeLendbookID},
				},
			},
			wantErr: ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &pgService{
				db: testDB,
			}
			err := s.Delete(context.Background(), tt.args.p)
			if err != nil && err != tt.wantErr {
				t.Errorf("pgService.Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && tt.wantErr != nil {
				t.Errorf("pgService.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewPGService(t *testing.T) {
	type args struct {
		db *gorm.DB
	}
	tests := []struct {
		name string
		args args
		want Service
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPGService(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPGService() = %v, want %v", got, tt.want)
			}
		})
	}
}
