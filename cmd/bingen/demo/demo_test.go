package demo

import (
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/niubaoshu/gotiny"
	"github.com/stretchr/testify/require"
)

func TestDemo(t *testing.T) {
	var a GetTab
	gofakeit.Struct(&a)
	_ = a
	// a = A{
	// 	A0: "hello",
	// 	A1: 1,
	// 	A2: []int{1, 2, 3},
	// 	A3: [][]int{{1, 2, 3}, {4, 5, 6}},
	// 	A4: [][][]int{{{1, 2, 3}, {4, 5, 6}}, {{7, 8, 9}, {10, 11, 12}}},
	// 	A5: map[string]int{"a": 1, "b": 2},
	// }
	a1, err := a.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	_ = a1

	a2 := new(GetTab)
	err = a2.UnmarshalBinary(a1)
	if err != nil {
		t.Fatal(err)
	}
	require.Equal(t, a, *a2)
}

func BenchmarkExample(b *testing.B) {
	var user GetTab
	gofakeit.Struct(&user)
	b.Run("MarshalBinary", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := user.MarshalBinary()
			if err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("gotiny MarshalBinary", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			gotiny.Marshal(&user)
		}
	})

	b.Run("UnmarshalBinary", func(b *testing.B) {
		r1, err := user.MarshalBinary()
		if err != nil {
			b.Fatal(err)
		}
		var user2 = &GetTab{}
		for i := 0; i < b.N; i++ {
			err = user2.UnmarshalBinary(r1)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("gotiny UnmarshalBinary", func(b *testing.B) {
		r1 := gotiny.Marshal(&user)
		var user2 = &GetTab{}
		for i := 0; i < b.N; i++ {
			gotiny.Unmarshal(r1, user2)
		}
	})
}
