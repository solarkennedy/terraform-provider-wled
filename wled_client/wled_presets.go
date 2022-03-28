package wled_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

type WLEDPresetID string

type WLEDPresets map[WLEDPresetID]WLEDPreset

type WLEDPreset struct {
	// Name of the preset
	Name string `json:"n"`
	// On/Off state of the light
	On bool `json:"on"`
	// Brightness of the light. If on is false, contains last brightness when light was on (aka brightness when on is set to true. Setting bri to 0 is supported but it is recommended to use the range 1-255 and use on: false to turn off. The state response will never havethe value 0 for bri.
	Brightness int `json:"bri"`
	// Duration of the crossfade between different colors/brightness levels. One unit is 100ms, so a value of 4 results in atransition of 400ms.
	Transition int `json:"transition"`
	// Main Segment ID
	Mainseg int `json:"mainseg"`
	// Segments are individual parts of the LED strip. In 0.9.0 this will enables running different effects on differentparts of the strip.
	Segments []WLEDSegment `json:"seg"`
}

// https://kno.wled.ge/interfaces/json-api/#contents-of-the-segment-object
type WLEDSegment struct {
	// Zero-indexed ID of the segment. May be omitted, in that case the ID will be inferred from the order of the segment objects in the seg array.
	ID int `json:"id"`
	// Grouping (how many consecutive LEDs of the same segment will be grouped to the same color)
	Grouping int `json:"grp"`
	// Spacing (how many LEDs are turned off and skipped between each group)
	Spacing int `json:"spc"`
	// Offset (how many LEDs to rotate the virtual start of the segments, available since 0.13.0)
	Offfset int `json:"of"`
	// Turns on and off the individual segment. (available since 0.10.0)
	On bool `json:"on"`
	// Segment contents will not be refreshed
	Feeze bool `json:"frz"`
	// Brightness of the light. If on is false, contains last brightness when light was on (aka brightness when on is set to true. Setting bri to 0 is supported but it is recommended to use the range 1-255 and use on: false to turn off. The state response will never havethe value 0 for bri.
	Brightness int `json:"bri"`
	// White spectrum color temperature (available since 0.13.0)
	ColorTemperature int `json:"cct"`
	// Array that has up to 3 color arrays as elements, the primary, secondary (background) and tertiary colors of the segment. Each color is an array of 3 or 4 bytes, which represent an RGB(W) color.
	ColorArray [][]int `json:"col"`
	// ID of the effect or ~ to increment, ~- to decrement, or r for random.
	EffectID int `json:"fx"`
	// Relative effect speed (0-255)
	EffectSpeed int `json:"sx"`
	// Effect intensity
	EffectIntensity int `json:"ix"`
	// ID of the color palette or ~ to increment, ~- to decrement, or r for random.
	PaletteID int `json:"pal"`
	// true if the segment is selected. Selected segments will have their state (color/FX) updated by APIs that don't support segments (e.g. UDP sync, HTTP API). If no segment is selected, the first segment (id:0) will behave as if selected. WLED will report the state of the first (lowest id) segment that is selected to APIs (HTTP, MQTT, Blynk...), or mainseg in case no segment is selected and for the UDP API. Live data is always applied to all LEDs regardless of segment configuration.
	Selected bool `json:"sel"`
	// Flips the segment, causing animations to change direction.
	Reversed bool `json:"rev"`
	// Mirrors the segment (available since 0.10.2)
	Mirrored bool `json:"mi"`
}

func (c *WLEDClient) GetPreset(id WLEDPresetID) (WLEDPreset, bool, error) {
	allPresets, err := c.GetPresets()
	if err != nil {
		return WLEDPreset{}, false, err
	}
	preset, ok := allPresets[id]
	return preset, ok, err
}

func (c *WLEDClient) SetPreset(id WLEDPresetID, preset WLEDPreset) error {
	allPresets, err := c.GetPresets()
	if err != nil {
		return err
	}
	allPresets[id] = preset
	return c.SetPresets(allPresets)
}

func (c *WLEDClient) GetPresets() (WLEDPresets, error) {
	presets := make(WLEDPresets)
	url := "http://" + c.host + "/presets.json"
	resp, err := http.Get(url)
	if err != nil {
		return presets, fmt.Errorf("Error fetching %s: %w", url, err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return presets, fmt.Errorf("Error fetching body %s: %w", url, err)
	}
	err = json.Unmarshal(body, &presets)
	if err != nil {
		return presets, fmt.Errorf("Error parsing json for '%s': %w", body, err)
	}
	return presets, nil
}

func (c *WLEDClient) SetPresets(presets WLEDPresets) error {
	fileName := "presets.json"
	url := "http://" + c.host + "/upload"

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("filename", fileName)
	if err != nil {
		return fmt.Errorf("Error creating form file: %w", err)
	}

	b, err := json.Marshal(presets)
	if err != nil {
		return fmt.Errorf("Error marshaling presets: %w", err)
	}
	file := bytes.NewReader(b)
	_, err = io.Copy(part, file)
	if err != nil {
		return fmt.Errorf("Error marshaling presets: %w", err)
	}
	writer.Close()

	r, err := http.NewRequest("POST", url, body)
	r.Header.Add("Content-Type", writer.FormDataContentType())
	if err != nil {
		return fmt.Errorf("Error creating form file request: %w", err)
	}
	client := &http.Client{}
	_, err = client.Do(r)
	return err
}
