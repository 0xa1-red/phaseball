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
		label := fmt.Sprintf("Roll: %d - Event: %s", tt.roll, tt.hit.Short())
		tf := func(t *testing.T) {
			t.Logf("Case %d - %s\n", i+1, label)
			{
				actualEvent, actualOut, actualExtra := defense(tt.hit, tt.roll)

				if expected, actual := tt.expectedEvent.Short(), actualEvent.Short(); expected != actual {
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
