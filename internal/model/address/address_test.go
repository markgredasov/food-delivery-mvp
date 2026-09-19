package address_test

import (
	"database/sql/driver"
	"encoding/json"
	"strings"
	"testing"

	errs "github.com/markgredasov/food-delivery-mvp/internal/errors"
	address "github.com/markgredasov/food-delivery-mvp/internal/model/address"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	okComment := "entrance from the yard"
	emptyComment := ""
	spacesComment := "     "
	okSpacesComment := "   entrance from the yard    "
	tests := []struct {
		name      string
		city      string
		street    string
		house     string
		apartment string
		comment   *string
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "valid address without comment",
			city:      "Moscow",
			street:    "Tverskaya",
			house:     "1",
			apartment: "10",
			comment:   nil,
			wantErr:   false,
		},
		{
			name:      "valid address with comment",
			city:      "Saint Petersburg",
			street:    "Nevsky Prospekt",
			house:     "22",
			apartment: "5",
			comment:   &okComment,
			wantErr:   false,
		},
		{
			name:      "empty city",
			city:      "",
			street:    "Tverskaya",
			house:     "1",
			apartment: "10",
			comment:   nil,
			wantErr:   true,
			errMsg:    "city must not be empty",
		},
		{
			name:      "city with only spaces",
			city:      "   ",
			street:    "Tverskaya",
			house:     "1",
			apartment: "10",
			comment:   nil,
			wantErr:   true,
			errMsg:    "city must not be empty",
		},
		{
			name:      "city too long",
			city:      string(make([]rune, 101)),
			street:    "Tverskaya",
			house:     "1",
			apartment: "10",
			comment:   nil,
			wantErr:   true,
			errMsg:    "city is too long",
		},
		{
			name:      "empty street",
			city:      "Moscow",
			street:    "",
			house:     "1",
			apartment: "10",
			comment:   nil,
			wantErr:   true,
			errMsg:    "street must not be empty",
		},
		{
			name:      "street too long",
			city:      "Moscow",
			street:    string(make([]rune, 101)),
			house:     "1",
			apartment: "10",
			comment:   nil,
			wantErr:   true,
			errMsg:    "street is too long",
		},
		{
			name:      "empty house",
			city:      "Moscow",
			street:    "Tverskaya",
			house:     "",
			apartment: "10",
			comment:   nil,
			wantErr:   true,
			errMsg:    "house must not be empty",
		},
		{
			name:      "house too long",
			city:      "Moscow",
			street:    "Tverskaya",
			house:     string(make([]rune, 101)),
			apartment: "10",
			comment:   nil,
			wantErr:   true,
			errMsg:    "house is too long",
		},
		{
			name:      "empty apartment",
			city:      "Moscow",
			street:    "Tverskaya",
			house:     "1",
			apartment: "",
			comment:   nil,
			wantErr:   true,
			errMsg:    "apartment must not be empty",
		},
		{
			name:      "apartment too long",
			city:      "Moscow",
			street:    "Tverskaya",
			house:     "1",
			apartment: string(make([]rune, 101)),
			comment:   nil,
			wantErr:   true,
			errMsg:    "apartment is too long",
		},
		{
			name:      "empty comment",
			city:      "Moscow",
			street:    "Tverskaya",
			house:     "1",
			apartment: "10",
			comment:   &emptyComment,
			wantErr:   true,
			errMsg:    "comment must not be empty",
		},
		{
			name:      "comment with only spaces",
			city:      "Moscow",
			street:    "Tverskaya",
			house:     "1",
			apartment: "10",
			comment:   &spacesComment,
			wantErr:   true,
			errMsg:    "comment must not be empty",
		},
		{
			name:      "trims spaces from all fields",
			city:      "  Moscow  ",
			street:    "  Tverskaya  ",
			house:     "  1  ",
			apartment: "  10  ",
			comment:   &okSpacesComment,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addr, err := address.New(tt.city, tt.street, tt.house, tt.apartment, tt.comment)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)

				assert.ErrorIs(t, err, errs.ErrInvalidArgument)
			} else {
				require.NoError(t, err)
				assert.Equal(t, strings.TrimSpace(tt.city), addr.City)
				assert.Equal(t, strings.TrimSpace(tt.street), addr.Street)
				assert.Equal(t, strings.TrimSpace(tt.house), addr.House)
				assert.Equal(t, strings.TrimSpace(tt.apartment), addr.Apartment)

				if tt.comment != nil {
					assert.Equal(t, strings.TrimSpace(*tt.comment), addr.Comment)
				} else {
					assert.Empty(t, addr.Comment)
				}
			}
		})
	}
}

