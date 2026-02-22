package cover

import "testing"

func TestCoverImageExtract(t *testing.T) {
	tests := []struct {
		name   string
		images []Image
		want   string
	}{
		{
			name: "prefer front cover over back cover",
			images: []Image{
				{URL: "https://example.com/back.jpg", Back: true, Width: 2000, Height: 2000},
				{URL: "https://example.com/front.jpg", Front: true, Width: 500, Height: 500},
			},
			want: "https://example.com/front.jpg",
		},
		{
			name: "prefer higher resolution when cover side is same",
			images: []Image{
				{URL: "https://example.com/front-small.jpg", Front: true, Width: 250, Height: 250},
				{URL: "https://example.com/front-large.jpg", Front: true, Width: 1200, Height: 1200},
			},
			want: "https://example.com/front-large.jpg",
		},
		{
			name: "prefer largest thumbnail when original dimensions are unknown",
			images: []Image{
				{
					URL:   "https://example.com/full.jpg",
					Front: true,
					Thumbnails: map[string]string{
						"250":  "https://example.com/250.jpg",
						"1200": "https://example.com/1200.jpg",
					},
				},
			},
			want: "https://example.com/1200.jpg",
		},
		{
			name: "return empty when no valid image exists",
			images: []Image{
				{URL: "", Thumbnails: map[string]string{"small": ""}},
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CoverImageExtract(tt.images)
			if got != tt.want {
				t.Fatalf("CoverImageExtract() = %q, want %q", got, tt.want)
			}
		})
	}
}
