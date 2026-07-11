package webapi

import (
	"context"
	"encoding/json"
	"net/http"
	"pulse/internal/entity"
	"pulse/pkg/httpclient"
	"time"
)

type Actor struct {
	Login string `json:"login"`
}

type Repo struct {
	Name string `json:"name"`
}

type githubEvent struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	Actor     Actor     `json:"actor"`
	Repo      Repo      `json:"repo"`
}

type Adapter struct {
	client  *httpclient.Client
	baseUrl string
	token   string
}

func NewAdapter(client *httpclient.Client, url string, token string) *Adapter {
	return &Adapter{
		client:  client,
		baseUrl: url,
		token:   token,
	}
}

func (a *Adapter) Fetch(ctx context.Context) ([]*entity.Event, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", a.baseUrl+"/events", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	if a.token != "" {
		req.Header.Set("Authorization", "Bearer "+a.token)
	}

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ghEvents []githubEvent

	err = json.NewDecoder(resp.Body).Decode(&ghEvents)
	if err != nil {
		return nil, err
	}

	var result []*entity.Event
	for _, item := range ghEvents {
		event := &entity.Event{
			ID:         item.ID,
			ExternalID: item.ID,
			Source:     entity.SourceGitHub,
			Type:       entity.EventType(item.Type),
			Title:      item.Repo.Name,
			OccuredAt:  item.CreatedAt,
			CollectedAt:  time.Now(),
		}
		result = append(result, event)
	}

	return result, nil
}
