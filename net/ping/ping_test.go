package ping

import (
	"reflect"
	"testing"
)

func TestPingOptions_Validate(t *testing.T) {
	type fields struct {
		Count    int
		Interval int
		Deadline int
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &PingOption{
				Count:    tt.fields.Count,
				Interval: tt.fields.Interval,
				Deadline: tt.fields.Deadline,
			}
			o.Validate()
		})
	}
}

func Test_parseLinuxPing(t *testing.T) {
	type args struct {
		output []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *PingResult
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseLinuxPing(tt.args.output)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseLinuxPing() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseLinuxPing() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseWindowsPing(t *testing.T) {
	type args struct {
		output []byte
	}
	tests := []struct {
		name    string
		args    args
		want    Result
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				output: []byte(`
Active code page: 437

Pinging www.google.com [216.58.221.164] with 32 bytes of data:
Reply from 216.58.221.164: bytes=32 time=37ms TTL=51
Request timed out.
Reply from 216.58.221.164: bytes=32 time=38ms TTL=51
Reply from 216.58.221.164: bytes=32 time=37ms TTL=51
Reply from 216.58.221.164: bytes=32 time=37ms TTL=51

Ping statistics for 216.58.221.164:
    Packets: Sent = 5, Received = 4, Lost = 1 (20% loss),
Approximate round trip times in milli-seconds:
    Minimum = 37ms, Maximum = 38ms, Average = 37ms

				`),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseWindowsPing(tt.args.output)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseWindowsPing() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseWindowsPing() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseDarwinPing(t *testing.T) {
	type args struct {
		output []byte
	}
	tests := []struct {
		name    string
		args    args
		want    Result
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				output: []byte(`
PING 192.168.1.1 (192.168.1.1): 56 data bytes
64 bytes from 192.168.1.1: icmp_seq=0 ttl=64 time=1.664 ms
64 bytes from 192.168.1.1: icmp_seq=1 ttl=64 time=0.971 ms
64 bytes from 192.168.1.1: icmp_seq=2 ttl=64 time=3.934 ms
64 bytes from 192.168.1.1: icmp_seq=3 ttl=64 time=3.539 ms
64 bytes from 192.168.1.1: icmp_seq=4 ttl=64 time=3.690 ms

--- 192.168.1.1 ping statistics ---
5 packets transmitted, 5 packets received, 0.0% packet loss
round-trip min/avg/max/stddev = 0.971/2.760/3.934/1.204 ms

				`),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseDarwinPing(tt.args.output)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseDarwinPing() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseDarwinPing() = %v, want %v", got, tt.want)
			}
		})
	}
}
