package index

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndex(t *testing.T) {
	idx := make(index)

	assert.Nil(t, idx.search("кофе"))
	assert.Nil(t, idx.search("булка"))

	idx.add([]document{{ID: 1, Text: "Съешь же ещё этих мягких французских булок, да выпей чаю!"}})
	// todo Починить stopwordFilter
	// assert.Nil(t, idx.search("же"))
	assert.Nil(t, idx.search("кофе"))
	assert.Equal(t, idx.search("чай"), []int{1})
	assert.Equal(t, idx.search("ЧаЙ"), []int{1})
	assert.Equal(t, idx.search("ВыПей"), []int{1})
	assert.Equal(t, idx.search("французский"), []int{1})

	// idx.add([]document{{ID: 2, Text: "Лошади кушают овес и сено."}})
	// todo Починить stopwordFilter
	//assert.Nil(t, idx.search("не"))
	// assert.Equal(t, idx.search("лошад"), []int{1, 2})
	// assert.Equal(t, idx.search("кушать"), []int{1, 2})
}
