package filehttp

import (
	"net/http"
	"time"

	_jwt "common/jwt"

	"file/internal/fileauthorization"

	"github.com/gin-gonic/gin"
)

type signedReadTokenIssuer interface {
	Issue(
		profileID int64,
		now time.Time,
	) fileauthorization.IssuedSignedReadToken
}

// SignatureHandler owns the HTTP boundary for issuing signed-read tokens.
//
// Token construction itself remains transport-independent in
// fileauthorization.SignedReadIssuer.
type SignatureHandler struct {
	issuer signedReadTokenIssuer
}

func NewSignatureHandler(
	issuer signedReadTokenIssuer,
) *SignatureHandler {
	return &SignatureHandler{
		issuer: issuer,
	}
}

type signatureResponse struct {
	S string
	I int64
	T int64
}

// IssueSignedReadToken preserves the existing /v1/file/signature HTTP
// response contract while moving token construction to its canonical owner.
//
// @Summary API sinh chữ ký để xem file được bảo vệ
// @Description
// @Security BearerAuth
// @Produce json
// @Router /signature [put]
func (h *SignatureHandler) IssueSignedReadToken(
	c *gin.Context,
) {
	profileID := int64(_jwt.GetProfileId(c))

	issued := h.issuer.Issue(
		profileID,
		time.Now(),
	)

	c.JSON(
		http.StatusOK,
		signatureResponse{
			S: issued.Encoded,
			I: issued.ProfileID,
			T: issued.ExpiresAt,
		},
	)
}
