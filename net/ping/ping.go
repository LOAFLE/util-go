package ping

import (
	"log"
	"strconv"
	"strings"
)

type Option interface {
	GetCount() int
	GetInterval() int
	GetDeadline() int
	Validate()
}

type Response interface {
	GetTTL() int
	GetTime() float64
	GetError() string
}

type Summary interface {
	GetSendCount() int
	GetReceiveCount() int
	GetLossPercent() float32
	GetMinTime() float64
	GetMaxTime() float64
	GetAvgTime() float64
}

type Result interface {
	GetResponses() map[int]Response
	GetSummary() Summary
	GetRaw() []string
}

type PingOption struct {
	Count    int `json:"count,omitempty"`
	Interval int `json:"interval,omitempty"`
	Deadline int `json:"deadline,omitempty"`
}

func (o *PingOption) GetCount() int {
	return o.Count
}
func (o *PingOption) GetInterval() int {
	return o.Interval
}
func (o *PingOption) GetDeadline() int {
	return o.Deadline
}
func (o *PingOption) Validate() {
	if 0 >= o.Count {
		o.Count = 1
	}
	if 0 >= o.Interval {
		o.Interval = 1
	}
	if 0 >= o.Deadline {
		o.Deadline = 1
	}
}

type PingResponse struct {
	TTL   int     `json:"ttl,omitempty"`
	Time  float64 `json:"time,omitempty"`
	Error string  `json:"error,omitempty"`
}

func (r *PingResponse) GetTTL() int {
	return r.TTL
}
func (r *PingResponse) GetTime() float64 {
	return r.Time
}
func (r *PingResponse) GetError() string {
	return r.Error
}

type PingSummary struct {
	SendCount    int     `json:"sendCount,omitempty"`
	ReceiveCount int     `json:"receiveCount,omitempty"`
	LossPercent  float32 `json:"lossPercent,omitempty"`
	MinTime      float64 `json:"minTime,omitempty"`
	MaxTime      float64 `json:"maxTime,omitempty"`
	AvgTime      float64 `json:"avgTime,omitempty"`
}

func (s *PingSummary) GetSendCount() int {
	return s.SendCount
}
func (s *PingSummary) GetReceiveCount() int {
	return s.ReceiveCount
}
func (s *PingSummary) GetLossPercent() float32 {
	return s.LossPercent
}
func (s *PingSummary) GetMinTime() float64 {
	return s.MinTime
}
func (s *PingSummary) GetMaxTime() float64 {
	return s.MaxTime
}
func (s *PingSummary) GetAvgTime() float64 {
	return s.AvgTime
}

type PingResult struct {
	Responses map[int]Response `json:"responses,omitempty"`
	Summary   Summary          `json:"summary,omitempty"`
	Raw       []string         `json:"raw,omitempty"`
}

func (r *PingResult) GetResponses() map[int]Response {
	return r.Responses
}
func (r *PingResult) GetSummary() Summary {
	return r.Summary
}
func (r *PingResult) GetRaw() []string {
	return r.Raw
}

// $ ping 192.168.1.1 -c 7 -i 1 -w 0.3
// PING 192.168.1.1 (192.168.1.1) 56(84) bytes of data.
// 64 bytes from 192.168.1.1: icmp_seq=1 ttl=64 time=0.325 ms
// 64 bytes from 192.168.1.1: icmp_seq=2 ttl=64 time=0.344 ms
// 64 bytes from 192.168.1.1: icmp_seq=3 ttl=64 time=0.163 ms
// 64 bytes from 192.168.1.1: icmp_seq=7 ttl=64 time=0.180 ms
// From 192.168.1.101 icmp_seq=1 Destination Host Unreachable

