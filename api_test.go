package ympp

import (
	"context"
	"reflect"
	"testing"
)

func TestLibraryAPI_GetOwnerInfo(t *testing.T) {
	type args struct {
		ctx   context.Context
		login string
	}
	tests := []struct {
		name     string
		api      *LibraryAPI
		args     args
		wantInfo OwnerInfo
		wantErr  bool
	}{
		{
			name: "existing owner info",
			api:  NewLibraryAPI(MusicHandlersURL),
			args: args{
				ctx:   context.Background(),
				login: "winpooh32",
			},
			wantInfo: OwnerInfo{
				Visibility: "public",
				HasTracks:  true,
				Owner: Owner{
					UID:   "280506128",
					Login: "winpooh32",
					Name:  "WinPooh32",
				},
			},
			wantErr: false,
		},
		{
			name: "not found",
			api:  NewLibraryAPI(MusicHandlersURL),
			args: args{
				ctx:   context.Background(),
				login: "winpooh32222",
			},
			wantInfo: OwnerInfo{},
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := tt.api
			gotInfo, err := api.GetOwnerInfo(tt.args.ctx, tt.args.login)
			if (err != nil) != tt.wantErr {
				t.Errorf("LibraryAPI.GetOwnerInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotInfo, tt.wantInfo) {
				t.Errorf("LibraryAPI.GetOwnerInfo() = %v, want %v", gotInfo, tt.wantInfo)
			}
		})
	}
}

func TestLibraryAPI_GetLibrary(t *testing.T) {
	type args struct {
		ctx   context.Context
		login string
	}
	tests := []struct {
		name    string
		api     *LibraryAPI
		args    args
		wantErr bool
	}{
		{
			name: "info",
			api:  NewLibraryAPI(MusicHandlersURL),
			args: args{
				ctx:   context.Background(),
				login: "winpooh32",
			},
			wantErr: false,
		},
		{
			name: "info not found",
			api:  NewLibraryAPI(MusicHandlersURL),
			args: args{
				ctx:   context.Background(),
				login: "winpooh32312304",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := tt.api
			gotLibrary, err := api.GetLibrary(tt.args.ctx, tt.args.login)
			if (err != nil) != tt.wantErr {
				t.Errorf("LibraryAPI.GetLibrary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(gotLibrary.Playlists) == 0 {
				t.Errorf("len(gotLibrary.Playlists) == 0")
			}
		})
	}
}

func TestLibraryAPI_GetPlaylist(t *testing.T) {
	type args struct {
		ctx      context.Context
		login    string
		playlist int
	}
	tests := []struct {
		name    string
		api     *LibraryAPI
		args    args
		wantErr bool
	}{
		{
			name: "tracklist found",
			api:  NewLibraryAPI(MusicHandlersURL),
			args: args{
				ctx:      context.Background(),
				login:    "winpooh32",
				playlist: 3,
			},
			wantErr: false,
		},
		{
			name: "not found",
			api:  NewLibraryAPI(MusicHandlersURL),
			args: args{
				ctx:      context.Background(),
				login:    "winpooh323123432",
				playlist: 3,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := tt.api
			gotPlaylist, err := api.GetPlaylist(tt.args.ctx, tt.args.login, tt.args.playlist)
			if (err != nil) != tt.wantErr {
				t.Errorf("LibraryAPI.GetPlaylist() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(gotPlaylist.Tracks) == 0 {
				t.Errorf("len(gotTracklist.Tracks) == 0")
			}
		})
	}
}
