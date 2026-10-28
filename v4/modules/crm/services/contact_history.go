package services

import (
	"context"
	"fmt"

	_dto "common/domain/dto"
	_enum "common/domain/enum"
	"common/requestlog"
	_utils "common/utils"
	domain "mdn4/modules/crm/entity"
)

// recordContactHistory writes one contact activity row (read back by GET /v2/crm/contact/{id}/history).
// Failures are logged and never returned so history cannot block the contact change itself.
func (s *ContactUsecase) recordContactHistory(ctx context.Context, contactID uint64, action _enum.EHistory, notes []string) {
	if s.historyWriter == nil || contactID == 0 {
		return
	}
	history := &_dto.HistoryDTO{
		Title:      _enum.HistoryMap[action],
		Note:       notes,
		TargetId:   contactID,
		TargetType: _enum.TargetHistoryContact,
		ActionType: action,
		OwnerOf:    _enum.EOwnerOfMember,
	}
	if profileID := _utils.GetProfileIdWithContext(ctx); profileID > 0 {
		history.OwnerID = &profileID
	}
	if err := s.historyWriter.CreateHistory(ctx, history); err != nil {
		requestlog.FromContext(ctx, nil).WarnContext(ctx, "contact history not recorded",
			"category", "application", "contact_id", contactID, "action_type", int32(action), "error", err.Error())
	}
}

// contactChangeNotes lists the fields that PostgreContact.Update overwrote with a different value.
func contactChangeNotes(before, after *domain.Contact) []string {
	if before == nil || after == nil {
		return nil
	}
	fields := []struct {
		label    string
		old, new string
	}{
		{"Họ tên", before.FullName, after.FullName},
		{"Số điện thoại", before.Phone, after.Phone},
		{"Email", before.Email, after.Email},
		{"Zalo", before.Zalo, after.Zalo},
		{"Công ty", before.Company, after.Company},
		{"Địa chỉ", before.Address, after.Address},
		{"Ghi chú", before.Note, after.Note},
	}
	notes := make([]string, 0, len(fields)+1)
	for _, f := range fields {
		if f.old != f.new {
			notes = append(notes, fmt.Sprintf("%s: %s → %s", f.label, orEmpty(f.old), orEmpty(f.new)))
		}
	}
	if before.Avatar != after.Avatar {
		notes = append(notes, "Cập nhật ảnh đại diện")
	}
	return notes
}

func orEmpty(v string) string {
	if v == "" {
		return "(trống)"
	}
	return v
}
