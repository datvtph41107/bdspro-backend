package mapper

import (
	userpb "pb/types/user"
	user_model "user/internal/models/user"
)

type ReportProofMapper struct{}

func NewReportProofMapper() *ReportProofMapper {
	return &ReportProofMapper{}
}

// EntityToPbReportProof chuyển domain ReportProof sang protobuf ReportProof
func (m *ReportProofMapper) EntityToPbReportProof(proof *user_model.ReportProof) *userpb.ReportProof {
	if proof == nil {
		return nil
	}

	pbProof := &userpb.ReportProof{
		Id:       proof.ID,
		ReportId: proof.ReportID,
		FileUrl:  proof.FileUrl,
		FileType: proof.FileType,
	}

	return pbProof
}

// EntityListToPbReportProofs chuyển đổi slice ReportProof domain sang slice protobuf ReportProof
func (m *ReportProofMapper) EntityListToPbReportProofs(proofs []user_model.ReportProof) []*userpb.ReportProof {
	if len(proofs) == 0 {
		return []*userpb.ReportProof{}
	}
	pbProofs := make([]*userpb.ReportProof, len(proofs))
	for i := range proofs {
		pbProofs[i] = m.EntityToPbReportProof(&proofs[i])
	}
	return pbProofs
}
