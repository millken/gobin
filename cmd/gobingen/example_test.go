package main

import (
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/niubaoshu/gotiny"
	"github.com/stretchr/testify/require"
)

func TestExample(t *testing.T) {
	var user AAA
	gofakeit.Struct(&user)
	r1, err := user.MarshalBinary()
	require.NoError(t, err)
	var user2 = &AAA{}
	err = user2.UnmarshalBinary(r1)
	require.NoError(t, err)
	require.Equal(t, user, *user2)
}

func TestTiny(t *testing.T) {
	var user AAA
	fmt.Println(gotiny.GetName(&user))
	gofakeit.Struct(&user)
	r1 := gotiny.Marshal(&user)
	var user2 = &AAA{}
	gotiny.Unmarshal(r1, user2)
	require.Equal(t, user, *user2)
}

func BenchmarkExample(b *testing.B) {
	var user AAA
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
		var user2 = &AAA{}
		for i := 0; i < b.N; i++ {
			err = user2.UnmarshalBinary(r1)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("gotiny UnmarshalBinary", func(b *testing.B) {
		r1 := gotiny.Marshal(&user)
		var user2 = &AAA{}
		for i := 0; i < b.N; i++ {
			gotiny.Unmarshal(r1, user2)
		}
	})
}
