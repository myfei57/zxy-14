package source

import (
	"maskhub/internal/detect"
	"maskhub/internal/store"
)

// Stats summarizes one source's fields and classification state.
type Stats struct {
	SourceID   string `json:"source_id"`
	Fields     int    `json:"fields"`
	Classified int    `json:"classified"`
	Reclassified int  `json:"reclassified"`
	Pending    int    `json:"pending"`
}

// StatsOf computes the field statistics of a source.
func StatsOf(state *store.State, sourceID string) Stats {
	fields := detect.Fields(state, sourceID)
	s := Stats{SourceID: sourceID, Fields: len(fields)}
	for _, f := range fields {
		switch f.Status {
		case store.FieldClassified:
			s.Classified++
		case store.FieldReclassified:
			s.Reclassified++
		default:
			s.Pending++
		}
	}
	return s
}

// ListStats returns statistics for every source.
func ListStats(state *store.State) []Stats {
	var out []Stats
	for _, src := range state.Sources() {
		out = append(out, StatsOf(state, src.ID))
	}
	return out
}
