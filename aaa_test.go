package aaa

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestApplicationPayload_Perms(t *testing.T) {

	tc := []ApplicationPayload{{
		UserID:           1,
		UserLogin:        "user_test",
		RoleID:           2,
		PermissionBitSet: json.RawMessage([]byte(`"e65f"`)),
		IsDebug:          false,
		ExtraPayload:     nil,
	},
	}

	for i := range tc {
		if !bytes.Equal(tc[i].Perms(), []byte(`e65f`)) {
			t.Error("failed check #1")
		}
	}

}
