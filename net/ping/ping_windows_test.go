package ping

import "testing"

func TestPing(t *testing.T) {
	type args struct {
		destination string
		option      Option
	}
	tests := []struct {
		name    string
		args    args
		want    Result
		wantErr bool
	}{
		{
			name: "192.168.1.1",
			args: args{
				destination: "192.168.1.1",
				option: &PingOption{
					Count: 4,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Ping(tt.args.destination, tt.args.option)
			if (err != nil) != tt.wantErr {
				t.Errorf("Ping() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Ping() = %v, want %v", got, tt.want)
			}
		})
	}
}
