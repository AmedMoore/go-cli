package cli

import "testing"

func Test_App_Vars(t *testing.T) {
	app := NewApp([]string{})

	name := "foo"
	app.Set("name", name)
	got, ok := app.Get("name").(string)

	if !ok || got != name {
		t.Errorf("app.Get(\"name\") = %s, expected %s", got, name)
	}
}
