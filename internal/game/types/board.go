package types

type Board struct {
	Data [][]string
}

func NewBoard(size int) *Board {
	data := make([][]string, size)
	for i := range data {
		data[i] = make([]string, size)
		for j := range data[i] {
			data[i][j] = ""
		}
	}
	return &Board{
		Data: data,
	}
}

func (b *Board) Size() int {
	return len(b.Data)
}

func (b *Board) FilledSquares() int {
	count := 0
	for i := range b.Data {
		for j := range b.Data[i] {
			if b.Data[i][j] != "" {
				count++
			}
		}
	}
	return count
}

func IsValidLetter(letter string) bool {
	return len(letter) == 1 && letter[0] >= 'A' && letter[0] <= 'Z'
}
