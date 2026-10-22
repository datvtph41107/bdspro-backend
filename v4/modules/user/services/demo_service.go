package services

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	_errors "common/errors"
	_utils "common/utils"
	bdspropb "pb/types/bdspro"

	"mdn4/data/demo"
	"mdn4/internal/interface/repo"
	bdsprodto "mdn4/modules/bdspro/dto"
	bdsprosvc "mdn4/modules/bdspro/services"
	crmsvc "mdn4/modules/crm/services"
	hubrepo "mdn4/modules/hub/repository"
	"mdn4/modules/user/domain/user"
	usecases "mdn4/modules/user/services/user"
)

// DemoService owns POST /v2/user/demo-mode seed/cleanup for modules/user.
type DemoService interface {
	SetDemoMode(ctx context.Context, enable bool) error
}

// LocalDemoService seeds contacts/products/deals/appointments via in-process
// modules (no TCP to bdspro/crm/tqd). Snapshots live under v4/data/demo.
// @bind: mdn4/modules/user/services.DemoService
type LocalDemoService struct {
	logs         repo.IDemoProductLogRepo
	profile      *usecases.ProfileUsecase
	contacts     crmsvc.ContactService
	products     bdsprosvc.ProductService
	deals        bdsprosvc.DealService
	appointments bdsprosvc.AppointmentService
	wards        hubrepo.LocationV2Repository
	sync         *_utils.SyncUtil
}

func NewLocalDemoService(
	logs repo.IDemoProductLogRepo,
	profile *usecases.ProfileUsecase,
	contacts crmsvc.ContactService,
	products bdsprosvc.ProductService,
	deals bdsprosvc.DealService,
	appointments bdsprosvc.AppointmentService,
	wards hubrepo.LocationV2Repository,
	sync *_utils.SyncUtil,
) *LocalDemoService {
	return &LocalDemoService{
		logs:         logs,
		profile:      profile,
		contacts:     contacts,
		products:     products,
		deals:        deals,
		appointments: appointments,
		wards:        wards,
		sync:         sync,
	}
}

// SetDemoMode seeds or cleans demo data then flips the profile demo flag.
func (s *LocalDemoService) SetDemoMode(ctx context.Context, enable bool) error {
	if enable {
		if err := s.seedDemoProducts(ctx); err != nil {
			return err
		}
	} else {
		if err := s.cleanupDemoProducts(ctx); err != nil {
			return err
		}
	}
	return s.profile.SetDemoMode(ctx, enable)
}

