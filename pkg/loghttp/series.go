package loghttp

import (
	"net/http"
	"sort"

	"github.com/grafana/loki/pkg/logproto"
)

type SeriesResponse struct {
	Status string     `json:"status"`
	Data   []LabelSet `json:"data"`
}

func ParseSeriesQuery(r *http.Request) (*logproto.SeriesRequest, error) {
	start, end, err := bounds(r)
	if err != nil {
		return nil, err
	}

	xs := r.Form["match"]
	// Prometheus encodes with `match[]`; we use both for compatibility.
	ys := r.Form["match[]"]

	deduped := union(xs, ys)
	sort.Strings(deduped)

	// ensure matchers are valid before fanning out to ingesters/store as well as returning valuable parsing errors
	// instead of 500s
	_, err = Match(deduped)
	if err != nil {
		return nil, err
	}

	return &logproto.SeriesRequest{
		Start:  start,
		End:    end,
		Groups: deduped,
	}, nil

}

func union(cols ...[]string) []string {
	m := map[string]struct{}{}

	for _, col := range cols {
		for _, s := range col {
			m[s] = struct{}{}
		}
	}

	res := make([]string, 0, len(m))
	for k := range m {
		res = append(res, k)
	}

	return res
}
