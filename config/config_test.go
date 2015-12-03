package config

import "testing"

func TestDefault(t *testing.T) {
	val := valuedef("1", "default")
	if val != "1" {
		t.Errorf("Excepted %q, got %q", "1", val)
	}
	val = valuedef(nil, "default")
	if val != "default" {
		t.Errorf("Excepted %q, got %q", "default", val)
	}
}

func TestIArr2Str(t *testing.T) {
	sarr := iArr2sArr(nil)
	if sarr != nil {
		t.Errorf("Excepted return nil iArr2sArr(nil),got %v", sarr)
	}
	iarr := []interface{}{
		"1", "2", "23", "34",
	}
	sarr = iArr2sArr(iarr)
	if len(sarr) != len(iarr) {
		t.Errorf("Excepted iArr2sArr len %d,got %d", len(sarr), len(iarr))
	}
	for i := 0; i < len(sarr); i++ {
		if sarr[i] != iarr[i].(string) {
			t.Errorf("Excepted iArr2sArr [%d] %s,got %s", i, iarr[i].(string), sarr[i])
		}
	}
}

func TestParse(t *testing.T) {
	data := []byte(`
	{
		"source": "zbar",
		"zbar": {
			"enabled": true,
			"location": "/usr/bin/zbarcam",		
			"device": "/dev/video1",
			"args": [
				"1",
				"2"
			]
		}
	}`)
	cfg, e := parse(data)
	if e != nil {
		t.Error(e)
	}
	if cfg == nil {
		t.Error("Excepted parse returns not nil")
	}
	if cfg.Source != SourceZBar {
		t.Errorf("Excepted %s, got %s", SourceZBar, cfg.Source)
	}
	if cfg.ZBar == nil {
		t.Error("Excepted Config.ZBar returns not nil")
	}
	if cfg.ZBar.Device != "/dev/video1" {
		t.Errorf("Excepted %s, got %s", "/dev/video1", cfg.ZBar.Device)
	}
	if cfg.ZBar.Location != "/usr/bin/zbarcam" {
		t.Errorf("Excepted %s, got %s", "/usr/bin/zbarcam", cfg.ZBar.Location)
	}
	if len(cfg.ZBar.Args) != 2 || cfg.ZBar.Args[0] != "1" || cfg.ZBar.Args[1] != "2" {
		t.Errorf("Excepted %v, got %v", []string{"1", "2"}, cfg.ZBar.Args)
	}
}
