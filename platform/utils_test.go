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
