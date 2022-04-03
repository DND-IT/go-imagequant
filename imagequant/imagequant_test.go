package imagequant

import (
	"image"
	"reflect"
	"testing"
)

func TestRun(t *testing.T) {
	type args struct {
		imgRGBA *image.RGBA
		gamma   float64
	}
	tests := []struct {
		name       string
		args       args
		wantPalImg image.Image
		wantErr    bool
	}{
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPalImg, err := Run(tt.args.imgRGBA, tt.args.gamma)
			if (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotPalImg, tt.wantPalImg) {
				t.Errorf("Run() gotPalImg = %v, want %v", gotPalImg, tt.wantPalImg)
			}
		})
	}
}

