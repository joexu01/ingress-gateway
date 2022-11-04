package cert

import "testing"

func TestPath(t *testing.T) {
	path := Path("ca.crt")
	t.Log(path)
}
