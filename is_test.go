package is_test

import (
	"testing"

	"github.com/matryer/is"
	"github.com/stretchr/testify/require"
)

type testType struct {
	Username string `is:"required,lower"`
	Email    string `is:"required,email"`
	Normal   string
	Number   int64 `is:"nonzero"`
}

// Valid

func TestIs(t *testing.T) {

	obj := testType{Username: "MatRyer", Email: "test@test.com", Number: 1}

	probs, err := is.Valid(&obj)

	require.NoError(t, err)
	require.Equal(t, 0, len(probs))

	require.Equal(t, "matryer", obj.Username)

}

func TestIsnt(t *testing.T) {

	obj := testType{Username: "", Email: "test@test"}

	probs, err := is.Valid(&obj)

	require.NoError(t, err)
	require.Equal(t, 3, len(probs))
	require.Equal(t, "Username is required", probs["Username"].Error())
	require.Equal(t, "Email is not a valid email address", probs["Email"].Error())
	require.Equal(t, "Number cannot be zero", probs["Number"].Error())

}

func TestLower(t *testing.T) {

	var v interface{}
	var err error
	v, err = is.DefaultValidator.Handlers["lower"]("MonKEY")
	require.NoError(t, err)
	require.Equal(t, "monkey", v)

}
