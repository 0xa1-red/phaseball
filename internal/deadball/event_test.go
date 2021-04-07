package deadball

import (
	"fmt"
	"testing"
)

func TestOut(t *testing.T) {
	tests := []struct {
		swing    int
		expected Event
	}{
		{
			swing:    10,
			expected: EventOutK,
		},
		{
			swing:    11,
			expected: EventOutK,
		},
		{
			swing:    12,
			expected: EventOutK,
		},
		{
			swing:    13,
			expected: EventOutG3,
		},
		{
			swing:    14,
			expected: EventOut43,
		},
		{
			swing:    15,
			expected: EventOut53,
		},
		{
			swing:    16,
			expected: EventOut63,
		},
		{
			swing:    17,
			expected: EventOutF7,
		},
		{
			swing:    18,
			expected: EventOutF8,
		},
		{
			swing:    19,
			expected: EventOutF9,
		},
	}

	for i, tt := range tests {
		tf := func(t *testing.T) {
			t.Logf("Case %d - Swing: %d\n", i+1, tt.swing)
			{
				if expected, actual := tt.expected, Out(tt.swing); expected != actual {
					t.Fatalf("Fail: expected event to be %s, got %s", expected, actual)
				}
				t.Log("Pass")
			}
		}

		t.Run(fmt.Sprintf("swing_%d", tt.swing), tf)
	}
}

func TestDefense(t *testing.T) {
	tests := []struct {
		roll          int
		hit           Event
		expectedEvent Event
		expectedOut   bool
		expectedExtra bool
	}{
		{
			roll:          12,
			hit:           EventHitSingle1B,
			expectedEvent: EventHitSingle1B,
			expectedOut:   true,
			expectedExtra: false,
		},
		{
			roll:          10,
			hit:           EventHitDoubleCF,
			expectedEvent: EventHitSingleAdv2,
			expectedOut:   false,
			expectedExtra: false,
		},
		{
			roll:          10,
			hit:           EventHitTripleCF,
			expectedEvent: EventHitDoubleAdv3,
			expectedOut:   false,
			expectedExtra: false,
		},
		{
			roll:          3,
			hit:           EventHitSingle1B,
			expectedEvent: EventHitSingle1B,
			expectedOut:   false,
			expectedExtra: false,
		},
		{
			roll:          1,
			hit:           EventHitSingle1B,
			expectedEvent: EventHitSingle1B,
			expectedOut:   false,
			expectedExtra: true,
		},
	}

	for i, tt := range tests {
		label := fmt.Sprintf("Roll: %d - Event: %s", tt.roll, tt.hit)
		tf := func(t *testing.T) {
			t.Logf("Case %d - %s\n", i+1, label)
			{
				actualEvent, actualOut, actualExtra := defense(tt.hit, tt.roll)

				if expected, actual := tt.expectedEvent, actualEvent; expected != actual {
					t.Fatalf("Fail: expected event %s, got %s", expected, actual)
				}

				if expected, actual := tt.expectedOut, actualOut; expected != actual {
					t.Fatalf("Fail: expected out to be %t, got %t", expected, actual)
				}

				if expected, actual := tt.expectedExtra, actualExtra; expected != actual {
					t.Fatalf("Fail: expected extra to be %t, got %t", expected, actual)
				}

				t.Log("Pass")
			}
		}

		t.Run(label, tf)
	}
}

func TestExtendedEvent_GetLong(t *testing.T) {
	if expected, actual := "Strikeout", EventOutK.Long; expected != actual {
		t.Fatalf("Fail: expected %s, got %s", expected, actual)
	}

	testEvent := Event{Label: "TEST"}
	if expected, actual := "TEST", testEvent.GetLong(); expected != actual {
		t.Fatalf("Fail: expected %s, got %s", expected, actual)
	}

	t.Log("Pass")
}

func TestHit(t *testing.T) {
	tests := []struct {
		swing          int
		crit           bool
		roll           int
		expectedResult Event
		expectedExtra  bool
		expectedOut    bool
	}{
		{
			roll:           19,
			expectedResult: EventHitHomeRun,
			expectedExtra:  false,
			expectedOut:    false,
		},
		{
			roll:           18,
			crit:           true,
			expectedResult: EventHitHomeRun,
			expectedExtra:  false,
			expectedOut:    false,
		},
		{
			roll:           16,
			expectedResult: EventHitDoubleAdv3,
			expectedExtra:  false,
			expectedOut:    false,
		},
		{
			roll:           15,
			crit:           true,
			expectedResult: EventHitDoubleAdv3,
			expectedExtra:  false,
			expectedOut:    false,
		},
		{
			roll:           8,
			expectedResult: EventHitSingleAdv2,
			expectedExtra:  false,
			expectedOut:    false,
		},
		{
			roll:           7,
			crit:           true,
			expectedResult: EventHitSingleAdv2,
			expectedExtra:  false,
			expectedOut:    false,
		},
		{
			roll:           1,
			expectedResult: EventHitSinglePlus,
			expectedExtra:  false,
			expectedOut:    false,
		},
	}
	for i, tt := range tests {
		label := fmt.Sprintf("Roll: %d - Crit: %t - Swing: %d", tt.roll, tt.crit, tt.swing)
		tf := func(t *testing.T) {
			t.Logf("Case %d - %s\n", i+1, label)
			{
				actualEvent, actualExtra, actualOut := hit(tt.swing, tt.crit, tt.roll)

				if expected, actual := tt.expectedResult, actualEvent; expected != actual {
					t.Fatalf("Fail: expected event %s, got %s", expected, actual)
				}

				if expected, actual := tt.expectedOut, actualOut; expected != actual {
					t.Fatalf("Fail: expected out to be %t, got %t", expected, actual)
				}

				if expected, actual := tt.expectedExtra, actualExtra; expected != actual {
					t.Fatalf("Fail: expected extra to be %t, got %t", expected, actual)
				}

				t.Log("Pass")
			}
		}

		t.Run(label, tf)
	}
}

func TestIsOutOutfield(t *testing.T) {
	tests := []struct {
		out        Event
		isOutfield bool
	}{
		{
			out:        EventOutG3,
			isOutfield: false,
		},
		{
			out:        EventOutF7,
			isOutfield: true,
		},
	}

	for i, tt := range tests {
		label := tt.out.Label
		tf := func(t *testing.T) {
			t.Logf("Case %d - %s\n", i+1, label)
			{
				if expected, actual := tt.isOutfield, IsOutOutfield(tt.out); expected != actual {
					t.Fatalf("Fail: expected %t, got %t", expected, actual)
				}

				t.Log("Pass")
			}
		}

		t.Run(label, tf)
	}
}

func TestIsOutInfield(t *testing.T) {
	tests := []struct {
		out       Event
		isInfield bool
	}{
		{
			out:       EventOutG3,
			isInfield: true,
		},
		{
			out:       EventOutF7,
			isInfield: false,
		},
	}

	for i, tt := range tests {
		label := tt.out.Label
		tf := func(t *testing.T) {
			t.Logf("Case %d - %s\n", i+1, label)
			{
				if expected, actual := tt.isInfield, IsOutInfield(tt.out); expected != actual {
					t.Fatalf("Fail: expected %t, got %t", expected, actual)
				}

				t.Log("Pass")
			}
		}

		t.Run(label, tf)
	}
}
