package postgres

import (
	"context"
	"errorHandler/errorHandler/storage"
	"testing"
)

func TestCreateUser(t *testing.T) {
	s := newTestStorage(t)

	tests := []struct {
		name    string
		in      storage.User
		want    int64
		wantErr bool
	}{
		{
			name: "===========CREATE_USER_SUCCESS===============",
			in: storage.User{
				Name:          "user",
				Email:         "user@gmail.com",
				Password:      "123456",
				EmailVerified: "active",
			},
			want: 1,
		},
		{
			name: "CREATE_BLOG_SUCCESS",
			in: storage.User{
				Name:          "user",
				Email:         "user1@gmail.com",
				Password:      "123456",
				EmailVerified: "active",
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.CreateUser(context.Background(), tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Storage.CreateUser() = %v, want %v", got, tt.want)
			}
		})
	}
}
