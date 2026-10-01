package enums

type EAppointmentStatus uint32

const (
	AppointmentStatusPending   EAppointmentStatus = 10 // Chờ xác nhận
	AppointmentStatusConfirmed EAppointmentStatus = 20 // Đã xác nhận
	AppointmentStatusCompleted EAppointmentStatus = 30 // Đã diễn ra
	AppointmentStatusCancelled EAppointmentStatus = 40 // Đã hủy
)

var AppointmentStatusMap = map[EAppointmentStatus]string{
	AppointmentStatusPending:   "Chờ xác nhận",
	AppointmentStatusConfirmed: "Đã xác nhận",
	AppointmentStatusCompleted: "Đã diễn ra",
	AppointmentStatusCancelled: "Đã hủy",
}

func (s EAppointmentStatus) IsValid() bool {
	return s == AppointmentStatusPending ||
		s == AppointmentStatusConfirmed ||
		s == AppointmentStatusCompleted ||
		s == AppointmentStatusCancelled
}