package requestcontext

import "net/http"

type operationIDGenerator func() (string, error)

// resolveOperationID is a temporary HTTP-compatibility seam.
// Non-missing input delegates to the canonical request package through
// ResolveOperationID. The injected generator is retained only so legacy HTTP
// atomicity tests can force generation failure until Gateway owns this boundary.
func resolveOperationID(
	values []string,
	generate operationIDGenerator,
) (OperationIDDecision, error) {
	if len(values) != 0 {
		return ResolveOperationID(values)
	}

	operationID, err := generate()
	if err != nil {
		return OperationIDDecision{}, err
	}
	if !IsValidOperationID(operationID) {
		return OperationIDDecision{}, ErrOperationIDGeneration
	}

	return OperationIDDecision{
		ID:     operationID,
		Source: OperationIDSourceGeneratedFallback,
		Reason: OperationIDReasonMissing,
	}, nil
}

/**
 * PrepareHTTPOperationID resolves and atomically binds the canonical
 * logical-operation identity at the HTTP ingress boundary.
 */
func PrepareHTTPOperationID(
	writer http.ResponseWriter,
	request *http.Request,
) (
	OperationIDDecision,
	error,
) {
	return prepareHTTPOperationID(
		writer,
		request,
		NewOperationID,
	)
}

func prepareHTTPOperationID(
	writer http.ResponseWriter,
	request *http.Request,
	generate operationIDGenerator,
) (
	OperationIDDecision,
	error,
) {
	if request == nil {
		panic(
			"requestcontext: nil HTTP request",
		)
	}

	if writer == nil {
		panic(
			"requestcontext: nil HTTP response writer",
		)
	}

	decision, err :=
		resolveOperationID(
			request.Header.Values(
				OperationIDHeader,
			),
			generate,
		)

	if err != nil {
		return OperationIDDecision{},
			err
	}

	boundContext, err :=
		BindOperationID(
			request.Context(),
			decision.ID,
		)

	if err != nil {
		return OperationIDDecision{},
			err
	}

	preparedRequest :=
		request.WithContext(
			boundContext,
		)

	preparedRequest.Header =
		request.Header.Clone()

	if preparedRequest.Header == nil {

		preparedRequest.Header = make(
			http.Header,
		)
	}

	preparedRequest.Header.Del(
		OperationIDHeader,
	)

	preparedRequest.Header.Set(
		OperationIDHeader,
		decision.ID,
	)

	// All semantic failure points have passed.
	// From here onward we commit the prepared boundary state.
	*request = *preparedRequest

	writer.Header().Set(
		OperationIDHeader,
		decision.ID,
	)

	return decision, nil
}
