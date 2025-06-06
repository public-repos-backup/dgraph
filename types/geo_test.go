/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package types

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	array := []string{
		`{'type':'Point','coordinates':[1,2]}`,
		`{'type':'MultiLineString','coordinates':[[[1,2,3],[4,5,6],[7,8,9],[1,2,3]]]}`,
	}
	for _, v := range array {
		src := Val{StringID, []byte(v)}

		if g, err := Convert(src, GeoID); err != nil {
			t.Errorf("Error parsing %s: %v", v, err)
		} else {
			// Marshal it back to text
			got := ValueForType(StringID)
			if err := Marshal(g, &got); err != nil || got.Value.(string) != v {
				t.Errorf("Marshal error expected %s, got %s. error %v", v, got.Value.(string), err)
			}

			wkb := ValueForType(BinaryID)
			// Marshal and unmarshal to WKB
			if err := Marshal(g, &wkb); err != nil {
				t.Errorf("Error marshaling to WKB: %v", err)
			}

			src := Val{GeoID, wkb.Value.([]byte)}
			if bg, err := Convert(src, GeoID); err != nil {
				t.Errorf("Error unmarshaling WKB: %v", err)
			} else if !reflect.DeepEqual(g, bg) {
				t.Errorf("Expected %#v, got %#v", g, bg)
			}
		}
	}
}

func TestParseGeoJsonErrors(t *testing.T) {
	array := []string{
		`{"type":"Curve","coordinates":[1,2]}`,
		`{"type":"Feature","geometry":{"type":"Point","coordinates":[125.6,10.1]},"properties":{"name":"Dinagat Islands"}}`,
		`{}`,
		`thisisntjson`,
	}
	for _, v := range array {
		src := Val{StringID, []byte(v)}
		if _, err := Convert(src, GeoID); err == nil {
			t.Errorf("Expected error parsing %s: %v", v, err)
		}
	}
}
