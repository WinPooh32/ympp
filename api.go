package ympp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/andelf/go-curl"
)

const MusicHandlersURL = "https://music.yandex.ru/handlers"

var ErrNotFound = errors.New("not found")

type LibraryAPI struct {
	userAgent string

	baseurl string
}

func NewDefaultLibraryAPI() *LibraryAPI {
	return NewLibraryAPI(
		MusicHandlersURL,
	)
}

func NewLibraryAPI(baseurl string) *LibraryAPI {
	return &LibraryAPI{
		userAgent: "curl/7.68.0",
		baseurl:   baseurl,
	}
}

func (api *LibraryAPI) SetUserAgent(ua string) {
	api.userAgent = ua
}

func (api *LibraryAPI) GetOwnerInfo(ctx context.Context, login string) (info OwnerInfo, err error) {
	q := url.Values{}
	q.Add("owner", login)

	err = api.get(ctx, "library.jsx", q, &info)
	if err != nil {
		return OwnerInfo{}, err
	}

	return info, nil
}

func (api *LibraryAPI) GetLibrary(ctx context.Context, login string) (library Library, err error) {
	q := url.Values{}
	q.Add("owner", login)
	q.Add("filter", "playlists")

	err = api.get(ctx, "library.jsx", q, &library)
	if err != nil {
		return Library{}, err
	}

	return library, nil
}

func (api *LibraryAPI) GetPlaylist(ctx context.Context, login string, playlist int) (playlistInfo PlaylistWithTracks, err error) {
	q := url.Values{}
	q.Add("owner", login)
	q.Add("kinds", strconv.Itoa(playlist))

	var tl Tracklist

	err = api.get(ctx, "playlist.jsx", q, &tl)
	if err != nil {
		return PlaylistWithTracks{}, err
	}

	return tl.Playlist, nil
}

func (api *LibraryAPI) get(ctx context.Context, endpoint string, query url.Values, dst interface{}) (err error) {
	endpoint, err = url.JoinPath(api.baseurl, endpoint)
	if err != nil {
		return fmt.Errorf("join url path: %w", err)
	}

	url, err := url.Parse(endpoint)
	if err != nil {
		return fmt.Errorf("parse url: %w", err)
	}

	url.RawQuery = query.Encode()

	// req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	// if err != nil {
	// 	return fmt.Errorf("new request: %w", err)
	// }

	// req.Header = api.header

	// resp, err := api.client.Do(req)
	// if err != nil {
	// 	return fmt.Errorf("do request: %w", err)
	// }
	// defer func() {
	// 	err = errors.Join(err, resp.Body.Close())
	// }()

	// switch status := resp.StatusCode; status {
	// case http.StatusOK:
	// 	break
	// case http.StatusNotFound:
	// 	return ErrNotFound
	// default:
	// 	return fmt.Errorf("unexpected status code = %d", status)
	// }

	// var reader io.ReadCloser
	// switch resp.Header.Get("Content-Encoding") {
	// case "gzip":
	// 	reader, err = gzip.NewReader(resp.Body)
	// 	defer reader.Close()
	// default:
	// 	reader = resp.Body
	// }

	// data, err := io.ReadAll(reader)
	// if err != nil {
	// 	return fmt.Errorf("read response body: %w", err)
	// }

	easy := curl.EasyInit()
	defer easy.Cleanup()

	err = easy.Setopt(curl.OPT_USERAGENT, api.userAgent)
	if err != nil {
		return fmt.Errorf("set user-agent: %w", err)
	}

	err = easy.Setopt(curl.OPT_HTTPHEADER, []string{
		"Accept: */*",
		"Referer: https://music.yandex.ru/",
	})
	if err != nil {
		return fmt.Errorf("set header: %w", err)
	}

	err = easy.Setopt(curl.OPT_URL, url.String())
	if err != nil {
		return fmt.Errorf("set url: %w", err)
	}

	var data []byte

	err = easy.Setopt(curl.OPT_WRITEFUNCTION, func(buf []byte, _ interface{}) bool {
		data = append(data, buf...)
		return true
	})
	if err != nil {
		return fmt.Errorf("set write callback: %w", err)
	}

	err = easy.Perform()
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}

	code, err := easy.Getinfo(curl.INFO_HTTP_CODE)
	if err != nil {
		return fmt.Errorf("get status code: %w", err)
	}

	switch status := code.(int); status {
	case http.StatusOK:
		break
	case http.StatusNotFound:
		return ErrNotFound
	default:
		return fmt.Errorf("unexpected status code = %d", status)
	}

	err = json.Unmarshal(data, dst)
	if err != nil {
		return fmt.Errorf("parse json: %w", err)
	}

	return nil
}
