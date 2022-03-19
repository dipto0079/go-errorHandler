package postgres

import (
	"context"
	"sort"
	"testing"

	"errorHandler/errorHandler/storage"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestErrors(t *testing.T) {
	ts := newTestStorage(t)
	var id string

	tests := []struct {
		name    string
		in      storage.ErrorHandler
		want    string
		wantErr bool
	}{
		{
			name: "Success_Error_Handler_Create",
			in: storage.ErrorHandler{
				ErrorCode:    "00087",
				ErrorDetails: "This is test Error",
				EnvType:      "development",
				CreatedBy:    "development",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := ts.CreateError(context.TODO(), tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.CreateError() error%v, wantErr %v", err, tt.wantErr)
				return
			}
			id = got
		})
	}

	getTests := []struct {
		name    string
		in      string
		want    *storage.ErrorHandler
		wantErr bool
	}{
		{
			name: "GET_Error_Handler_SUCCESS",
			in:   id,
			want: &storage.ErrorHandler{
				ID:           id,
				ErrorCode:    "00087",
				ErrorDetails: "This is test Error",
				EnvType:      "development",
				CreatedBy:    "development",
			},
		},
	}

	for _, tt := range getTests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := ts.GetError(context.TODO(), tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.GetError() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("Diff = %v", cmp.Diff(got, tt.want))
			}
		})
	}

	allError := []struct {
		name    string
		want    []storage.ErrorHandler
		wantErr bool
	}{
		{
			name: "All_Error_Handler_SUCCESS",
			want: []storage.ErrorHandler{
				{
					ID:           id,
					ErrorCode:    "00087",
					ErrorDetails: "This is test Error",
					EnvType:      "development",
					CreatedBy:    "development",
				},
			},
		},
	}

	for _, tt := range allError {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			gotList, err := ts.ListError(context.Background())

			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.AllErrorHandler() error = %v, wantErr %v", err, tt.wantErr)
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

				tOps := []cmp.Option{
					cmpopts.IgnoreFields(storage.ErrorHandler{}, "CreatedAt", "CreatedBy", "DeletedAt", "DeleteByEnvType"),
				}

				if !tt.wantErr && !cmp.Equal(tt.want[i], got, tOps...) {
					t.Errorf("Storage.GetErrorHandler() = + got, - want: %+v", cmp.Diff(tt.want[i], got))
				}
			}
		})
	}

	deleteTests := []struct {
		name    string
		in      string
		wantErr bool
	}{
		{
			name: "DELETE_Error_Handler_SUCCESS",
			in:   id,
		},
	}

	for _, tt := range deleteTests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := ts.DeleteError(context.TODO(), tt.in, "1")
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.DELETEErrorHandler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}

}
