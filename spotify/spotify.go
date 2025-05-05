package spotify

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"go.uber.org/zap"
	"golang.org/x/oauth2/clientcredentials"
)

type SpotifyService interface {
	GetObjectName(ctx context.Context, url string) (string, error)
	GetObjectType(ctx context.Context, url string) (SpotifyObjectType, error)
	GetPlaylistTracks(ctx context.Context, url string) ([]spotify.PlaylistItem, error)
}

type spotifyService struct {
	ClientID      string
	ClientSecret  string
	spotifyClient *spotify.Client
	log           *zap.Logger
}

func NewSpotifyService(ctx context.Context, clientID, clientSecret string, log *zap.Logger) SpotifyService {
	spotifyConfig := clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     spotifyauth.TokenURL,
	}

	token, err := spotifyConfig.Token(ctx)
	if err != nil {
		log.Fatal("failed to get token", zap.Error(err))
	}

	httpClient := spotifyauth.New().Client(context.Background(), token)
	spotifyClient := spotify.New(httpClient)

	return &spotifyService{
		spotifyClient: spotifyClient,
		log:           log,
	}
}

func (s *spotifyService) refreshToken(ctx context.Context) error {
	spotifyConfig := clientcredentials.Config{
		ClientID:     s.ClientID,
		ClientSecret: s.ClientSecret,
		TokenURL:     spotifyauth.TokenURL,
	}

	token, err := spotifyConfig.Token(ctx)
	if err != nil {
		s.log.Error("failed to refresh token", zap.Error(err))
		return err
	}

	httpClient := spotifyauth.New().Client(context.Background(), token)
	s.spotifyClient = spotify.New(httpClient)

	return nil
}

func (s *spotifyService) GetObjectName(ctx context.Context, url string) (string, error) {
	s.refreshToken(ctx)

	if !s.isValidSpotifyURL(url) {
		s.log.Error("invalid spotify url", zap.String("url", url))
		return "", errors.New("invalid spotify url")
	}

	objectType, err := s.GetObjectType(ctx, url)
	if err != nil {
		s.log.Error("failed to get object type", zap.Error(err))
		return "", err
	}

	id := s.getSpotifyID(url)
	if id == "" {
		s.log.Error("failed to get spotify id", zap.String("url", url))
		return "", errors.New("invalid spotify url")
	}

	var name string
	switch objectType {
	case SpotifyObjectTypePlaylist:
		playlist, err := s.spotifyClient.GetPlaylist(ctx, id)
		if err != nil {
			s.log.Error("failed to get playlist", zap.Error(err), zap.String("id", string(id)))
			return "", err
		}
		name = playlist.Name
	case SpotifyObjectTypeAlbum:
		album, err := s.spotifyClient.GetAlbum(ctx, id)
		if err != nil {
			s.log.Error("failed to get album", zap.Error(err), zap.String("id", string(id)))
			return "", err
		}
		name = album.Name
	case SpotifyObjectTypeTrack:
		track, err := s.spotifyClient.GetTrack(ctx, id)
		if err != nil {
			s.log.Error("failed to get track", zap.Error(err), zap.String("id", string(id)))
			return "", err
		}
		name = track.Name
	default:
		return "", errors.New("unknown object type")
	}

	return name, nil
}

func (s *spotifyService) GetObjectType(ctx context.Context, url string) (SpotifyObjectType, error) {
	if strings.Contains(url, "playlist") {
		return SpotifyObjectTypePlaylist, nil
	}
	if strings.Contains(url, "album") {
		return SpotifyObjectTypeAlbum, nil
	}
	if strings.Contains(url, "track") {
		return SpotifyObjectTypeTrack, nil
	}
	if strings.Contains(url, "artist") {
		return SpotifyObjectTypeArtist, nil
	}
	return "", errors.New("unknown object type")
}

func (s *spotifyService) isValidSpotifyURL(url string) bool {
	// Check if the URL starts with "https://open.spotify.com/"
	return strings.HasPrefix(url, "https://open.spotify.com/")
}

func (s *spotifyService) getSpotifyID(url string) spotify.ID {
	id := strings.Split(strings.Split(url, "/")[4], "?")[0]
	return spotify.ID(id)
}

func (s *spotifyService) GetPlaylistTracks(ctx context.Context, url string) ([]spotify.PlaylistItem, error) {
	s.refreshToken(ctx)

	if !s.isValidSpotifyURL(url) {
		return nil, errors.New("invalid spotify url")
	}

	id := s.getSpotifyID(url)
	if id == "" {
		return nil, errors.New("invalid spotify url")
	}

	var playlistItems []spotify.PlaylistItem
	itemsPage, err := s.spotifyClient.GetPlaylistItems(context.Background(), id)
	if err != nil {
		fmt.Println("Error getting playlist items:", err)
		return nil, err
	}

	if itemsPage.Total > spotify.Numeric(len(itemsPage.Items)) {
		total := int(itemsPage.Total)
		for i := 0; i < total; i += int(itemsPage.Limit) {
			items, err := s.spotifyClient.GetPlaylistItems(context.Background(), id, spotify.Limit(int(itemsPage.Limit)), spotify.Offset(i))
			if err != nil {
				fmt.Println("Error getting playlist items:", err)
				return nil, err
			}
			playlistItems = append(playlistItems, items.Items...)
		}
	} else {
		playlistItems = itemsPage.Items
	}

	return playlistItems, nil
}