func (s *LocalDemoService) seedDemoProducts(ctx context.Context) error {
	if s.products == nil {
		return _errors.ReturnError(500, "bdspro product service chưa được cấu hình")
	}
	if s.logs == nil {
		return _errors.ReturnError(500, "demo product log repo chưa được cấu hình")
	}

	profileID := _utils.GetProfileIdWithContext(ctx)
	if profileID == 0 {
		return _errors.ReturnError(401, "profileId không hợp lệ")
	}

	existing, err := s.logs.ListProductIDsByProfileID(ctx, profileID)
	if err != nil {
		return err
	}
	if len(existing) > 0 {
		contactIDs, listErr := s.logs.ListContactIDsByProfileID(ctx, profileID)
		if listErr != nil {
			return listErr
		}
		contactProfileIDs := s.resolveContactProfileIDs(ctx, contactIDs)
		if err := s.seedDemoDealsAndAppointments(ctx, profileID, existing, contactProfileIDs); err != nil {
			return err
		}
		s.bumpDemoSyncKeys(ctx, existing)
		return nil
	}

	originIDs, contactIDs, contactProfileIDs, err := s.seedDemoContacts(ctx, profileID)
	if err != nil {
		return err
	}
	if len(contactIDs) > 0 {
		logs := make([]*user.DemoContactLog, 0, len(contactIDs))
		for _, contactID := range contactIDs {
			logs = append(logs, &user.DemoContactLog{
				ProfileID: profileID,
				ContactID: contactID,
			})
		}
		if err := s.logs.CreateContactBatch(ctx, logs); err != nil {
			s.bestEffortDeleteContacts(ctx, contactIDs)
			return err
		}
	}

	products, err := demo.LoadDemoProducts()
	if err != nil {
		return _errors.ReturnError(500, "không đọc được dữ liệu demo: "+err.Error())
	}

	medias, err := demo.LoadDemoMedias()
	if err != nil {
		return _errors.ReturnError(500, "không đọc được dữ liệu demo media: "+err.Error())
	}

	wardIDsByProvince, err := s.loadDemoWardIDsByProvince(ctx, products)
	if err != nil {
		s.bestEffortDeleteContacts(ctx, contactIDs)
		_ = s.logs.DeleteContactsByProfileID(ctx, profileID)
		return err
	}

	// ESourceStatus: 10 Chờ xác minh, 20 Thiết giấy tờ, 30 Đã xác minh số.
	// Shuffle rồi round-robin để mỗi lần bật demo đều có đủ 3 trạng thái.
	sourceStatuses := []uint32{10, 20, 30}
	rand.Shuffle(len(sourceStatuses), func(i, j int) {
		sourceStatuses[i], sourceStatuses[j] = sourceStatuses[j], sourceStatuses[i]
	})

	createdIDs := make([]uint64, 0, len(products))
	for i, product := range products {
		product.OwnerType = 10 // EOwnerOfMember — ProductMe / ownership filter
		product.SourceStatus = sourceStatuses[i%len(sourceStatuses)]
		// sourceTime: ngày hiện tại trừ random 0–100 ngày
		sourceTime := time.Now().AddDate(0, 0, -rand.Intn(101))
		product.SourceTime = _utils.FormatTimeToString(&sourceTime)
		product.MediaItems = demo.MediaItemsForProduct(medias, i)
		if len(originIDs) > 0 {
			originID := originIDs[i%len(originIDs)]
			product.SourceContactId = &originID
		}
		s.assignDemoWard(product, wardIDsByProvince)
		resp, err := s.products.ProductNew(ctx, product)
		if err != nil {
			s.bestEffortDeleteProducts(ctx, createdIDs)
			s.bestEffortDeleteContacts(ctx, contactIDs)
			_ = s.logs.DeleteContactsByProfileID(ctx, profileID)
			return err
		}
		if resp != nil && resp.Id > 0 {
			createdIDs = append(createdIDs, resp.Id)
		}
	}

	logs := make([]*user.DemoProductLog, 0, len(createdIDs))
	for _, productID := range createdIDs {
		logs = append(logs, &user.DemoProductLog{
			ProfileID: profileID,
			ProductID: productID,
		})
	}
	if err := s.logs.CreateBatch(ctx, logs); err != nil {
		s.bestEffortDeleteProducts(ctx, createdIDs)
		s.bestEffortDeleteContacts(ctx, contactIDs)
		_ = s.logs.DeleteContactsByProfileID(ctx, profileID)
		return err
	}

	if err := s.seedDemoDealsAndAppointments(ctx, profileID, createdIDs, contactProfileIDs); err != nil {
		return err
	}
	s.bumpDemoSyncKeys(ctx, createdIDs)
	return nil
}