// --- 192.168.1.1 ping statistics ---
// 7 packets transmitted, 4 received, 42% packet loss, time 6010ms
// rtt min/avg/max/mdev = 0.163/0.253/0.344/0.082 ms
func parseLinuxPing(output []byte) (Result, error) {
	result := &PingResult{
		Responses: make(map[int]Response, 0),
		Summary:   &PingSummary{},
		Raw:       make([]string, 0),
	}
	lines := strings.Split(string(output), "\n")

LOOP:
	for _, line := range lines {
		result.Raw = append(result.Raw, line)
		fields := strings.Fields(line)
		switch len(fields) {
		case 5:
			if "rtt" != fields[0] {
				continue LOOP
			}

			times := strings.Split(fields[3], "/")

			minTime, err := strconv.ParseFloat(times[0], 64)
			if nil != err {
				log.Print(err)
				continue LOOP
			}
			result.Summary.(*PingSummary).MinTime = minTime

			maxTime, err := strconv.ParseFloat(times[2], 64)
			if nil != err {
				log.Print(err)
				continue LOOP
			}
			result.Summary.(*PingSummary).MaxTime = maxTime

			avgTime, err := strconv.ParseFloat(times[1], 64)
			if nil != err {
				log.Print(err)
				continue LOOP
			}
			result.Summary.(*PingSummary).AvgTime = avgTime

		case 8:
			if "bytes" != fields[1] || "from" != fields[2] {
				continue LOOP
			}

			seqs := strings.Split(fields[4], "=")
			ttls := strings.Split(fields[5], "=")
			times := strings.Split(fields[6], "=")
			seq, err := strconv.Atoi(seqs[1])
			if nil != err {
				log.Print(err)
				continue LOOP
			}
			ttl, err := strconv.Atoi(ttls[1])
			if nil != err {
				log.Print(err)
				continue LOOP
			}
			_time, err := strconv.ParseFloat(times[1], 64)
			if nil != err {
				log.Print(err)
				continue LOOP
			}

			result.Responses[seq] = &PingResponse{
				TTL:  ttl,
				Time: _time,
			}

		case 10:
			sendCount, err := strconv.Atoi(fields[0])
			if nil != err {
				log.Print(err)
				continue LOOP
			}
			result.Summary.(*PingSummary).SendCount = sendCount

			receiveCount, err := strconv.Atoi(fields[3])
			if nil != err {
				log.Print(err)
				continue LOOP
			}
			result.Summary.(*PingSummary).ReceiveCount = receiveCount

			lossPercent, err := strconv.ParseFloat(strings.Replace(fields[5], "%", "", -1), 32)
			if nil != err {
				log.Print(err)
				continue LOOP
			}
			result.Summary.(*PingSummary).LossPercent = float32(lossPercent)
		}
	}

	return result, nil
}

// Windows 10
// Active code page: 437

// Pinging 192.168.1.1 with 32 bytes of data:
// Reply from 192.168.1.1: bytes=32 time<1ms TTL=64
// Reply from 192.168.1.1: bytes=32 time<1ms TTL=64
// Reply from 192.168.1.1: bytes=32 time<1ms TTL=64
// Reply from 192.168.1.1: bytes=32 time<1ms TTL=64
// Reply from 192.168.1.1: bytes=32 time<1ms TTL=64

// Ping statistics for 192.168.1.1:
//     Packets: Sent = 5, Received = 5, Lost = 0 (0% loss),
// Approximate round trip times in milli-seconds:
//     Minimum = 0ms, Maximum = 0ms, Average = 0ms

// Active code page: 437

// Pinging www.google.com [216.58.221.164] with 32 bytes of data:
// Reply from 216.58.221.164: bytes=32 time=37ms TTL=51
// Request timed out.
// Reply from 216.58.221.164: bytes=32 time=38ms TTL=51
// Reply from 216.58.221.164: bytes=32 time=37ms TTL=51
// Reply from 216.58.221.164: bytes=32 time=37ms TTL=51

// Ping statistics for 216.58.221.164:
//     Packets: Sent = 5, Received = 4, Lost = 1 (20% loss),
// Approximate round trip times in milli-seconds:
//     Minimum = 37ms, Maximum = 38ms, Average = 37ms

func parseWindowsPing(output []byte) (Result, error) {
	result := &PingResult{
		Responses: make(map[int]Response, 0),
		Summary:   &PingSummary{},
		Raw:       make([]string, 0),
	}
	lines := strings.Split(string(output), "\n")

	seq := 1
LOOP:
	for _, line := range lines {
		result.Raw = append(result.Raw, line)
		fields := strings.Fields(line)
		switch len(fields) {
		case 3:
			if "Request timed out." != line {
				continue LOOP
			}
			// result.Responses[seq] = nil
			seq = seq + 1
		case 6:
			if "Reply" != fields[0] || "from" != fields[1] {
				continue LOOP
			}
			times := strings.Replace(fields[4], "time", "", -1)
			times = strings.Replace(times, "<", "", -1)
			times = strings.Replace(times, "=", "", -1)
			times = strings.Replace(times, "ms", "", -1)

			ttls := strings.Split(fields[5], "=")

			ttl, err := strconv.Atoi(ttls[1])
			if nil != err {
				log.Print(err)
				continue LOOP
			}
			_time, err := strconv.ParseFloat(times, 64)
			if nil != err {
				log.Print(err)
				continue LOOP
			}

			result.Responses[seq] = &PingResponse{
				TTL:  ttl,
				Time: _time,
			}
			seq = seq + 1
		case 9:
			if "Minimum" != fields[0] {
				continue LOOP
			}

			minTimes := strings.Replace(fields[2], "ms", "", -1)
			minTimes = strings.Replace(minTimes, ",", "", -1)
			minTime, err := strconv.ParseFloat(minTimes, 64)
			if nil != err {
				log.Print(err)
				continue LOOP
			}
			result.Summary.(*PingSummary).MinTime = minTime

			maxTimes := strings.Replace(fields[5], "ms", "", -1)
			maxTimes = strings.Replace(maxTimes, ",", "", -1)
			maxTime, err := strconv.ParseFloat(maxTimes, 64)
			if nil != err {
				log.Print(err)
				continue LOOP
			}
			result.Summary.(*PingSummary).MaxTime = maxTime

			avgTimes := strings.Replace(fields[8], "ms", "", -1)
			avgTime, err := strconv.ParseFloat(avgTimes, 64)
			if nil != err {
				log.Print(err)
				continue LOOP
			}
			result.Summary.(*PingSummary).AvgTime = avgTime

		case 12:
			if "Packets:" != fields[0] {
				continue LOOP
			}
			sendCount, err := strconv.Atoi(strings.Replace(fields[3], ",", "", -1))
			if nil != err {
				log.Print(err)
				continue LOOP
			}
			result.Summary.(*PingSummary).SendCount = sendCount

			receiveCount, err := strconv.Atoi(strings.Replace(fields[6], ",", "", -1))
			if nil != err {
				log.Print(err)
				continue LOOP
			}
			result.Summary.(*PingSummary).ReceiveCount = receiveCount

			lossPercents := strings.Replace(fields[10], "(", "", -1)
			lossPercents = strings.Replace(lossPercents, "%", "", -1)
			lossPercent, err := strconv.ParseFloat(lossPercents, 32)
			if nil != err {
				log.Print(err)
				continue LOOP
			}
			result.Summary.(*PingSummary).LossPercent = float32(lossPercent)
		}
	}

	return result, nil
}

