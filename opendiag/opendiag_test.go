package opendiag

import (
	_ "embed"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

//go:embed appLog-2023-09-16-16-25-12.log
var testFile string

func TestDecode(t *testing.T) {
	parsed, err := DateFromFileName("appLog-2023-09-16-16-25-12.log")
	require.NoError(t, err)
	require.Equal(t, time.Date(2023, 9, 16, 16, 25, 12, 0, time.Local), parsed)

	expected := Log{
		Header: `AppVersion: 2.17.13
Android SDK: 31 (12)
Android device: Xiaomi M2007J3SY
ECU: ВАЗ: BOSCH MP7.0 E2
Connect: Bluetooth
State: connecting 16:25:15,885
Device: 66:1E:11:0F:71:77
State: connected 16:25:16,833`,
		Entries: Entries{
			{
				Time: time.Date(2023, 9, 16, 16, 25, 16, 934000000, time.Local),
				Send: "Send:	AT@1",
				Receive: `Receive: OBDII TO RS232 INTERPRETER
123`,
			},
		},
	}

	olog, err := Decode(parsed, testFile)
	require.NoError(t, err)
	require.Equal(t, expected, olog)

	expectedData := `AppVersion: 2.17.13
Android SDK: 31 (12)
Android device: Xiaomi M2007J3SY
ECU: ВАЗ: BOSCH MP7.0 E2
Connect: Bluetooth
State: connecting 16:25:15,885
Device: 66:1E:11:0F:71:77
State: connected 16:25:16,833
Time:	16:25:16,934
Send:	AT@1
Receive: OBDII TO RS232 INTERPRETER
123
`

	name, data := olog.Encode()
	require.Equal(t, "appLog-2023-09-16-16-25-16.log", name)
	require.Equal(t, expectedData, data)
}
