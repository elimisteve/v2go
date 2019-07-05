// Steve Phillips / elimisteve
// 2019.07.05

package translate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTranslateAndRunFiles(t *testing.T) {
	correct := "Hello, World!\n"
	passingHelloWorld := []string{
		"../test_v_files/hello_world.v",
		"../test_v_files/hello_world_interpolated.v",
		"../test_v_files/hello_world_interpolated2.v",
	}
	for _, vFilename := range passingHelloWorld {
		out, err := TranslateAndRunFiles([]string{vFilename})
		if err != nil {
			t.Fatalf("Error running TranslateAndRunFiles(%q): %v", vFilename, err)
		}
		assert.Equal(t, correct, string(out))
	}
}
