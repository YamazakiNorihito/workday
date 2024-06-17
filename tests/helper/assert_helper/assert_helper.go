package assert_helper

import (
	"testing"
	"time"

	"github.com/YamazakiNorihito/workday/internal/domain/metadata"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func EqualUUID(t *testing.T, expected uuid.UUID, actual interface{}, msgAndArgs ...interface{}) bool {
	t.Helper()
	var actualUUID uuid.UUID
	var err error

	switch v := actual.(type) {
	/*case []byte:
	actualUUID, err = uuid.FromBytes(v)*/
	case string:
		actualUUID, err = uuid.Parse(v)
	default:
		t.Fatalf("unsupported type for actual: %T", actual)
		return false
	}

	if err != nil {
		t.Fatalf("failed to convert ID to UUID: %v", err)
		return false
	}

	return assert.Equal(t, expected, actualUUID, msgAndArgs...)
}

func EqualUserMeta(t *testing.T, expected metadata.UserMeta, actual interface{}, msgAndArgs ...interface{}) bool {
	t.Helper()
	userMetaMap := map[string]interface{}{"id": expected.ID, "name": expected.Name}
	return assert.Subset(t, userMetaMap, actual, msgAndArgs...)
}

func EqualUnixTime(t *testing.T, expected time.Time, actual interface{}, msgAndArgs ...interface{}) bool {
	t.Helper()
	expectedUnix := expected.Unix()

	var actualUnix int64
	switch v := actual.(type) {
	case float64:
		actualUnix = int64(v)
	case int64:
		actualUnix = v
	default:
		t.Fatalf("unsupported type for actual: %T", actual)
		return false
	}

	return assert.Equal(t, expectedUnix, actualUnix, msgAndArgs...)
}
