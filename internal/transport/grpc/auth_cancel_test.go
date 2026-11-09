package grpc

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/forgelab-dev/lfa-forgenet/internal/auth"
	"github.com/forgelab-dev/lfa-forgenet/internal/store"
	forgenetv1 "github.com/forgelab-dev/lfa-forgenet/pkg/gen/forgenet/v1"
)

// failingUsers échoue toujours avec err.
type failingUsers struct{ err error }

func (f failingUsers) GetUserByID(context.Context, string) (*store.User, error) { return nil, f.err }

func TestGetCurrentUser_ClientCancellationIsNotAnOutage(t *testing.T) {
	cases := []struct {
		name    string
		err     error
		code    codes.Code
		journal bool
	}{
		{"appel abandonné par le client", fmt.Errorf("lecture de l'utilisateur : %w", context.Canceled), codes.Canceled, false},
		{"délai du client dépassé", context.DeadlineExceeded, codes.DeadlineExceeded, false},
		{"base en panne", errors.New("disque plein"), codes.Unavailable, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var logs bytes.Buffer
			srv := NewAuthServiceServer(nil, failingUsers{err: tc.err}, zerolog.New(&logs))
			ctx := auth.WithUserContext(context.Background(), &auth.Claims{UserID: "usr_1", Role: store.RoleMember})

			_, err := srv.GetCurrentUser(ctx, &forgenetv1.GetCurrentUserRequest{})

			assert.Equal(t, tc.code, status.Code(err))
			assert.Equal(t, tc.journal, logs.Len() > 0, "journal : %s", logs.String())
		})
	}
}
