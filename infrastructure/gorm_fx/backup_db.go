package gorm_fx

import (
	"errors"
	"time"

	"github.com/zsmartex/pkg/v2/config"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

var BackupModule = fx.Module("backup_db",
	fx.Provide(NewBackupDatabase),
)

type DatabaseBackup struct {
	*gorm.DB
}

func NewBackupDatabase(cfg config.BackupPostgres) (*DatabaseBackup, error) {
	db, err := New(Config{
		Config: config.Postgres{
			Host:            cfg.Host,
			Port:            cfg.Port,
			User:            cfg.User,
			Pass:            cfg.Pass,
			Name:            cfg.Name,
			ApplicationName: cfg.ApplicationName,
			SSLMode:         cfg.SSLMode,
		},
	})
	if err != nil {
		return nil, err
	}

	return &DatabaseBackup{DB: db}, nil
}

type PGStatWalReceiverRes struct {
	PID                int
	Status             string
	ReceiveStartLSN    int
	ReceiveStartTLI    int
	WrittenLSN         int
	FlushedLSN         int
	ReceivedTLI        int
	LastMsgSendTime    time.Time
	LastMsgReceiptTime time.Time
	LatestEndLSN       int
	LatestEndTime      time.Time
	SlotName           string
	SenderHost         string
	SenderPort         int
	ConnInfo           string
}

func (db *DatabaseBackup) PGStatWalReceiver() (*PGStatWalReceiverRes, error) {
	var res []*PGStatWalReceiverRes

	err := db.Raw(`SELECT * FROM pg_stat_wal_receiver;`).Scan(&res).Error
	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, errors.New("pg_stat_wal_receiver is empty")
	}

	for _, r := range res {
		if r.Status == "streaming" {
			return r, nil
		}
	}

	return nil, errors.New("slot not found")
}
