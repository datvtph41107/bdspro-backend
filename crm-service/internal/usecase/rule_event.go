package usecase

import (
	_enum "common/domain/enum"
	"context"
	"crm/internal/dto"
	"crm/internal/enums"
	"crm/internal/interface/provider"
	"crm/internal/repo"
	"fmt"
	"strconv"
	"time"
)

type RuleEventUsecase struct {
	ruleRepo             repo.RuleRepo
	customerRepo         repo.LeadRepo
	notificationProvider provider.NotificationProvider
}

func NewRuleEventUsecase(ruleRepo repo.RuleRepo,
	customerRepo repo.LeadRepo,
	client provider.NotificationProvider) *RuleEventUsecase {
	return &RuleEventUsecase{ruleRepo: ruleRepo,
		customerRepo:         customerRepo,
		notificationProvider: client,
	}
}

func (u *RuleEventUsecase) TriggerJobAutoEvent(c context.Context) error {
	fmt.Println("Trigger Rules", time.Now().Format("2006-01-02 15:04:05"))
	customers, err := u.customerRepo.GetOwnerWithRule(c)
	if err != nil {
		return err
	}
	for _, customer := range customers {
		if customer.TriggerToDay {
			continue
		}
		switch customer.Condition {
		case enums.ConditionDateNotHistory:
			u.CheckEventNotHistory(c, customer)
		case enums.ConditionStage:
			u.CheckEventStage(c, customer)
		case enums.ConditionTag:
			u.CheckEventTag(c, customer)
		case enums.ConditionBirth:
			u.CheckEventBirth(c, customer)
		case enums.ConditionEventDate:
			u.CheckEventEventDate(c, customer)
		}
	}

	return nil
}

func (u *RuleEventUsecase) CheckEventNotHistory(c context.Context, customer dto.CustomerWithRuleDTO) error {
	rangeDate, err := strconv.Atoi(customer.ConditionValue)
	if err != nil {
		return err
	}

	// thời điểm $rangeDate ngày trước đó
	fromDate := time.Now().AddDate(0, 0, -rangeDate)

	// nếu khách hàng mới tạo và chưa quá $rangeDate ngày trước đó thì không trigger
	if customer.UpdatedAt.After(fromDate) {
		return nil
	}

	// existed, err := u.notificationProvider.ExistedFromDate(c, customer.CustomerID, fromDate)
	// if err != nil {
	// 	return err
	// }
	// if existed {
	// 	return nil
	// }

	u.TriggerEvent(c, customer)

	return nil
}

func (u *RuleEventUsecase) CheckEventStage(c context.Context, customer dto.CustomerWithRuleDTO) error {
	stageID, err := strconv.ParseUint(customer.TriggerValue, 10, 64)
	if err != nil {
		return nil
	}

	if customer.StageID != nil && *customer.StageID == stageID {
		u.TriggerEvent(c, customer)
	}

	return nil
}

func (u *RuleEventUsecase) CheckEventTag(c context.Context, customer dto.CustomerWithRuleDTO) error {
	// u.TriggerEvent(c, customer)
	return nil
}

func (u *RuleEventUsecase) CheckEventBirth(c context.Context, customer dto.CustomerWithRuleDTO) error {
	birthday := customer.Birthday
	if birthday == nil {
		return nil
	}

	now := time.Now()
	birthdayDate := birthday.Day()
	birthdayMonth := birthday.Month()

	chargeId, err := strconv.ParseUint(customer.TriggerValue, 10, 64)
	if err != nil {
		return err
	}
	if now.Day() == birthdayDate && now.Month() == birthdayMonth && customer.ChargePersonId != chargeId {
		u.TriggerEvent(c, customer)
	}

	return nil
}

func (u *RuleEventUsecase) CheckEventEventDate(c context.Context, customer dto.CustomerWithRuleDTO) error {
	eventDate, err := time.Parse(time.RFC3339, customer.TriggerValue)
	if err != nil {
		return nil
	}

	now := time.Now()
	if now.Day() == eventDate.Day() && now.Month() == eventDate.Month() && now.Year() == eventDate.Year() {
		u.TriggerEvent(c, customer)
	}

	return nil
}

