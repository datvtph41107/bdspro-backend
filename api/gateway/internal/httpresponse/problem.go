package httpresponse

import (
	"context"
	"net/http"

	sharedhttp "common/httpresponse"
)

const ProblemMediaType = sharedhttp.ProblemMediaType

type FieldProblem = sharedhttp.FieldProblem
type Problem = sharedhttp.Problem

func NewProblem(
	status int,
	code string,
	detail string,
) Problem {
	return sharedhttp.NewProblem(
		status,
		code,
		detail,
	)
}

func TypeForCode(code string) string {
	return sharedhttp.TypeForCode(code)
}

func ProblemFromError(err error) Problem {
	return sharedhttp.ProblemFromError(err)
}

func WriteProblem(
	ctx context.Context,
	w http.ResponseWriter,
	problem Problem,
) {
	sharedhttp.WriteProblem(
		ctx,
		w,
		problem,
	)
}
