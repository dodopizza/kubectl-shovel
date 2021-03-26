package deadman

import (
	"os"
	"time"
)

const operatorIsDeadAfterNSeconds int64 = 30

// IsOperatorAlive - watch if operator still alive. Check if file exists and its mod time is not obsolete
func IsOperatorAlive() (bool, error) {
	stat, err := os.Stat(aliveFile)
	if os.IsNotExist(err) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	if stat.ModTime().Unix()+operatorIsDeadAfterNSeconds < time.Now().Unix() {
		return false, nil
	}

	return true, nil
}
