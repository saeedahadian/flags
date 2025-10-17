package flags

import "testing"

func TestWhenMixedFlagTypesPassedIsParsed(t *testing.T) {
	args := []string{"--all", "-w", "12", "--parallel", "true"}

	flags, _ := ParseFlags(args)

	if len(flags) != 3 {
		t.Fatalf("Failed to generate 3 flags as should be.\n")
	}

	expectedFlags := map[string]*Flag{
		"all":      MustNewBoolFlag("all", "true"),
		"w":        MustNewIntFlag("w", "12"),
		"parallel": MustNewBoolFlag("parallel", "true"),
	}

	for _, flag := range flags {
		expectedFlag, exists := expectedFlags[flag.Key]
		if !exists {
			t.Fatalf("%s flag not found but expected.", flag.Key)
		}

		if expectedFlag.Value.String() != flag.Value.String() {
			t.Fatalf("%s flag is expected to have %s value but got %v.", flag.Key, expectedFlag.Value.String(), flag.Value.String())
		}
	}
}

func TestWhenSingleBooleanFlagPassedIsParsed(t *testing.T) {
	args := []string{"--all"}

	flags, err := ParseFlags(args)
	if err != nil {
		t.Fatalf("Failed while parsing flags with error: %s", err.Error())
	}

	if len(flags) != 1 {
		t.Fatalf("Should only parse a single flag.")
	}

	flag := flags[0]
	if flag.Key != "all" || flag.Value.String() != "true" {
		t.Fatalf("Didn't parse flags correctly.")
	}
}

func TestWhenEmptyArgsPassedNoFlagsAreParsed(t *testing.T) {
	args := []string{}

	flags, err := ParseFlags(args)
	if err != nil {
		t.Fatalf("Error while parsing the arguments: %s", err.Error())
	}

	if len(flags) != 0 {
		t.Fatalf("No flags should've been parsed.")
	}
}