// seedDemoContacts tạo liên hệ CRM từ snapshot contacts.json; trả originId + contactId + profileId.
func (s *LocalDemoService) seedDemoContacts(ctx context.Context, currentProfileID uint64) ([]uint64, []uint64, []uint64, error) {
	if s.contacts == nil {
		return nil, nil, nil, _errors.ReturnError(500, "crm contact service chưa được cấu hình")
	}

	contacts, err := demo.LoadDemoContacts()
	if err != nil {
		return nil, nil, nil, _errors.ReturnError(500, "không đọc được dữ liệu demo contacts: "+err.Error())
	}

	originIDs := make([]uint64, 0, len(contacts))
	contactIDs := make([]uint64, 0, len(contacts))
	contactProfileIDs := make([]uint64, 0, len(contacts))
	seenContact := make(map[uint64]struct{}, len(contacts))

	for _, contact := range contacts {
		if contact.ProfileID == 0 || contact.ProfileID == currentProfileID {
			continue
		}
		if contact.Phone == "" || contact.FullName == "" {
			continue
		}

		originID, contactID, profileID, err := s.contacts.CreateContact(ctx, contact.FullName, contact.Phone, contact.Email, contact.Avatar, contact.ProfileID)
		if err != nil || contactID == 0 {
			// Đã tồn tại / lỗi tạo — tìm lại theo SĐT thuộc owner hiện tại.
			existingID, findErr := s.contacts.FindContactIDByPhone(ctx, contact.Phone)
			if findErr != nil || existingID == 0 {
				log.Printf("seedDemoContacts skip phone=%s: create=%v find=%v", contact.Phone, err, findErr)
				continue
			}
			contactID = existingID
			if profileID == 0 {
				profileID, _ = s.contacts.GetContactProfileID(ctx, contactID)
			}
			if originID == 0 && profileID > 0 {
				originID, _ = s.contacts.GetOriginID(ctx, profileID)
			}
			if originID == 0 {
				originID, _ = s.contacts.GetOriginID(ctx, contact.ProfileID)
			}
		}
		if originID == 0 && profileID > 0 {
			originID, _ = s.contacts.GetOriginID(ctx, profileID)
		}
		if originID > 0 {
			originIDs = append(originIDs, originID)
		}
		if contactID > 0 {
			if _, ok := seenContact[contactID]; !ok {
				seenContact[contactID] = struct{}{}
				contactIDs = append(contactIDs, contactID)
				if profileID > 0 {
					contactProfileIDs = append(contactProfileIDs, profileID)
				}
			}
		}
	}
	return originIDs, contactIDs, contactProfileIDs, nil
}

// loadDemoWardIDsByProvince lấy ward id từ hub ward_v2 trên DSN v4 (db_v4),
// không gọi TQD gRPC — dùng để random khi seed sản phẩm demo.
func (s *LocalDemoService) loadDemoWardIDsByProvince(
	ctx context.Context,
	products []*bdspropb.ProductSaveRequest,
) (map[uint64][]uint64, error) {
	if s.wards == nil {
		return nil, _errors.ReturnError(500, "hub location repo chưa được cấu hình")
	}

	provinceIDs := make(map[uint64]struct{})
	for _, product := range products {
		provinceID := product.GetProvinceId()
		if provinceID == 0 {
			provinceID = 1
		}
		provinceIDs[provinceID] = struct{}{}
	}

	out := make(map[uint64][]uint64, len(provinceIDs))
	for provinceID := range provinceIDs {
		wardIDs, err := s.wards.ListWardIDsByProvinceID(ctx, provinceID)
		if err != nil {
			return nil, _errors.ReturnError(500, "không lấy được ward từ db_v4: "+err.Error())
		}
		if len(wardIDs) == 0 {
			return nil, _errors.ReturnError(500, fmt.Sprintf("db_v4 không có ward cho provinceId=%d", provinceID))
		}
		out[provinceID] = wardIDs
	}
	return out, nil
}

func (s *LocalDemoService) assignDemoWard(
	product *bdspropb.ProductSaveRequest,
	wardIDsByProvince map[uint64][]uint64,
) {
	if product == nil {
		return
	}
	provinceID := product.GetProvinceId()
	if provinceID == 0 {
		provinceID = 1
		product.ProvinceId = &provinceID
	}
	wardIDs := wardIDsByProvince[provinceID]
	if len(wardIDs) == 0 {
		return
	}
	wardID := wardIDs[rand.Intn(len(wardIDs))]
	product.WardId = &wardID
}

