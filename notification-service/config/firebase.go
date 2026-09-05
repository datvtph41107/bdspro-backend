package config

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	firebase "firebase.google.com/go"
	"github.com/spf13/viper"
	"google.golang.org/api/option"
)

// FirebaseConfig là capability tùy chọn tại process boundary. Thiếu credential
// sẽ tắt push delivery nhưng không làm Notification history/inbox fail. File
// service-account thật phải nằm ngoài Git và được inject bằng env.
type FirebaseConfig struct {
	CredentialFile string
	App            *firebase.App
}

func NewFirebaseConfig() (*FirebaseConfig, error) {
	credentialFile := strings.TrimSpace(os.Getenv("NOTIFICATION_FIREBASE_CREDENTIAL_FILE"))
	if credentialFile == "" {
		credentialFile = strings.TrimSpace(viper.GetString("firebase.fileUrl"))
	}
	cfg := &FirebaseConfig{CredentialFile: credentialFile}
	if credentialFile == "" {
		return cfg, nil
	}
	if !filepath.IsAbs(credentialFile) {
		wd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		credentialFile = filepath.Join(wd, credentialFile)
	}
	serviceAccountJSON, err := os.ReadFile(credentialFile)
	if os.IsNotExist(err) {
		return cfg, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read Firebase credential: %w", err)
	}
	app, err := firebase.NewApp(context.Background(), nil, option.WithCredentialsJSON(serviceAccountJSON))
	if err != nil {
		return nil, fmt.Errorf("initialize Firebase: %w", err)
	}
	cfg.App = app
	return cfg, nil
}

func (fc *FirebaseConfig) FirebaseApp() *firebase.App {
	if fc == nil {
		return nil
	}
	return fc.App
}
