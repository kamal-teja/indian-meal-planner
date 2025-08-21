package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFlexibleDate_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected time.Time
		wantErr  bool
	}{
		{
			name:     "RFC3339 with timezone",
			input:    `"2023-10-15T14:30:00+07:00"`,
			expected: time.Date(2023, 10, 15, 14, 30, 0, 0, time.FixedZone("", 7*3600)),
			wantErr:  false,
		},
		{
			name:     "RFC3339 UTC",
			input:    `"2023-10-15T14:30:00Z"`,
			expected: time.Date(2023, 10, 15, 14, 30, 0, 0, time.UTC),
			wantErr:  false,
		},
		{
			name:     "ISO8601 without timezone",
			input:    `"2023-10-15T14:30:00"`,
			expected: time.Date(2023, 10, 15, 14, 30, 0, 0, time.UTC),
			wantErr:  false,
		},
		{
			name:     "SQL datetime format",
			input:    `"2023-10-15 14:30:00"`,
			expected: time.Date(2023, 10, 15, 14, 30, 0, 0, time.UTC),
			wantErr:  false,
		},
		{
			name:     "Date only YYYY-MM-DD",
			input:    `"2023-10-15"`,
			expected: time.Date(2023, 10, 15, 0, 0, 0, 0, time.UTC),
			wantErr:  false,
		},
		{
			name:    "Invalid format",
			input:   `"invalid-date"`,
			wantErr: true,
		},
		{
			name:    "Empty string",
			input:   `""`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fd FlexibleDate
			err := fd.UnmarshalJSON([]byte(tt.input))

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "unable to parse date")
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected.Year(), fd.Time.Year())
				assert.Equal(t, tt.expected.Month(), fd.Time.Month())
				assert.Equal(t, tt.expected.Day(), fd.Time.Day())
				assert.Equal(t, tt.expected.Hour(), fd.Time.Hour())
				assert.Equal(t, tt.expected.Minute(), fd.Time.Minute())
				assert.Equal(t, tt.expected.Second(), fd.Time.Second())
			}
		})
	}
}

func TestFlexibleDate_MarshalJSON(t *testing.T) {
	// Arrange
	testTime := time.Date(2023, 10, 15, 14, 30, 0, 0, time.UTC)
	fd := FlexibleDate{Time: testTime}

	// Act
	data, err := fd.MarshalJSON()

	// Assert
	require.NoError(t, err)
	
	// Unmarshal to verify the format
	var result string
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)
	
	// Should be in RFC3339 format
	assert.Equal(t, "2023-10-15T14:30:00Z", result)
}

func TestFlexibleDate_String(t *testing.T) {
	// Arrange
	testTime := time.Date(2023, 10, 15, 14, 30, 0, 0, time.UTC)
	fd := FlexibleDate{Time: testTime}

	// Act
	result := fd.String()

	// Assert
	assert.Equal(t, "2023-10-15T14:30:00Z", result)
}

func TestFlexibleDate_JSONRoundTrip(t *testing.T) {
	// Test different date formats can be unmarshaled and then marshaled consistently
	testCases := []string{
		`"2023-10-15"`,
		`"2023-10-15T14:30:00"`,
		`"2023-10-15T14:30:00Z"`,
		`"2023-10-15 14:30:00"`,
	}

	for _, input := range testCases {
		t.Run("roundtrip_"+input, func(t *testing.T) {
			// Unmarshal
			var fd FlexibleDate
			err := fd.UnmarshalJSON([]byte(input))
			require.NoError(t, err)

			// Marshal
			data, err := fd.MarshalJSON()
			require.NoError(t, err)

			// Unmarshal again
			var fd2 FlexibleDate
			err = fd2.UnmarshalJSON(data)
			require.NoError(t, err)

			// Should be equal
			assert.True(t, fd.Time.Equal(fd2.Time))
		})
	}
}

func TestFlexibleDate_StructUsage(t *testing.T) {
	// Test that FlexibleDate works properly in struct serialization/deserialization
	type TestStruct struct {
		Date FlexibleDate `json:"date"`
		Name string       `json:"name"`
	}

	// Test with different date formats
	testCases := []struct {
		name string
		json string
	}{
		{
			name: "date only",
			json: `{"date":"2023-10-15","name":"test"}`,
		},
		{
			name: "full datetime",
			json: `{"date":"2023-10-15T14:30:00Z","name":"test"}`,
		},
		{
			name: "SQL format",
			json: `{"date":"2023-10-15 14:30:00","name":"test"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Unmarshal
			var ts TestStruct
			err := json.Unmarshal([]byte(tc.json), &ts)
			require.NoError(t, err)
			assert.Equal(t, "test", ts.Name)
			assert.Equal(t, 2023, ts.Date.Time.Year())
			assert.Equal(t, time.October, ts.Date.Time.Month())
			assert.Equal(t, 15, ts.Date.Time.Day())

			// Marshal back
			data, err := json.Marshal(ts)
			require.NoError(t, err)
			
			// Should contain RFC3339 format
			var result map[string]interface{}
			err = json.Unmarshal(data, &result)
			require.NoError(t, err)
			
			dateStr, ok := result["date"].(string)
			require.True(t, ok)
			
			// Should be in RFC3339 format regardless of input format
			_, err = time.Parse(time.RFC3339, dateStr)
			assert.NoError(t, err)
		})
	}
}

func TestFlexibleDate_ErrorMessage(t *testing.T) {
	// Test that error message is informative
	var fd FlexibleDate
	err := fd.UnmarshalJSON([]byte(`"not-a-date"`))
	
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unable to parse date")
	assert.Contains(t, err.Error(), "not-a-date")
	assert.Contains(t, err.Error(), "supported formats")
}