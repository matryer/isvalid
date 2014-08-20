package is_test

import (
	"testing"

	"github.com/matryer/is"
	"github.com/stretchr/testify/require"
)

type testType struct {
	Username string `is:"required,lower"`
	Email    string `is:"required,email"`
}

func TestIs(t *testing.T) {

	obj := &testType{Username: "MatRyer", Email: "test@test.com"}

	probs, err := is.Valid(obj)

	require.NoError(t, err)
	require.Equal(t, 0, len(probs))

	require.Equal(t, "matryer", obj.Username)

}

func TestIsnt(t *testing.T) {

	obj := &testType{Username: "", Email: "test@test"}

	probs, err := is.Valid(obj)

	require.NoError(t, err)
	require.Equal(t, 2, len(probs))
	require.Equal(t, "Username is required", probs["Username"].Error())
	require.Equal(t, "Email is not a valid email address", probs["Email"].Error())

}
