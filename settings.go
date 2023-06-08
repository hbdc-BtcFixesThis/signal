package main

import (
	"fmt"

	"encoding/json"

	bolt "go.etcd.io/bbolt"
)

var (
	defaultPasswordLen  = 20
	settingsBucketName  = []byte("settings")
	settingsDBRecordKey = []byte("settings")
)

type SignalSettings struct {
	Port                string   `json:"port"`
	UiDir               string   `json:"ui_dir"`
	IsPublicNode        bool     `json:"is_public"`
	TlsCrtPath          string   `json:"crt"`
	TlsKeyPath          string   `json:"key"`
	MaxStorageSizeBytes uint64   `json:"max_storage_size"`
	DbPath              string   `json:"signal_db_path"`
	PassHash            string   `json:"pass"`
	Peers               []string `json:"peers"`
}

func LoadSettings(db *bolt.DB) (*SignalSettings, error) {
	var settings SignalSettings
	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(settingsBucketName)
		if err != nil {
			return err
		}

		settingsBytes := b.Get(settingsDBRecordKey)
		if err != nil {
			return err
		}

		if settingsBytes == nil {
			defaultSettings, err := defaultSettingsBytes()
			if err != nil {
				return err
			}
			b.Put(settingsDBRecordKey, defaultSettings)
			// this is a bit redundant since the settings
			// were just saved but this only happens when
			// settings are set for the first time. Since
			// this method will be accessed infrequently
			// retrieving the settings after saving could
			// help detect issues dispite the extra lookup
			settingsBytes = b.Get(settingsDBRecordKey)
		}
		err = json.Unmarshal(settingsBytes, &settings)
		return err
	})

	if err != nil {
		return nil, err
	}

	return &settings, err
}

func defaultSettingsBytes() ([]byte, error) {
	// only called when new settings conf is being generated
	pw, err := GenRandStr(defaultPasswordLen)
	fmt.Printf(`

		------ No config found! ------

		A new one has been generated for you. To access and update your settings,
		including the password gennerated below, visit https://localhost:8888 on
		your prefered browser and use the following pw:

		########################################################
		####                                                ####
		####        password: %s%s
		####                                                ####
		########################################################

		In the event that you would like to reset your settings to their defaults
		you can do so by deleting the file signal_conf.db and restarting the
		process again. This will regenerate a conf file as well as a new password.


		NOTE: The password above is stored as a hash. If you forget the password
		above, or any future password/s set, you WILL NOT be able to retrieve them
		here!


`, pw, "          ####")
	if err != nil {
		return nil, err
	}

	// 1gb default max size
	return []byte(`{
		"port": ":8888",
		"ui_dir": "static",
		"is_public": true,
		"crt": "signal.crt.pem",
		"key": "signal.key.pem",
		"max_storage_size": 1000000000,
		"signal_db_path": "signal.db",
		"pass": "` + SHA256(pw) + `",
		"peers": []
	}`), nil
}
