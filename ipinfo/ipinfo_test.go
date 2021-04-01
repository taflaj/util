// ipinfo_test.go

package ipinfo

import "testing"

func TestIP(t *testing.T) {
	tests := []struct {
		ip      string
		err     bool
		bogon   bool
		country string
	}{
		{"70.188.219.55", false, false, "US"},
		{"127.0.0.1", false, true, ""},
		{"185.220.101.204", false, false, "DE"},
		{"192.168.1.1", false, true, ""},
		{"209.127.17.234", false, false, "CA"},
		{"AS15169", true, false, ""},
	}
	for _, test := range tests {
		info, err := GetInfo(test.ip)
		if err != nil {
			t.Error(err)
		} else if info == nil {
			t.Error("No response returned")
		} else if test.err {
			if info.Error == "" {
				t.Error("Expected error but didn't get any")
			}
		} else if test.ip != info.IP {
			t.Errorf("Requested %s but got %s", test.ip, info.IP)
		} else if test.bogon != info.Bogon {
			t.Errorf("Bogon expected %v but got %v", test.bogon, info.Bogon)
		} else if test.country != info.Country {
			t.Errorf("Country expected %s but got %s", test.country, info.Country)
		}
	}
}
