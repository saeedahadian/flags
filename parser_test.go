package flags

import "testing"

func TestWhenMixedFlagTypesPassedIsParsed(t *testing.T) {
	args := []string{"--all", "-w", "12", "--parallel", "true"}

	flags, _ := ParseFlags(args)

	if len(flags) != 3 {
		t.Errorf("Failed to generate 3 flags as should be.\n")
	}

	expectedFlags := map[string]*Flag{
		"all":      MustNewBoolFlag("all", "true"),
		"w":        MustNewIntFlag("w", "12"),
		"parallel": MustNewBoolFlag("parallel", "true"),
	}

	for _, flag := range flags {
		expectedFlag, exists := expectedFlags[flag.Key]
		if !exists {
			t.Errorf("%s flag not found but expected.", flag.Key)
		}

		if expectedFlag.Value.String() != flag.Value.String() {
			t.Errorf("%s flag is expected to have %s value but got %v.", flag.Key, expectedFlag.Value.String(), flag.Value.String())
		}
	}
}
