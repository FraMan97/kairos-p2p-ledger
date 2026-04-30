package config

import (
	"log"
	"os"
)

var (
	OtsDigestURL  string
	OtsUpgradeURL string
)

func InitConfig() {
	OtsDigestURL = os.Getenv("OTS_DIGEST_URL")
	if OtsDigestURL == "" {
		OtsDigestURL = "https://a.pool.opentimestamps.org/digest"
	}

	OtsUpgradeURL = os.Getenv("OTS_UPGRADE_URL")
	if OtsUpgradeURL == "" {
		OtsUpgradeURL = "https://a.pool.opentimestamps.org/upgrade"
	}

	log.Printf("[Config] OTS_DIGEST_URL configured to: %s", OtsDigestURL)
	log.Printf("[Config] OTS_UPGRADE_URL configured to: %s", OtsUpgradeURL)
}