// ------------------------------------------------------------

func (u *RuleEventUsecase) TriggerEvent(c context.Context, customer dto.CustomerWithRuleDTO) error {
	fmt.Println("Event Auto Interact Customer Trigger:", customer.CustomerID, customer.FullName, customer.RuleName)

	switch customer.Trigger {
	case enums.RuleThenReminder:
		u.TriggerReminder(c, customer)
	case enums.RuleThenLabel:
		u.TriggerLabel(c, customer)
	case enums.RuleThenAssignLead:
		u.TriggerAssignLead(c, customer)
	case enums.RuleThenSwitchStage:
		u.TriggerSwitchStage(c, customer)
	}

	// u.notificationProvider.CreateCRMHistory(c, &base_dto.HistoryDTO{
	// 	TargetId:   customer.CustomerID,
	// 	TargetType: base_enum.TargetHistoryLead,
	// 	ActionType: base_enum.HistoryRuleEvent,
	// 	OwnerID:    &customer.OwnerID,
	// 	OwnerOf:    base_enum.EOwnerOf(customer.OwnerType),
	// 	Title:      "Tương tác tự động với khách hàng thành công",
	// 	Note: []string{"Tương tác tự động với khách hàng",
	// 		customer.FullName,
	// 		"với quy tắc",
	// 		customer.RuleName,
	// 		"Đã thực hiện",
	// 		enums.RuleThenMap[customer.Trigger],
	// 	},
	// })
	// u.notificationProvider.CreateCRMHistory(c, &base_dto.HistoryDTO{
	// 	TargetId:   customer.RuleID,
	// 	TargetType: base_enum.TargetHistoryRule,
	// 	ActionType: base_enum.HistoryRuleEvent,
	// 	OwnerID:    &customer.OwnerID,
	// 	OwnerOf:    base_enum.EOwnerOf(customer.OwnerType),
	// 	Title:      "Quy tắc \"" + customer.RuleName + "\" đã kích hoạt cho \"" + customer.FullName + "\"",
	// 	Note:       []string{},
	// })
	return nil
}

func (u *RuleEventUsecase) TriggerReminder(c context.Context, customer dto.CustomerWithRuleDTO) error {
	// TODO: gửi thông báo nhắc nhở
	u.notificationProvider.CreateNotification(c,
		"Nhắc nhở",
		"",
		[]string{fmt.Sprintf("Khách hàng %s đã quá %s ngày chưa tương tác", customer.FullName, customer.ConditionValue)},
		_enum.NotificationRuleEvent,
		&customer.ChargePersonId,
		customer.OwnerID,
		_enum.EOwnerOfMember,
		[]string{},
	)
	return nil
}

func (u *RuleEventUsecase) TriggerLabel(c context.Context, customer dto.CustomerWithRuleDTO) error {
	return nil
}

func (u *RuleEventUsecase) TriggerAssignLead(c context.Context, customer dto.CustomerWithRuleDTO) error {
	// TODO: gán nhân viên
	chargeId, err := strconv.ParseUint(customer.TriggerValue, 10, 64)
	if err != nil {
		return err
	}

	_, err = u.customerRepo.Assign(c, customer.CustomerID, &dto.LeadDTO{
		OwnerID:        &customer.OwnerID,
		OwnerType:      customer.OwnerType,
		ChargePersonId: &chargeId,
		AssignNote:     fmt.Sprintf("Tương tác tự động gán cho thành viên phụ trách %s", "."),
	})
	return err
}

func (u *RuleEventUsecase) TriggerSwitchStage(c context.Context, customer dto.CustomerWithRuleDTO) error {
	// TODO: chuyển trạng thái
	stageID, err := strconv.ParseUint(customer.TriggerValue, 10, 64)
	if err != nil {
		return err
	}
	_, err = u.customerRepo.SwitchStage(c, customer.CustomerID, &stageID, "Tương tác hệ thống")
	return err
}
