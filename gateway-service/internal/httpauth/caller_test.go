package httpauth

import (
	"common/identity"
	"common/jwtverify"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClassifyCaller_TokenTypes(t *testing.T) {
	t.Parallel()
	base := jwtverify.Principal{AuthID: 1, ProfileID: 2, SessionID: 3}

	access := base
	access.Type = string(jwtverify.AccessToken)
	caller, err := ClassifyCaller(CallerInput{Principal: &access})
	require.NoError(t, err)
	assert.Equal(t, identity.CallerUser, caller.Kind)

	temp := base
	temp.Type = string(jwtverify.TempToken)
	_, err = ClassifyCaller(CallerInput{Principal: &temp})
	require.NoError(t, err)

	refresh := base
	refresh.Type = string(jwtverify.RefreshToken)
	_, err = ClassifyCaller(CallerInput{Principal: &refresh})
	assert.ErrorIs(t, err, ErrTokenTypeNotAllowed)

	unknown := base
	unknown.Type = "OTHER"
	_, err = ClassifyCaller(CallerInput{Principal: &unknown})
	assert.ErrorIs(t, err, ErrTokenTypeNotAllowed)
}
