package postgres

import (
	"context"
	"errorHandler/errorHandler/storage"
	//"log"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCreateCategory(t *testing.T) {
	s := newTestStorage(t)

	tests := []struct {
		name    string
		in      storage.Category
		want    int64
		wantErr bool
	}{
		{
			name: "CREATE_BLOG_SUCCESS",
			in: storage.Category{
				Title: "This is title",
			},
			want: 1,
		},
		{
			name: "CREATE_BLOG_SUCCESS",
			in: storage.Category{
				Title: "This is title 1",
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.Create_sto(context.Background(), tt.in)
			//log.Printf("%#v", got)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Storage.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCategory(t *testing.T) {
	s := newTestStorage(t)
	tests := []struct {
		name    string
		in      int64
		want    storage.Category
		wantErr bool
	}{
		{
			name: "GET_Category_SUCCESS",
			in:   1,
			want: storage.Category{
				ID:    1,
				Title: "This is title",
			},
		},
		{
			name:    "FAILED_Category_DOES_NOT_EXIST",
			in:      100,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			category, err := s.Get_sto(context.Background(), tt.in)
			//log.Printf("%#v", category)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !cmp.Equal(category, tt.want) {
				t.Errorf("Diff: got -, want += %v", cmp.Diff(category, tt.want))
			}
		})
	}
}

func TestAllDataCategory(t *testing.T) {
	s := newTestStorage(t)
	tests := []struct {
		name    string
		want    []storage.Category
		wantErr bool
	}{
		{
			name: "=======GET_All_data_SUCCESS========",
			want: []storage.Category{
				{
					ID:    1,
					Title: "This is title",
				},
				{
					ID:    2,
					Title: "This is title 1",
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			gotList, err := s.Get_all_Data(context.Background())
			//log.Printf("======???========== %#v", gotList)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			sort.Slice(tt.want, func(i, j int) bool {
				return tt.want[i].ID < tt.want[j].ID
			})

			sort.Slice(gotList, func(i, j int) bool {
				return gotList[i].ID < gotList[j].ID
			})

			for i, got := range gotList {

				if !cmp.Equal(got, tt.want[i]) {
					t.Errorf("Diff: got -, want += %v", cmp.Diff(got, tt.want[i]))
				}

			}

		})
	}
}

func TestUpdateCategory(t *testing.T) {
	s := newTestStorage(t)
	tests := []struct {
		name    string
		in      storage.Category
		want    *storage.Category
		wantErr bool
	}{
		{
			name: "UPDATE_Category_SUCCESS",
			in: storage.Category{
				ID:    1,
				Title: "This is title Update",
			},
			want: &storage.Category{
				ID:    1,
				Title: "This is title Update",
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := s.Update(context.Background(), tt.in)
			//log.Printf("%#v", got)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestDeleteCategory(t *testing.T) {
	s := newTestStorage(t)
	tests := []struct {
		name    string
		in      int64
		wantErr bool
	}{
		{
			name: "DELETE_Category_SUCCESS",
			in:   2,
		},
		{
			name:    "FAILED_TO_DELETE_Category_DOES_NOT_EXIST",
			in:      100,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := s.Delete(context.Background(), tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
