package address

import (
	"testing"

	"github.com/payshares/go/support/errors"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	cases := []struct {
		Name            string
		Domain          string
		ExpectedAddress string
	}{
		{"scott", "payshares.org", "scott*payshares.org"},
		{"", "payshares.org", "*payshares.org"},
		{"scott", "", "scott*"},
	}

	for _, c := range cases {
		actual := New(c.Name, c.Domain)
		assert.Equal(t, actual, c.ExpectedAddress)
	}
}

func TestSplit(t *testing.T) {
	cases := []struct {
		CaseName       string
		Address        string
		ExpectedName   string
		ExpectedDomain string
		ExpectedError  error
	}{
		{"happy path", "scott*payshares.org", "scott", "payshares.org", nil},
		{"blank", "", "", "", ErrInvalidAddress},
		{"blank name", "*payshares.org", "", "", ErrInvalidName},
		{"blank domain", "scott*", "", "", ErrInvalidDomain},
		{"invalid domain", "scott*--3.com", "", "", ErrInvalidDomain},
	}

	for _, c := range cases {
		name, domain, err := Split(c.Address)

		if c.ExpectedError == nil {
			assert.Equal(t, name, c.ExpectedName)
			assert.Equal(t, domain, c.ExpectedDomain)
		} else {
			assert.Equal(t, errors.Cause(err), c.ExpectedError)
		}
	}
}
