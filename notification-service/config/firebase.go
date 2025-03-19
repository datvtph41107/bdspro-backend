package config

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	firebase "firebase.google.com/go"
	"github.com/spf13/viper"
	"google.golang.org/api/option"
)

// FirebaseConfig lưu trữ thông tin cấu hình Firebase
type FirebaseConfig struct {
	FileURL        string
	ConfigUsername string
	ConfigPassword string
	CredentialFile string
	App            *firebase.App
}

// NewFirebaseConfig khởi tạo FirebaseConfig
func NewFirebaseConfig() (*FirebaseConfig, error) {
	fileURL := viper.GetString("firebase.file")
	credentialFile := viper.GetString("firebase.fileUrl")

	fbConfig := &FirebaseConfig{
		FileURL: fileURL,
		// ConfigUsername: username,
		// ConfigPassword: password,
		CredentialFile: credentialFile,
	}

	err := fbConfig.initFirebase()
	if err != nil {
		return nil, err
	}

	return fbConfig, nil
}

// initFirebase khởi tạo Firebase App
func (fc *FirebaseConfig) initFirebase() error {
	// 1. Lấy nội dung JSON từ Config Server (nếu cần)
	// Nếu lấy từ URL:
	// serviceAccountJSON, err := fc.fetchFirebaseConfigFromServer()
	// Nếu lấy từ file:
	wd, err := os.Getwd()
	serviceAccountJSON, err := os.ReadFile(wd + "/" + fc.CredentialFile)
	if err != nil {
		return fmt.Errorf("failed to read Firebase config file: %w", err)
	}

	// 2. Khởi tạo Firebase
	opt := option.WithCredentialsJSON(serviceAccountJSON)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return fmt.Errorf("failed to initialize Firebase: %w", err)
	}

	fc.App = app
	return nil
}

// fetchFirebaseConfigFromServer lấy file JSON từ Config Server
func (fc *FirebaseConfig) fetchFirebaseConfigFromServer() ([]byte, error) {
	req, err := http.NewRequest("GET", fc.FileURL, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(fc.ConfigUsername, fc.ConfigPassword)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch Firebase config: %s", resp.Status)
	}

	return io.ReadAll(resp.Body)
}

// FirebaseApp trả về Firebase App đã khởi tạo
func (fc *FirebaseConfig) FirebaseApp() *firebase.App {
	return fc.App
}
