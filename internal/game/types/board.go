package types

type Board struct {
	Data [][]string
}

func NewBoard(size uint) *Board {
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
