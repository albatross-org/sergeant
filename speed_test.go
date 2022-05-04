package sergeant_test

import (
	"testing"

	"github.com/albatross-org/go-albatross/albatross"
	"github.com/albatross-org/sergeant"
)

func BenchmarkConfigLoad(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, err := sergeant.LoadConfig("")
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkWebServer(b *testing.B) {
	for n := 0; n < b.N; n++ {
		config, err := sergeant.LoadConfig("")
		if err != nil {
			b.Error(err)
			break
		}

		underlyingStore, err := albatross.FromConfig(config.Store)
		if err != nil {
			b.Error(err)
			break
		}

		store := sergeant.NewStore(underlyingStore, config)

		_, _, err = store.Set("all")
		if err != nil {
			b.Error(err)
			break
		}
	}
}
