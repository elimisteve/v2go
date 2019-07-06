// Steve Phillips / elimisteve
// 2019.07.05

package translate

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTranslateAndRunFiles(t *testing.T) {
	correct := "Hello, World!\n"
	passingHelloWorld := []string{
		"../test_v_files/hello_world.v",
		"../test_v_files/hello_world_interpolated.v",
		"../test_v_files/hello_world_interpolated2.v",
		"../test_v_files/hello_world_module.v",
	}
	for _, vFilename := range passingHelloWorld {
		out, err := TranslateAndRunFiles([]string{vFilename})
		if err != nil {
			t.Fatalf("Error running TranslateAndRunFiles(%q): %v", vFilename, err)
		}
		assert.Equal(t, correct, string(out))
	}
}

func TestTranslateAndRunFiles2(t *testing.T) {
	correctLen := 31
	passing := []string{
		"../test_v_files/links_scraper2.v",
	}
	for _, vFilename := range passing {
		out, err := TranslateAndRunFiles([]string{vFilename})
		if err != nil {
			t.Fatalf("Error running TranslateAndRunFiles(%q): %v", vFilename, err)
		}
		assert.Equal(t, correctLen, len(strings.Split(string(out), "\n")))
	}
}