func (s *LocalDemoService) cleanupDemoProducts(ctx context.Context) error {
	if s.products == nil {
		return _errors.ReturnError(500, "bdspro product service chưa được cấu hình")
	}
	if s.logs == nil {
		return _errors.ReturnError(500, "demo product log repo chưa được cấu hình")
	}

	profileID := _utils.GetProfileIdWithContext(ctx)
	if profileID == 0 {
		return _errors.ReturnError(401, "profileId không hợp lệ")
	}

	loggedIDs, err := s.logs.ListProductIDsByProfileID(ctx, profileID)
	if err != nil {
		return err
	}
	loggedSet := make(map[uint64]struct{}, len(loggedIDs))
	for _, id := range loggedIDs {
		if id > 0 {
			loggedSet[id] = struct{}{}
		}
	}

	// Orphan: sản phẩm demo còn sống nhưng không còn (hoặc chưa từng) trong demo_product_log.
	orphanIDs, orphanErr := s.listActiveDemoProductIDs(ctx)
	if orphanErr != nil {
		log.Printf("cleanupDemoProducts list orphan: %v", orphanErr)
	}

	seen := make(map[uint64]struct{}, len(loggedIDs)+len(orphanIDs))
	toDelete := make([]uint64, 0, len(loggedIDs)+len(orphanIDs))
	for _, id := range append(loggedIDs, orphanIDs...) {
		if id == 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		toDelete = append(toDelete, id)
	}

	if err := s.cleanupDemoAppointments(ctx, profileID); err != nil {
		return err
	}
	if err := s.cleanupDemoDeals(ctx, profileID); err != nil {
		return err
	}

	deletedLoggedIDs := make([]uint64, 0, len(loggedIDs))
	var firstErr error
	for _, productID := range toDelete {
		if err := s.products.DeleteProduct(ctx, productID); err != nil {
			log.Printf("cleanupDemoProducts delete id=%d: %v", productID, err)
			if firstErr == nil {
				firstErr = err
			}
			continue
		}
		// Chỉ xoá log khi service con (bdspro) xoá thành công.
		if _, ok := loggedSet[productID]; ok {
			deletedLoggedIDs = append(deletedLoggedIDs, productID)
		}
	}
	if err := s.logs.DeleteByProfileIDAndProductIDs(ctx, profileID, deletedLoggedIDs); err != nil {
		return err
	}

	if err := s.cleanupDemoContacts(ctx, profileID); err != nil {
		return err
	}
	s.bumpDemoSyncKeys(ctx, deletedLoggedIDs)
	if firstErr != nil {
		return _errors.ReturnError(500, "không xoá hết sản phẩm demo: "+firstErr.Error())
	}
	return nil
}

// listActiveDemoProductIDs tìm sản phẩm đang active của user trùng tên snapshot demo.
func (s *LocalDemoService) listActiveDemoProductIDs(ctx context.Context) ([]uint64, error) {
	products, err := demo.LoadDemoProducts()
	if err != nil {
		return nil, err
	}
	demoNames := make(map[string]struct{}, len(products))
	for _, product := range products {
		if product == nil || product.Name == "" {
			continue
		}
		demoNames[product.Name] = struct{}{}
	}
	if len(demoNames) == 0 {
		return nil, nil
	}

	const pageSize uint32 = 100
	var (
		page uint32
		ids  []uint64
	)
	for {
		resp, err := s.products.ProductMe(ctx, bdsprodto.SyncListRequest{Page: page, Size: pageSize})
		if err != nil {
			return nil, err
		}
		if resp == nil || len(resp.Data) == 0 {
			break
		}
		for _, item := range resp.Data {
			if item == nil || item.Id == 0 {
				continue
			}
			if _, ok := demoNames[item.Name]; ok {
				ids = append(ids, item.Id)
			}
		}
		if uint32(len(resp.Data)) < pageSize {
			break
		}
		if resp.TotalElements > 0 && int64((page+1)*pageSize) >= resp.TotalElements {
			break
		}
		page++
		if page > 50 {
			break
		}
	}
	return ids, nil
}

func (s *LocalDemoService) cleanupDemoContacts(ctx context.Context, profileID uint64) error {
	if s.contacts == nil {
		return nil
	}

	contactIDs, err := s.logs.ListContactIDsByProfileID(ctx, profileID)
	if err != nil {
		return err
	}
	loggedSet := make(map[uint64]struct{}, len(contactIDs))
	seen := make(map[uint64]struct{}, len(contactIDs))
	deletedLoggedIDs := make([]uint64, 0, len(contactIDs))

	for _, contactID := range contactIDs {
		if contactID == 0 {
			continue
		}
		loggedSet[contactID] = struct{}{}
		seen[contactID] = struct{}{}
		if err := s.contacts.DeleteContact(ctx, contactID); err != nil {
			log.Printf("cleanupDemoContacts delete id=%d: %v", contactID, err)
			continue
		}
		deletedLoggedIDs = append(deletedLoggedIDs, contactID)
	}

	// Best-effort: xoá thêm contact demo cũ chưa kịp log (theo SĐT trong contacts.json).
	contacts, loadErr := demo.LoadDemoContacts()
	if loadErr == nil {
		for _, contact := range contacts {
			if contact.Phone == "" {
				continue
			}
			existingID, findErr := s.contacts.FindContactIDByPhone(ctx, contact.Phone)
			if findErr != nil || existingID == 0 {
				continue
			}
			if _, ok := seen[existingID]; ok {
				continue
			}
			seen[existingID] = struct{}{}
			if err := s.contacts.DeleteContact(ctx, existingID); err != nil {
				log.Printf("cleanupDemoContacts delete phone=%s id=%d: %v", contact.Phone, existingID, err)
				continue
			}
			if _, ok := loggedSet[existingID]; ok {
				deletedLoggedIDs = append(deletedLoggedIDs, existingID)
			}
		}
	}

	return s.logs.DeleteContactsByProfileIDAndContactIDs(ctx, profileID, deletedLoggedIDs)
}

