package di114514

import (
	"math/rand"
	"testing"
)

func TestDefineInt(t *testing.T) {
	c := NewContainer()

	err := c.Define("name", 101)
	if err != nil {
		t.Fail()
	}

	if v, ok := c.GetInstance("name").(int); ok {
		if v != 101 {
			t.Fail()
		}
	} else {
		t.Fail()
	}

	if c.GetInstance("not_existing_name") != nil {
		t.Fail()
	}

	if c.NewInstance("name") != nil {
		t.Fail()
	}
}

func TestDefineString(t *testing.T) {
	c := NewContainer()

	err := c.Define("name2", "alpen is not alpine")
	if err != nil {
		t.Fail()
	}

	if v, ok := c.GetInstance("name2").(string); ok {
		if v != "alpen is not alpine" {
			t.Fail()
		}
	} else {
		t.Fail()
	}

	if c.GetInstance("not_existing_name") != nil {
		t.Fail()
	}

	if c.NewInstance("name2") != nil {
		t.Fail()
	}
}

func TestDefine(t *testing.T) {
	c := NewContainer()
	err := c.Define("runes10", func() (interface{}) {
		return randStringRunes(10)
	})
	if err != nil {
		t.Fail()
	}

	if runes10A, ok := c.GetInstance("runes10").(string); ok {
		if runes10B, ok := c.GetInstance("runes10").(string); ok {
			if runes10A != runes10B {
				t.Fail()
			}
		} else {
			t.Fail()
		}
	} else {
		t.Fail()
	}

	if runes10C, ok := c.NewInstance("runes10").(string); ok {
		if runes10D, ok := c.NewInstance("runes10").(string); ok {
			if runes10C == runes10D {
				t.Fail()
			}
		} else {
			t.Fail()
		}
	} else {
		t.Fail()
	}
}

func TestDefineWithContainer(t *testing.T) {
	c := NewContainer()
	err := c.Define("runes10", func() (interface{}) {
		return randStringRunes(10)
	})
	if err != nil {
		t.Fail()
	}

	err = c.Define("runes20", func(c ContainerInterface) (interface{}) {
		if runes10, ok := c.GetInstance("runes10").(string); ok {
			return "[[[[[" + runes10 + "]]]]]";
		}

		return nil
	})
	if err != nil {
		t.Fail()
	}

	if runes10, ok := c.GetInstance("runes10").(string); ok {
		if runes20, ok := c.GetInstance("runes20").(string); ok {
			if "[[[[[" + runes10 + "]]]]]" != runes20 {
				t.Fail()
			}
		} else {
			t.Fail()
		}
	} else {
		t.Fail()
	}
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}