package search

import (
	"context"
	"strconv"

	"github.com/elastic/elastic-transport-go/v8/elastictransport"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/rs/zerolog/log"
)

const index = "track"

type SearchService struct {
	client *elasticsearch.TypedClient
}

func New(cfg elasticsearch.Option) *SearchService {
	client, err := elasticsearch.NewTyped(
		cfg,
		elasticsearch.WithTransportOptions(elastictransport.WithTransport(&Transport{})),
	)
	if err != nil {
		panic(err)
	}

	return &SearchService{client: client}
}

func (s *SearchService) IndexTrack(ctx context.Context, id int, title string) error {

	_, err := s.client.Index(index).
		Id(strconv.Itoa(id)).
		Document(map[string]any{
			"id":    id,
			"title": title,
		}).Do(ctx)
	return err
}

func (s *SearchService) SearchTrack(ctx context.Context, q string, limit, offset int) ([]int, error) {
	log.Debug().Msg(q)

	resp, err := s.client.Search().
		Index(index).
		Query(&types.Query{
			MatchPhrasePrefix: map[string]types.MatchPhrasePrefixQuery{
				"title": {
					Query: q,
				},
			},
		}).
		Do(ctx)

	if err != nil {
		return nil, err
	}
	log.Debug().Int("res count", int(resp.Hits.Total.Value)).Send()
	ids := make([]int, 0, resp.Hits.Total.Value)

	for _, hit := range resp.Hits.Hits {
		id, _ := strconv.Atoi(*hit.Id_)
		ids = append(ids, id)
	}

	return ids, nil
}