func (s *LocalDemoService) bumpDemoSyncKeys(ctx context.Context, productIDs []uint64) {
	if s.sync == nil {
		return
	}
	now := time.Now().UnixMilli()
	_ = s.sync.PutTimestamp(ctx, s.sync.GetKey(ctx, _utils.SyncKeyContactMe, 0), now)
	_ = s.sync.PutTimestamp(ctx, s.sync.GetKey(ctx, _utils.SyncKeyProductMe, 0), now)
	_ = s.sync.PutTimestamp(ctx, s.sync.GetKey(ctx, _utils.SyncKeyPropertyMe, 0), now)
	_ = s.sync.PutTimestamp(ctx, s.sync.GetKey(ctx, _utils.SyncKeyAppointmentMe, 0), now)
	_ = s.sync.PutTimestamp(ctx, s.sync.GetKey(ctx, _utils.SyncKeyDealMe, 0), now)
	seen := make(map[uint64]struct{}, len(productIDs))
	for _, productID := range productIDs {
		if productID == 0 {
			continue
		}
		if _, ok := seen[productID]; ok {
			continue
		}
		seen[productID] = struct{}{}
		_ = s.sync.PutTimestamp(ctx, s.sync.GetKey(ctx, _utils.SyncKeyProductDeals, productID), now)
		_ = s.sync.PutTimestamp(ctx, s.sync.GetKey(ctx, _utils.SyncKeyProductAppointments, productID), now)
		_ = s.sync.PutTimestamp(ctx, s.sync.GetKey(ctx, _utils.SyncKeyProductContacts, productID), now)
	}
}

func (s *LocalDemoService) resolveContactProfileIDs(ctx context.Context, contactIDs []uint64) []uint64 {
	if s.contacts == nil || len(contactIDs) == 0 {
		return nil
	}
	out := make([]uint64, 0, len(contactIDs))
	seen := make(map[uint64]struct{}, len(contactIDs))
	for _, contactID := range contactIDs {
		profileID, err := s.contacts.GetContactProfileID(ctx, contactID)
		if err != nil || profileID == 0 {
			continue
		}
		if _, ok := seen[profileID]; ok {
			continue
		}
		seen[profileID] = struct{}{}
		out = append(out, profileID)
	}
	return out
}

func (s *LocalDemoService) seedDemoDealsAndAppointments(
	ctx context.Context,
	profileID uint64,
	productIDs []uint64,
	contactProfileIDs []uint64,
) error {
	dealIDs, err := s.seedDemoDeals(ctx, profileID, productIDs, contactProfileIDs)
	if err != nil {
		return err
	}
	if err := s.seedDemoAppointments(ctx, profileID, productIDs, contactProfileIDs, dealIDs); err != nil {
		s.bestEffortDeleteDeals(ctx, dealIDs)
		_ = s.logs.DeleteDealsByProfileID(ctx, profileID)
		return err
	}
	return nil
}

