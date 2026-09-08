package filemedia

import "testing"

func TestSupportedCodecsAcceptsCurrentH264AACContract(
	t *testing.T,
) {
	video, audio := supportedCodecs(`
		Stream #0:0: Video: h264 (High)
		Stream #0:1: Audio: aac (LC)
	`)

	if !video {
		t.Fatal("video = false, want h264 accepted")
	}
	if !audio {
		t.Fatal("audio = false, want aac accepted")
	}
}

func TestSupportedCodecsAcceptsLegacyAliases(
	t *testing.T,
) {
	video, audio := supportedCodecs(`
		Video: avc
		Audio: mp4a
	`)

	if !video || !audio {
		t.Fatalf(
			"video=%v audio=%v, want both accepted",
			video,
			audio,
		)
	}
}

func TestSupportedCodecsRejectsUnsupportedCodec(
	t *testing.T,
) {
	video, audio := supportedCodecs(`
		Video: hevc
		Audio: opus
	`)

	if video || audio {
		t.Fatalf(
			"video=%v audio=%v, want both rejected",
			video,
			audio,
		)
	}
}
