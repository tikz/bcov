package utils

import "testing"

func TestSHA256FileHash(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"little.bam open", args{path: "x"}, "", true},
		{"little.bam hash", args{path: "../tests/little.bam"}, "238598669fe0d1a510b0ed76577943b6142e6fef29f082536a399ef3225c4ad4", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SHA256FileHash(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("SHA256FileHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SHA256FileHash() = %v, want %v", got, tt.want)
			}
		})
	}
}
