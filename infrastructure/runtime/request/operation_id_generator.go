package request

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

const operationIDRandomBytes = 16

type operationIDGenerator func() (string, error)

/**
 * NewOperationID creates a new logical-operation identity.
 */
func NewOperationID() (
	string,
	error,
) {
	operationID, err :=
		newOperationID(
			rand.Reader,
		)

	if err != nil {
		return "",
			fmt.Errorf(
				"%w: %v",
				ErrOperationIDGeneration,
				err,
			)
	}

	return operationID, nil
}

/**
 * newOperationID creates an Operation-ID from the supplied entropy source.
 */
func newOperationID(
	reader io.Reader,
) (
	string,
	error,
) {
	var raw [operationIDRandomBytes]byte

	if _, err :=
		io.ReadFull(
			reader,
			raw[:],
		); err != nil {

		return "",
			err
	}

	return "op_" +
			hex.EncodeToString(
				raw[:],
			),
		nil
}