func (s *LocalDemoService) seedDemoDeals(
	ctx context.Context,
	profileID uint64,
	productIDs []uint64,
	contactProfileIDs []uint64,
) ([]uint64, error) {
	if s.deals == nil {
		return nil, _errors.ReturnError(500, "bdspro deal service chưa được cấu hình")
	}

	existing, err := s.logs.ListDealIDsByProfileID(ctx, profileID)
	if err != nil {
		return nil, err
	}
	if len(existing) > 0 {
		for i, dealID := range existing {
			var customerID uint64
			if len(contactProfileIDs) > 0 {
				customerID = contactProfileIDs[i%len(contactProfileIDs)]
			}
			dealType := uint32(40)
			if i%2 == 1 {
				dealType = 50
			}
			if err := s.deals.SeedDemoRelated(ctx, dealID, profileID, customerID, dealType); err != nil {
				return nil, err
			}
		}
		return existing, nil
	}

	dealNames := []string{"Thương vụ demo bán căn hộ", "Thương vụ demo cho thuê nhà phố"}
	dealTypes := []uint32{40, 50} // Sale, Rent
	createdIDs := make([]uint64, 0, 2)

	for i := 0; i < 2; i++ {
		var customers []uint64
		if len(contactProfileIDs) > 0 {
			customers = []uint64{contactProfileIDs[i%len(contactProfileIDs)]}
		}
		var products []uint64
		if len(productIDs) > 0 {
			products = []uint64{productIDs[i%len(productIDs)]}
		}
		var customerID uint64
		if len(customers) > 0 {
			customerID = customers[0]
		}
		req := &bdspropb.CreateGroupDealRequest{
			Name:           dealNames[i],
			Status:         10, // Processing
			DealType:       dealTypes[i%len(dealTypes)],
			TargetProfit:   float64(500000000 + rand.Intn(500000000)),
			Note:           "Dữ liệu demo — sẽ xoá khi tắt chế độ demo",
			Description:    "Thương vụ được tạo tự động khi bật demo",
			CustomerIds:    customers,
			ProductIds:     products,
			OwnerType:      10, // EOwnerOfMember
			ChargePersonId: profileID,
			OwnerId:        profileID,
		}
		dealID, err := s.deals.CreateGroupDeal(ctx, req)
		if err != nil {
			s.bestEffortDeleteDeals(ctx, createdIDs)
			return nil, err
		}
		if dealID > 0 {
			if err := s.deals.SeedDemoRelated(ctx, dealID, profileID, customerID, dealTypes[i%len(dealTypes)]); err != nil {
				s.bestEffortDeleteDeals(ctx, append(createdIDs, dealID))
				return nil, err
			}
			createdIDs = append(createdIDs, dealID)
		}
	}

	if len(createdIDs) == 0 {
		return nil, nil
	}
	logs := make([]*user.DemoDealLog, 0, len(createdIDs))
	for _, dealID := range createdIDs {
		logs = append(logs, &user.DemoDealLog{ProfileID: profileID, DealID: dealID})
	}
	if err := s.logs.CreateDealBatch(ctx, logs); err != nil {
		s.bestEffortDeleteDeals(ctx, createdIDs)
		return nil, err
	}
	return createdIDs, nil
}

