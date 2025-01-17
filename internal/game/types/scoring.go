package types

type ScoringDirection string

const (
	ScoringDirectionHorizontal ScoringDirection = "horizontal"
	ScoringDirectionVertical   ScoringDirection = "vertical"
)

type ScoredWord struct {
	Word       string           `json:"word"`
	Score      int              `json:"score"`
	Direction  ScoringDirection `json:"direction"`
	StartIndex int              `json:"start_index"`
	EndIndex   int              `json:"end_index"`
}
