package is_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/matryer/is"
	"github.com/stretchr/testify/require"
)

type test1 struct {
	Username string  `is:"username,required,nonzero"`
	Code     float64 `is:"code,required"`
	Email    string  `is:"email,required,email"`
}

func TestIsProblems(t *testing.T) {

	data := map[string]interface{}{"username": "", "code": float64(1234), "email": "test@test.com"}
	var target test1
	problems, err := is.NewMSIDecoder(data).Decode(&target)
	require.NoError(t, err)

	require.Equal(t, 1, len(problems))
	require.Equal(t, "cannot be empty", problems["username"].Error())

}

func TestIsSuccess(t *testing.T) {

	data := map[string]interface{}{"username": "mat", "code": float64(1234), "email": "test@test.com"}
	var target test1

	problems, err := is.NewMSIDecoder(data).Decode(&target)
	require.NoError(t, err)

	// make sure values we correctly set
	require.Equal(t, target.Username, "mat")
	require.Equal(t, target.Code, 1234)

	// make sure there were no errors in the problems
	require.Equal(t, 0, len(problems))

}

func TestIsSuccessFromJSON(t *testing.T) {

	var target test1
	problems, err := is.NewJsonDecoder(json.NewDecoder(strings.NewReader(`{"username":"mat","code":1234,"email":"test@test.com"}`))).Decode(&target)
	require.NoError(t, err)

	// make sure values we correctly set
	require.Equal(t, target.Username, "mat")
	require.Equal(t, target.Code, 1234)

	// make sure there were no errors in the problems
	require.Equal(t, 0, len(problems))

}

func TestEmail(t *testing.T) {

	var target struct {
		Email string `is:"email,required,email"`
	}

	data := map[string]interface{}{"email": "nope"}
	probs, err := is.NewMSIDecoder(data).Decode(&target)
	require.NoError(t, err)
	require.Equal(t, 1, len(probs))
	require.Equal(t, "is not a valid email address", probs["email"].Error())

	data = map[string]interface{}{"email": "nope@"}
	probs, err = is.NewMSIDecoder(data).Decode(&target)
	require.NoError(t, err)
	require.Equal(t, 1, len(probs))
	require.Equal(t, "is not a valid email address", probs["email"].Error())

	data = map[string]interface{}{"email": "@nope"}
	probs, err = is.NewMSIDecoder(data).Decode(&target)
	require.NoError(t, err)
	require.Equal(t, 1, len(probs))
	require.Equal(t, "is not a valid email address", probs["email"].Error())

	data = map[string]interface{}{"email": "no.pe@nope"}
	probs, err = is.NewMSIDecoder(data).Decode(&target)
	require.NoError(t, err)
	require.Equal(t, 1, len(probs))
	require.Equal(t, "is not a valid email address", probs["email"].Error())

	data = map[string]interface{}{"email": "nope@me"}
	probs, err = is.NewMSIDecoder(data).Decode(&target)
	require.NoError(t, err)
	require.Equal(t, 1, len(probs))
	require.Equal(t, "is not a valid email address", probs["email"].Error())

	data = map[string]interface{}{"email": "nope@me."}
	probs, err = is.NewMSIDecoder(data).Decode(&target)
	require.NoError(t, err)
	require.Equal(t, 1, len(probs))
	require.Equal(t, "is not a valid email address", probs["email"].Error())

	data = map[string]interface{}{"email": "yes@me.ok"}
	probs, err = is.NewMSIDecoder(data).Decode(&target)
	require.NoError(t, err)
	require.Equal(t, 0, len(probs))

}