func (s *LocalDemoService) seedDemoAppointments(
	ctx context.Context,
	profileID uint64,
	productIDs []uint64,
	contactProfileIDs []uint64,
	dealIDs []uint64,
) error {
	if s.appointments == nil {
		return _errors.ReturnError(500, "bdspro appointment service chưa được cấu hình")
	}

	existing, err := s.logs.ListAppointmentIDsByProfileID(ctx, profileID)
	if err != nil {
		return err
	}
	if len(existing) > 0 {
		return nil
	}

	titles := []string{
		"Xem nhà demo cùng khách",
		"Ký gửi sản phẩm demo",
		"Tư vấn đầu tư demo",
	}
	places := []string{
		"Văn phòng giao dịch Quận 1",
		"Căn hộ mẫu Thuận An",
		"Nhà phố Thủ Đức",
	}
	statuses := []uint32{10, 20, 30} // pending, confirmed, completed
	createdIDs := make([]uint64, 0, 3)

	for i := 0; i < 3; i++ {
		start := time.Now().Add(time.Duration(24*(i+1)+rand.Intn(48)) * time.Hour)
		start = time.Date(start.Year(), start.Month(), start.Day(), 9+rand.Intn(8), 0, 0, 0, start.Location())
		end := start.Add(time.Duration(1+rand.Intn(2)) * time.Hour)

		var participants []int64
		if len(contactProfileIDs) > 0 {
			participants = []int64{int64(contactProfileIDs[i%len(contactProfileIDs)])}
		}
		var products []int64
		if len(productIDs) > 0 {
			products = []int64{int64(productIDs[i%len(productIDs)])}
		}
		var dealID uint32
		if len(dealIDs) > 0 {
			dealID = uint32(dealIDs[i%len(dealIDs)])
		}

		appointmentID, err := s.appointments.CreateAppointment(
			ctx,
			titles[i],
			"Lịch hẹn demo — sẽ xoá khi tắt chế độ demo",
			places[i],
			start.Format(time.RFC3339),
			end.Format(time.RFC3339),
			participants,
			dealID,
			products,
			statuses[i%len(statuses)],
		)
		if err != nil {
			s.bestEffortDeleteAppointments(ctx, createdIDs)
			return err
		}
		if appointmentID > 0 {
			createdIDs = append(createdIDs, appointmentID)
		}
	}

	if len(createdIDs) == 0 {
		return nil
	}
	logs := make([]*user.DemoAppointmentLog, 0, len(createdIDs))
	for _, appointmentID := range createdIDs {
		logs = append(logs, &user.DemoAppointmentLog{ProfileID: profileID, AppointmentID: appointmentID})
	}
	if err := s.logs.CreateAppointmentBatch(ctx, logs); err != nil {
		s.bestEffortDeleteAppointments(ctx, createdIDs)
		return err
	}
	return nil
}

func (s *LocalDemoService) cleanupDemoAppointments(ctx context.Context, profileID uint64) error {
	if s.appointments == nil {
		return nil
	}
	appointmentIDs, err := s.logs.ListAppointmentIDsByProfileID(ctx, profileID)
	if err != nil {
		return err
	}
	deleted := make([]uint64, 0, len(appointmentIDs))
	for _, appointmentID := range appointmentIDs {
		if appointmentID == 0 {
			continue
		}
		if err := s.appointments.DeleteAppointment(ctx, appointmentID); err != nil {
			log.Printf("cleanupDemoAppointments delete id=%d: %v", appointmentID, err)
			continue
		}
		deleted = append(deleted, appointmentID)
	}
	return s.logs.DeleteAppointmentsByProfileIDAndIDs(ctx, profileID, deleted)
}

func (s *LocalDemoService) cleanupDemoDeals(ctx context.Context, profileID uint64) error {
	if s.deals == nil {
		return nil
	}
	dealIDs, err := s.logs.ListDealIDsByProfileID(ctx, profileID)
	if err != nil {
		return err
	}
	deleted := make([]uint64, 0, len(dealIDs))
	for _, dealID := range dealIDs {
		if dealID == 0 {
			continue
		}
		if err := s.deals.DeleteGroupDeal(ctx, dealID); err != nil {
			log.Printf("cleanupDemoDeals delete id=%d: %v", dealID, err)
			continue
		}
		deleted = append(deleted, dealID)
	}
	return s.logs.DeleteDealsByProfileIDAndIDs(ctx, profileID, deleted)
}

func (s *LocalDemoService) bestEffortDeleteAppointments(ctx context.Context, appointmentIDs []uint64) {
	if s.appointments == nil {
		return
	}
	for _, appointmentID := range appointmentIDs {
		_ = s.appointments.DeleteAppointment(ctx, appointmentID)
	}
}

func (s *LocalDemoService) bestEffortDeleteDeals(ctx context.Context, dealIDs []uint64) {
	if s.deals == nil {
		return
	}
	for _, dealID := range dealIDs {
		_ = s.deals.DeleteGroupDeal(ctx, dealID)
	}
}

func (s *LocalDemoService) bestEffortDeleteProducts(ctx context.Context, productIDs []uint64) {
	for _, productID := range productIDs {
		_ = s.products.DeleteProduct(ctx, productID)
	}
}

func (s *LocalDemoService) bestEffortDeleteContacts(ctx context.Context, contactIDs []uint64) {
	if s.contacts == nil {
		return
	}
	for _, contactID := range contactIDs {
		_ = s.contacts.DeleteContact(ctx, contactID)
	}
}
