//go:build linux

package main

import (
	"testing"
)

func TestUtils_RebuildStreamURL(t *testing.T) {
	urlSamples := []struct {
		url     string
		rebuild string
	}{
		{url: "rtsp://121.1.2.3", rebuild: "rtsp://121.1.2.3"},
		{url: "rtsp://121.1.2.3/Streaming/Channels/101", rebuild: "rtsp://121.1.2.3/Streaming/Channels/101"},
		{url: "rtsp://121.1.2.3:554/Streaming/Channels/101", rebuild: "rtsp://121.1.2.3:554/Streaming/Channels/101"},
		{url: "rtsp://121.1.2.3:554/Streaming/Channels/101?k=v", rebuild: "rtsp://121.1.2.3:554/Streaming/Channels/101?k=v"},
		{url: "rtsp://CamViewer:abc123@121.1.2.3:554/Streaming/Channels/101", rebuild: "rtsp://CamViewer:abc123@121.1.2.3:554/Streaming/Channels/101"},
		{url: "rtsp://CamViewer:abc123?!@121.1.2.3:554/Streaming/Channels/101", rebuild: "rtsp://CamViewer:abc123%3F%21@121.1.2.3:554/Streaming/Channels/101"},
		{url: "rtsp://CamViewer:abc123@?!@121.1.2.3:554/Streaming/Channels/101", rebuild: "rtsp://CamViewer:abc123%40%3F%21@121.1.2.3:554/Streaming/Channels/101"},
		{url: "rtsp://CamViewer:abc123@?!@121.1.2.3:554/Streaming/Channels/101?k=v", rebuild: "rtsp://CamViewer:abc123%40%3F%21@121.1.2.3:554/Streaming/Channels/101?k=v"},
		{url: "rtsp://CamViewer:abc123@?!@121.1.2.3:554", rebuild: "rtsp://CamViewer:abc123%40%3F%21@121.1.2.3:554"},
		{url: "rtsp://Cam@Viewer:abc123@?!@121.1.2.3:554", rebuild: "rtsp://Cam%40Viewer:abc123%40%3F%21@121.1.2.3:554"},
		{url: "rtsp://CamViewer:abc123@?!~#$%^&*()_+-=\\|?@121.1.2.3:554/Streaming/Channels/101", rebuild: "rtsp://CamViewer:abc123%40%3F%21~%23$%25%5E&%2A%28%29_+-=%5C%7C%3F@121.1.2.3:554/Streaming/Channels/101"},
		{url: "rtsp://CamViewer:abc123@347?1!@121.1.2.3:554/Streaming/Channels/101", rebuild: "rtsp://CamViewer:abc123%40347%3F1%21@121.1.2.3:554/Streaming/Channels/101"},
		{url: "srt://213.171.194.158:10080", rebuild: "srt://213.171.194.158:10080"},
		{url: "srt://213.171.194.158:10080?streamid=#!::r=live/primary,latency=20,m=request", rebuild: "srt://213.171.194.158:10080?streamid=#!::r=live/primary,latency=20,m=request"},
	}
	for _, urlSample := range urlSamples {
		if r0, err := RebuildStreamURL(urlSample.url); err != nil {
			t.Errorf("Fail for err %+v", err)
			return
		} else if rebuild := r0.String(); rebuild != urlSample.rebuild {
			t.Errorf("rebuild url %v failed, expect %v, actual %v",
				urlSample.url, urlSample.rebuild, rebuild)
			return
		}
	}
}

func TestUtils_ParseFFmpegLogs(t *testing.T) {
	for _, e := range []struct {
		log   string
		ts    string
		speed string
	}{
		{log: "time=00:10:09.138 speed=1x", ts: "00:10:09.138", speed: "1x"},
		{log: "size=18859kB time=00:10:09.138 speed=1x", ts: "00:10:09.138", speed: "1x"},
		{log: "size=18859kB time=00:10:09.138 speed=1x dup=1", ts: "00:10:09.138", speed: "1x"},
		{log: "size=18859kB time=00:10:09.138 bitrate=253.5kbits/s speed=1x dup=1", ts: "00:10:09.138", speed: "1x"},
		{log: "size=18859kB time=00:10:09.38 bitrate=253.5kbits/s speed=1x", ts: "00:10:09.38", speed: "1x"},
		{log: "frame=184 fps=9.7 q=28.0 size=364kB time=00:00:19.41 bitrate=153.7kbits/s dup=0 drop=235 speed=1.03x", ts: "00:00:19.41", speed: "1.03x"},
	} {
		if ts, speed, err := ParseFFmpegCycleLog(e.log); err != nil {
			t.Errorf("Fail parse %v for err %+v", e, err)
		} else if ts != e.ts {
			t.Errorf("Fail for ts %v of %v", ts, e)
		} else if speed != e.speed {
			t.Errorf("Fail for speed %v of %v", speed, e)
		}
	}
}
