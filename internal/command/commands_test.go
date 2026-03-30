package command

import "testing"

func TestParseBadInputs(t *testing.T) {

	inputs := []string{"", "GET"}

	for _, bad := range inputs {

		_, err := Parse(bad)

		if err == nil {
			t.Errorf("Parsing message '%s' should have raised an error", bad)
		}
	}
}

func TestParseSingleArgCommand(t *testing.T) {
	cmd, err := Parse("GET this_key")

	if err != nil {
		t.Errorf("Parsing error: %s", err.Error())
	}

	name, arg1 := cmd.Name(), cmd.Arg1()

	if name != "GET" {
		t.Errorf("Name does not match. Actual: '%s', Expected: '%s'", name, "GET")
	}

	if arg1 != "this_key" {
		t.Errorf("Arg1 does not match. Actual: '%s', Expected: '%s'", arg1, "this_key")
	}
}

func TestParseDualArgCommand(t *testing.T) {
	cmd, err := Parse("SET key value")

	if err != nil {
		t.Errorf("Parsing error: %s", err.Error())
	}

	name, arg1, arg2 := cmd.Name(), cmd.Arg1(), cmd.Arg2()

	if name != "SET" {
		t.Errorf("Name does not match. Actual '%s', Expected: '%s'", name, "SET")
	}

	if arg1 != "key" {
		t.Errorf("Name does not match. Actual '%s', Expected: '%s'", arg1, "key")
	}

	if arg2 != "value" {
		t.Errorf("Name does not match. Actual '%s', Expected: '%s'", arg2, "value")
	}
}
