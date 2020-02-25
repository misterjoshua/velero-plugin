package snapshot

import (
	"github.com/libopenstorage/openstorage/api"
	"github.com/sirupsen/logrus"
	"time"
)

// BackupInfoCache is a cache of backup info.
type BackupInfoCache struct {
	backups       []api.CloudBackupInfo
	lastRefreshed time.Time
	credID        string
}

var backups map[string]*BackupInfoCache = make(map[string]*BackupInfoCache)

// GetBackupInfoCacheByCredID gets a the right cache for a credential id
func GetBackupInfoCacheByCredID(credID string) *BackupInfoCache {
	if val, ok := backups[credID]; ok {
		return val
	}

	backups[credID] = &BackupInfoCache{
		credID: credID,
	}
	return backups[credID]
}

// GetCachedBackupInfo gets all backup info from cache first.
func (c *BackupInfoCache) GetCachedBackupInfo(logger logrus.FieldLogger) ([]api.CloudBackupInfo, error) {
	// Check how long it's been since we updated the cache.
	// Note: If the cache has never been set, lastRefreshed will be nil
	//   which is Jan 1st of year 1 (2000+ years ago), so we can check
	//   minutes here and it'll surely be larger than five minutes.
	durationSinceLastRefreshed := time.Now().Sub(c.lastRefreshed)
	minutesSinceLastRefresh := durationSinceLastRefreshed.Minutes()
	logger.Infof("Minutes since last refresh: %v", minutesSinceLastRefresh)

	if minutesSinceLastRefresh > 5 {
		logger.Infof("Refreshing backup info cache")
		backups, err := enumerateBackups(logger, c.credID)
		if err != nil {
			return nil, err
		}

		logger.Infof("Setting new cache")
		c.backups = backups
		c.lastRefreshed = time.Now()
	} else {
		logger.Info("Using backup info cache")
	}

  logger.Infof("There are %v backups in the cache", len(c.backups))
	return c.backups, nil
}

func enumerateBackups(logger logrus.FieldLogger, credID string) ([]api.CloudBackupInfo, error) {
	volDriver, err := getVolumeDriver()
	if err != nil {
		return nil, err
	}

	var continuationToken string
	var backups []api.CloudBackupInfo

PageTraversal:
	for {
		if continuationToken == "" {
			logger.Infof("Querying first page of backups")
		} else {
			logger.Infof("Querying backups with continuation token %v", continuationToken)
		}

		// Enumerating can be expensive but we need to do it to get the original
		// volume name. Ark already has it so it can pass it down to us.
		// CloudBackupRestore can also be updated to restore to the original volume
		// name.
		enumRequest := &api.CloudBackupEnumerateRequest{}
		enumRequest.CredentialUUID = credID
		enumRequest.All = true
		enumRequest.ContinuationToken = continuationToken
		enumRequest.MaxBackups = 10000 // XXX: This is a hack because it looks like pagination doesn't work on the px api at all.
		enumResponse, err := volDriver.CloudBackupEnumerate(enumRequest)
		if err != nil {
			logger.WithError(err)
			return nil, err
		}

		logger.Infof("There are %v backups on this page", len(enumResponse.Backups))

		backups = append(backups, enumResponse.Backups...)

		continuationToken = enumResponse.ContinuationToken
		if continuationToken == "" {
			// No continuations. We're on the last page.
			break PageTraversal
		}
	}

	return backups, nil
}
