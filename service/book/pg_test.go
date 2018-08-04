// +build integration

package book

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

	category := domain.Category{}
	errCategory := testDB.Create(&category).Error
	if errCategory != nil {
		t.Fatalf("Failed to create book by error %v", errCategory)
	}

	fakeBookID := domain.MustGetUUIDFromString("1698bbd6-e0c8-4957-a5a9-8c536970994b")

	type args struct {
		p *domain.Book
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				&domain.Book{
					Name:        "Create New Book 1",
					Description: "example@gmail.com",
					CategoryID:  category.ID,
				},
			},
		},
		{
			name: "Failed",
			args: args{
				&domain.Book{
					Name:        "Create New Book 1",
					Description: "example@gmail.com",
					CategoryID:  fakeBookID,
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

	englishCategory := domain.Category{}
	errEngCategory := testDB.Create(&englishCategory).Error
	if errEngCategory != nil {
		t.Fatalf("Failed to create category by error %v", errEngCategory)
	}

	mathCategory := domain.Category{}
	errMathCategory := testDB.Create(&mathCategory).Error
	if errMathCategory != nil {
		t.Fatalf("Failed to create category by error %v", errMathCategory)
	}

	book := domain.Book{CategoryID: englishCategory.ID}
	errBook := testDB.Create(&book).Error
	if errBook != nil {
		t.Fatalf("Failed to create category by error %v", errBook)
	}

	fakeBookID := domain.MustGetUUIDFromString("1698bbd6-e0c8-4957-a5a9-8c536970994b")

	type args struct {
		p *domain.Book
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "success update",
			args: args{
				&domain.Book{
					Model:       domain.Model{ID: book.ID},
					Name:        "book Name 1",
					Description: "example@gmail.com",
					CategoryID:  mathCategory.ID,
				},
			},
		},
		{
			name: "invalid category_id",
			args: args{
				&domain.Book{
					Model:       domain.Model{ID: book.ID},
					Name:        "book Name 1",
					Description: "example@gmail.com",
					CategoryID:  fakeBookID,
				},
			},
			wantErr: ErrNotFound,
		},
		{
			name: "success update",
			args: args{
				&domain.Book{
					Model:       domain.Model{ID: book.ID},
					Name:        "book Name 1",
					Description: "example@gmail.com",
					CategoryID:  englishCategory.ID,
				},
			},
		},
		{
			name: "failed update",
			args: args{
				&domain.Book{
					Model:       domain.Model{ID: fakeBookID},
					Name:        "book Name 1",
					Description: "example@gmail.com",
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

	book := domain.Book{}
	err = testDB.Create(&book).Error
	if err != nil {
		t.Fatalf("Failed to create book by error %v", err)
	}

	fakeBookID := domain.MustGetUUIDFromString("1698bbd6-e0c8-4957-a5a9-8c536970994b")

	type args struct {
		p *domain.Book
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.Book
		wantErr error
	}{
		{
			name: "success find correct book",
			args: args{
				&domain.Book{
					Model: domain.Model{ID: book.ID},
				},
			},
			want: &domain.Book{
				Model: domain.Model{ID: book.ID},
			},
		},
		{
			name: "failed find book by not exist book id",
			args: args{
				&domain.Book{
					Model: domain.Model{ID: fakeBookID},
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
		want    []domain.Book
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

	book := domain.Book{}
	err = testDB.Create(&book).Error
	if err != nil {
		t.Fatalf("Failed to create book by error %v", err)
	}

	fakeBookID := domain.MustGetUUIDFromString("1698bbd6-e0c8-4957-a5a9-8c536970994b")

	type args struct {
		p *domain.Book
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "success delete",
			args: args{
				&domain.Book{
					Name:        "This is book Name",
					Model:       domain.Model{ID: book.ID},
					Description: "example@gmail.com",
				},
			},
		},
		{
			name: "failed delete by not exist book id",
			args: args{
				&domain.Book{
					Model:       domain.Model{ID: fakeBookID},
					Name:        "This is book Name",
					Description: "example@gmail.com",
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
