package wled_client

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// https://kno.wled.ge/interfaces/http-api/#xml-response
type WLEDSettings struct {
	XMLName     xml.Name `xml:"vs"`
	Description string   `xml:"ds"`
	// 	<ac>0</ac>
	// <cl>255</cl>
	// <cl>255</cl>
	// <cl>255</cl>
	// <cs>0</cs>
	// <cs>0</cs>
	// <cs>0</cs>
	// <ns>0</ns>
	// <nr>1</nr>
	// <nl>0</nl>
	// <nf>1</nf>
	// <nd>60</nd>
	// <nt>0</nt>
	// <fx>0</fx>
	// <sx>132</sx>
	// <ix>255</ix>
	// <fp>1</fp>
	// <wv>-1</wv>
	// <ws>0</ws>
	// <ps>0</ps>
	// <cy>0</cy>
	// <ss>0</ss>
}

func (c *WLEDClient) GetSettings() (WLEDSettings, error) {
	settings := WLEDSettings{}
	url := c.getSettingsURL("win")

	resp, err := http.Get(url)
	if err != nil {
		return settings, fmt.Errorf("Error fetching %s: %w", url, err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return settings, fmt.Errorf("Error fetching body %s: %w", url, err)
	}
	err = xml.Unmarshal(body, &settings)
	if err != nil {
		return settings, fmt.Errorf("Error parsing xml for '%s': %w", body, err)
	}
	return settings, nil
}

func (c *WLEDClient) SetSettings(new_settings WLEDSettings) error {
	surl := c.getSettingsURL("settings/ui?")

	data := url.Values{}

	if new_settings.Description != "" {
		data["DS"] = []string{new_settings.Description}
	}

	_, err := http.PostForm(surl, data)
	if err != nil {
		return fmt.Errorf("Error posting body %+v to %s: %w", data, surl, err)
	}
	return nil
}

func (c *WLEDClient) getSettingsURL(extra string) string {
	return "http://" + c.host + "/" + extra
}