func TestAddress_Scan(t *testing.T) {
	tests := []struct {
		name    string
		src     any
		want    address.Address
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid JSON bytes",
			src:  []byte(`{"city":"Moscow","street":"Tverskaya","house":"1","apartment":"10","comment":"test"}`),
			want: address.Address{
				City:      "Moscow",
				Street:    "Tverskaya",
				House:     "1",
				Apartment: "10",
				Comment:   "test",
			},
			wantErr: false,
		},
		{
			name: "valid JSON string",
			src:  `{"city":"Moscow","street":"Tverskaya","house":"1","apartment":"10"}`,
			want: address.Address{
				City:      "Moscow",
				Street:    "Tverskaya",
				House:     "1",
				Apartment: "10",
				Comment:   "",
			},
			wantErr: false,
		},
		{
			name:    "nil source",
			src:     nil,
			want:    address.Address{},
			wantErr: true,
			errMsg:  "cannot scan nil into Address",
		},
		{
			name:    "invalid type",
			src:     123,
			want:    address.Address{},
			wantErr: true,
			errMsg:  "cannot scan int into Address",
		},
		{
			name:    "invalid JSON",
			src:     []byte(`{invalid json}`),
			want:    address.Address{},
			wantErr: true,
			errMsg:  "failed to unmarshal address",
		},
		{
			name: "partial JSON (missing fields)",
			src:  []byte(`{"city":"Moscow"}`),
			want: address.Address{
				City:      "Moscow",
				Street:    "",
				House:     "",
				Apartment: "",
				Comment:   "",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var addr address.Address
			err := addr.Scan(tt.src)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, addr)
			}
		})
	}
}

func TestAddress_Value(t *testing.T) {
	tests := []struct {
		name    string
		addr    address.Address
		want    driver.Value
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid address with all fields",
			addr: address.Address{
				City:      "Moscow",
				Street:    "Tverskaya",
				House:     "1",
				Apartment: "10",
				Comment:   "test",
			},
			want:    []byte(`{"city":"Moscow","street":"Tverskaya","house":"1","apartment":"10","comment":"test"}`),
			wantErr: false,
		},
		{
			name: "valid address without comment",
			addr: address.Address{
				City:      "Moscow",
				Street:    "Tverskaya",
				House:     "1",
				Apartment: "10",
				Comment:   "",
			},
			want:    []byte(`{"city":"Moscow","street":"Tverskaya","house":"1","apartment":"10","comment":""}`),
			wantErr: false,
		},
		{
			name: "missing city",
			addr: address.Address{
				City:      "",
				Street:    "Tverskaya",
				House:     "1",
				Apartment: "10",
			},
			wantErr: true,
			errMsg:  "invalid address: city, street and house are required",
		},
		{
			name: "missing street",
			addr: address.Address{
				City:      "Moscow",
				Street:    "",
				House:     "1",
				Apartment: "10",
			},
			wantErr: true,
			errMsg:  "invalid address: city, street and house are required",
		},
		{
			name: "missing house",
			addr: address.Address{
				City:      "Moscow",
				Street:    "Tverskaya",
				House:     "",
				Apartment: "10",
			},
			wantErr: true,
			errMsg:  "invalid address: city, street and house are required",
		},
		{
			name: "address with special characters",
			addr: address.Address{
				City:      "Moscow",
				Street:    "Tverskaya 1/2",
				House:     "1A",
				Apartment: "10/2",
				Comment:   "Entrance #2, door with code 1234",
			},
			want:    []byte(`{"apartment":"10/2", "city":"Moscow", "comment":"Entrance #2, door with code 1234", "house":"1A", "street":"Tverskaya 1/2"}`),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, err := tt.addr.Value()

			if tt.wantErr {
				require.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				require.NoError(t, err)

				var expectedMap, actualMap map[string]any

				if tt.want != nil {
					err = json.Unmarshal(tt.want.([]byte), &expectedMap)
					require.NoError(t, err)
				}

				err = json.Unmarshal(val.([]byte), &actualMap)
				require.NoError(t, err)

				assert.Equal(t, expectedMap, actualMap)
			}
		})
	}
}

func TestAddress_ScanValueRoundTrip(t *testing.T) {
	original := address.Address{
		City:      "Moscow",
		Street:    "Tverskaya",
		House:     "1",
		Apartment: "10",
		Comment:   "test comment",
	}

	val, err := original.Value()
	require.NoError(t, err)

	var scanned address.Address
	err = scanned.Scan(val)
	require.NoError(t, err)

	assert.Equal(t, original, scanned)
}

func TestAddress_EmptyCommentHandling(t *testing.T) {
	emptyComment := ""
	_, err := address.New("Moscow", "Tverskaya", "1", "10", &emptyComment)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "comment must not be empty")

	addr, err := address.New("Moscow", "Tverskaya", "1", "10", nil)
	require.NoError(t, err)
	assert.Empty(t, addr.Comment)
}

func TestNew_MaxLengthBoundary(t *testing.T) {
	longString := string(make([]rune, 100))

	addr, err := address.New(longString, "Street", "1", "10", nil)
	require.NoError(t, err)
	assert.Equal(t, longString, addr.City)

	tooLongString := string(make([]rune, 101))
	_, err = address.New(tooLongString, "Street", "1", "10", nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "city is too long")
}
