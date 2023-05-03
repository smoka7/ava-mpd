package brainz

import (
	"reflect"
	"testing"
)

const resp = `{"images":[{"approved":true,"back":false,"comment":"","edit":83638502,"front":true,"id":30688716524,"image":"http://coverartarchive.org/release/73d5bede-e281-4e08-b00e-d81dcae2e75c/30688716524.jpg","thumbnails":{"1200":"http://coverartarchive.org/release/73d5bede-e281-4e08-b00e-d81dcae2e75c/30688716524-1200.jpg","250":"http://coverartarchive.org/release/73d5bede-e281-4e08-b00e-d81dcae2e75c/30688716524-250.jpg","500":"http://coverartarchive.org/release/73d5bede-e281-4e08-b00e-d81dcae2e75c/30688716524-500.jpg","large":"http://coverartarchive.org/release/73d5bede-e281-4e08-b00e-d81dcae2e75c/30688716524-500.jpg","small":"http://coverartarchive.org/release/73d5bede-e281-4e08-b00e-d81dcae2e75c/30688716524-250.jpg"},"types":["Front"]}],"release":"https://musicbrainz.org/release/73d5bede-e281-4e08-b00e-d81dcae2e75c"}`

func TestFetch(t *testing.T) {
	type args struct {
		releaseID string
	}
	tests := []struct {
		name     string
		args     args
		wantBody []byte
		wantErr  bool
	}{{name: "one", args: args{releaseID: CAURL + "73d5bede-e281-4e08-b00e-d81dcae2e75c"}, wantBody: []byte(resp), wantErr: false}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBody, err := fetch(tt.args.releaseID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Fetch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotBody, tt.wantBody) {
				t.Errorf("Fetch() = %v, want %v", gotBody, tt.wantBody)
			}
		})
	}
}

func TestGetImageURL(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		wantURL string
		args    args
	}{
		{
			name: "one",
			args: args{
				b: []byte(resp),
			},
			wantURL: "http://coverartarchive.org/release/73d5bede-e281-4e08-b00e-d81dcae2e75c/30688716524.jpg",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotURL := getCoverURL(tt.args.b); gotURL != tt.wantURL {
				t.Errorf("GetImageURL() = %v, want %v", gotURL, tt.wantURL)
			}
		})
	}
}