// $ ping 192.168.1.1 -c 5 -i 1
// PING 192.168.1.1 (192.168.1.1): 56 data bytes
// 64 bytes from 192.168.1.1: icmp_seq=0 ttl=64 time=1.664 ms
// 64 bytes from 192.168.1.1: icmp_seq=1 ttl=64 time=0.971 ms
// 64 bytes from 192.168.1.1: icmp_seq=2 ttl=64 time=3.934 ms
// 64 bytes from 192.168.1.1: icmp_seq=3 ttl=64 time=3.539 ms
// 64 bytes from 192.168.1.1: icmp_seq=4 ttl=64 time=3.690 ms

// --- 192.168.1.1 ping statistics ---
// 5 packets transmitted, 5 packets received, 0.0% packet loss
// round-trip min/avg/max/stddev = 0.971/2.760/3.934/1.204 ms
func parseDarwinPing(output []byte) (Result, error) {
	result := &PingResult{
		Responses: make(map[int]Response, 0),
		Summary:   &PingSummary{},
		Raw:       make([]string, 0),
	}
	lines := strings.Split(string(output), "\n")

LOOP:
	for _, line := range lines {
		result.Raw = append(result.Raw, line)
		fields := strings.Fields(line)
		switch len(fields) {
		case 5:
			if "round-trip" != fields[0] {
				continue LOOP
			}

			times := strings.Split(fields[3], "/")

			minTime, err := strconv.ParseFloat(times[0], 64)
			if nil != err {
				log.Print(err)
				continue LOOP
			}
			result.Summary.(*PingSummary).MinTime = minTime

			maxTime, err := strconv.ParseFloat(times[2], 64)
			if nil != err {
				log.Print(err)
				continue LOOP
			}
			result.Summary.(*PingSummary).MaxTime = maxTime

			avgTime, err := strconv.ParseFloat(times[1], 64)
			if nil != err {
				log.Print(err)
				continue LOOP
			}
			result.Summary.(*PingSummary).AvgTime = avgTime

		case 8:
			if "bytes" != fields[1] || "from" != fields[2] {
				continue LOOP
			}

			seqs := strings.Split(fields[4], "=")
			ttls := strings.Split(fields[5], "=")
			times := strings.Split(fields[6], "=")
			seq, err := strconv.Atoi(seqs[1])
			if nil != err {
				log.Print(err)
				continue LOOP
			}
			ttl, err := strconv.Atoi(ttls[1])
			if nil != err {
				log.Print(err)
				continue LOOP
			}
			_time, err := strconv.ParseFloat(times[1], 64)
			if nil != err {
				log.Print(err)
				continue LOOP
			}

			result.Responses[seq] = &PingResponse{
				TTL:  ttl,
				Time: _time,
			}

		case 9:
			sendCount, err := strconv.Atoi(fields[0])
			if nil != err {
				log.Print(err)
				continue LOOP
			}
			result.Summary.(*PingSummary).SendCount = sendCount

			receiveCount, err := strconv.Atoi(fields[3])
			if nil != err {
				log.Print(err)
				continue LOOP
			}
			result.Summary.(*PingSummary).ReceiveCount = receiveCount

			lossPercent, err := strconv.ParseFloat(strings.Replace(fields[6], "%", "", -1), 32)
			if nil != err {
				log.Print(err)
				continue LOOP
			}
			result.Summary.(*PingSummary).LossPercent = float32(lossPercent)
		}
	}

	return result, nil
}
