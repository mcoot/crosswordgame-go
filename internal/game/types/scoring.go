package types

type ScoringDirection string

const (
	ScoringDirectionHorizontal ScoringDirection = "horizontal"
	ScoringDirectionVertical   ScoringDirection = "vertical"
)

type ScoredWord struct {
	Word        string           `json:"word"`
	Score       int              `json:"score"`
	Direction   ScoringDirection `json:"direction"`
	StartRow    int              `json:"start_row"`
	StartColumn int              `json:"start_column"`
}
