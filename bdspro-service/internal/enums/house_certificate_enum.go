package enums

type EHouseCertificate uint32

const (
	EHouseCertificateNone            EHouseCertificate = 0  // Chưa xác định
	EHouseCertificateRedBook         EHouseCertificate = 10 // Sổ đỏ
	EHouseCertificatePinkBook        EHouseCertificate = 20 // Sổ hồng
	EHouseCertificateNotarizedRecord EHouseCertificate = 30 // Vi bằng
)

var EHouseCertificateNames = map[EHouseCertificate]string{
	EHouseCertificateNone:            "Chưa xác định",
	EHouseCertificateRedBook:         "Sổ đỏ",
	EHouseCertificatePinkBook:        "Sổ hồng",
	EHouseCertificateNotarizedRecord: "Vi bằng",
}
