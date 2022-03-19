package postgres

import (
	"context"
	"errorHandler/errorHandler/storage"

	"sort"

	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCreateBlog(t *testing.T) {
	s := newTestStorage(t)
	tests := []struct {
		name    string
		in      storage.Blog
		want    int64
		wantErr bool
	}{
		{
			name: "CREATE_Blog_SUCCESS",
			in: storage.Blog{
				CatID:       1,
				Title:       "This is Title",
				Description: "This is description",
				Image:       "This is image",
			},
			want: 1,
		},
		{
			name: "CREATE_Blog_SUCCESS",
			in: storage.Blog{
				CatID:       1,
				Title:       "This is Title 2",
				Description: "This is description 2",
				Image:       "This is image 2",
			},
			want: 2,
		},
		{
			name: "FAILED_DUPLICATE_TITLE",
			in: storage.Blog{
				CatID:       1,
				Title:       "This is Title",
				Description: "This is description",
				Image:       "This is image",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			got, err := s.Create(context.Background(), tt.in)
			//fmt.Println(got,err)
			//log.Fatal(got,err)

			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.Create Blog() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Storage.Create Blog() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestListBlog(t *testing.T) {
	s := newTestStorage(t)
	tests := []struct {
		name    string
		want    []storage.Blog
		wantErr bool
	}{
		{
			name: "GET_ALL_Blog_SUCCESS",

			want: []storage.Blog{
				{
					ID:          1,
					CatID:       1,
					Title:       "This is Title",
					Description: "This is description",
					Image:       "This is image",
					CatName:     "This is title Update",
				},
				{
					ID:          2,
					CatID:       1,
					Title:       "This is Title 2",
					Description: "This is description 2",
					Image:       "This is image 2",
					CatName:     "This is title Update",
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			gotList, err := s.ListBlog(context.Background())

			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			wantList := tt.want

			sort.Slice(wantList, func(i, j int) bool {
				return wantList[i].ID < wantList[j].ID
			})

			sort.Slice(gotList, func(i, j int) bool {
				return gotList[i].ID < gotList[j].ID
			})

			for i, got := range gotList {

				if !cmp.Equal(got, wantList[i]) {
					t.Errorf("Diff: got -, want += %v", cmp.Diff(got, wantList[i]))
				}

			}
		})
	}
}

func TestGetBlog(t *testing.T) {

	s := newTestStorage(t)

	tests := []struct {
		name    string
		in      int64
		want    storage.Blog
		wantErr bool
	}{
		{
			name: "GET_Blog_SUCCESS",
			in:   1,
			want: storage.Blog{
				ID:          1,
				CatID:       1,
				Title:       "This is Title",
				Description: "This is description",
				Image:       "This is image",
			},
		},
		{
			name: "GET_Blog_SUCCESS",
			in:   2,
			want: storage.Blog{
				ID:          2,
				CatID:       1,
				Title:       "This is Title 2",
				Description: "This is description 2",
				Image:       "This is image 2",
			},
		},
		{
			name:    "FAILED_TO_GET_Blog",
			in:      3,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.GetBlog(context.Background(), tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("Diff: got -, want += %v", cmp.Diff(got, tt.want))
			}

		})
	}
}

func TestUpdateBlog(t *testing.T) {

	s := newTestStorage(t)

	tests := []struct {
		name    string
		in      storage.Blog
		wantErr bool
	}{
		{
			name: "UPDATE_Blog_SUCCESS",
			in: storage.Blog{
				ID:          1,
				CatID:       1,
				Title:       "This is Title updated",
				Description: "This is description updated",
				Image:       "This is image updated",
			},
		},
		{
			name: "FAILED_TO_UPDATE_Blog",
			in: storage.Blog{
				ID:          4,
				CatID:       1,
				Title:       "This is Title 3",
				Description: "This is description 3",
				Image:       "This is image 3",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := s.UpdateBlog(context.Background(), tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
func TestDeleteBlog(t *testing.T) {

	s := newTestStorage(t)

	tests := []struct {
		name    string
		in      int64
		want    bool
		wantErr bool
	}{
		{
			name: "DELETE_BLOG_SUCCESS",
			in:   1,
			want: true,
		},
		{
			name:    "FAILED_TO_DELETE_POST",
			in:      3,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := s.BlogDelete(context.Background(), tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
