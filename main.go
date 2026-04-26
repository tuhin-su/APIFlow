package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var version = "dev"
var vars = map[string]string{}

type Config struct {
	AppURL    string            `json:"app_url"`
	Loop      int               `json:"loop"`
	Delay     int               `json:"delay"`
	Headers   map[string]string `json:"headers"`
	Proxy     *ProxyConfig      `json:"proxy,omitempty"`
	Endpoints []Endpoint        `json:"endpoints"`
}

type ProxyConfig struct {
	URL      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Endpoint struct {
	Name   string                 `json:"name"`
	URL    string                 `json:"url"`
	Method string                 `json:"method"`
	Body   map[string]interface{} `json:"body,omitempty"`
	Set    map[string]string      `json:"set,omitempty"`
}

func main() {
	configPath := flag.String("load", "", "Path to config file")
	configPathShort := flag.String("l", "", "Path to config file (shorthand)")
	showVersion := flag.Bool("v", false, "Show version")

	flag.Parse()

	if *showVersion {
		fmt.Println("Version:", version)
		return
	}

	path := *configPath
	if path == "" {
		path = *configPathShort
	}

	if path == "" {
		fmt.Println("❌ Config file required")
		fmt.Println("Usage: flowgo --load config.json")
		os.Exit(1)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("❌ Failed to read file:", err)
		os.Exit(1)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		fmt.Println("❌ Invalid JSON:", err)
		os.Exit(1)
	}

	run(cfg)
}

func run(cfg Config) {
	client := buildHTTPClient(cfg.Proxy)

	loop := cfg.Loop
	if loop == -1 {
		loop = int(^uint(0) >> 1)
	}

	for i := 0; i < loop; i++ {
		fmt.Println("🔁 Loop:", i+1)

		for _, ep := range cfg.Endpoints {
			err := executeEndpoint(client, cfg, ep)
			if err != nil {
				fmt.Println("❌ Error:", err)
				break
			}
		}

		time.Sleep(time.Duration(cfg.Delay) * time.Millisecond)
	}
}

func buildHTTPClient(proxyCfg *ProxyConfig) *http.Client {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: 10 * time.Second,
		}).DialContext,
	}

	if proxyCfg != nil && proxyCfg.URL != "" {
		proxyURL, _ := url.Parse(proxyCfg.URL)

		if proxyCfg.Username != "" {
			proxyURL.User = url.UserPassword(proxyCfg.Username, proxyCfg.Password)
		}

		transport.Proxy = http.ProxyURL(proxyURL)
	}

	return &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}
}

func executeEndpoint(client *http.Client, cfg Config, ep Endpoint) error {
	fullURL := cfg.AppURL + ep.URL

	var body io.Reader

	if ep.Body != nil {
		bodyMap := replaceVarsInMap(ep.Body)
		jsonBody, _ := json.Marshal(bodyMap)
		body = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(ep.Method, fullURL, body)
	if err != nil {
		return err
	}

	for k, v := range cfg.Headers {
		req.Header.Set(k, replaceVars(v))
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	fmt.Println("➡️", ep.Name, resp.Status)

	var jsonResp map[string]interface{}
	json.Unmarshal(respBody, &jsonResp)

	for key, path := range ep.Set {
		val := extractFromJSON(jsonResp, path)
		if val != "" {
			vars[key] = val
			fmt.Println("🔐 SET", key, "=", val)
		}
	}

	return nil
}

func replaceVars(input string) string {
	for k, v := range vars {
		input = strings.ReplaceAll(input, k, v)
	}
	return input
}

func replaceVarsInMap(m map[string]interface{}) map[string]interface{} {
	newMap := make(map[string]interface{})
	for k, v := range m {
		switch val := v.(type) {
		case string:
			newMap[k] = replaceVars(val)
		default:
			newMap[k] = val
		}
	}
	return newMap
}

func extractFromJSON(data map[string]interface{}, path string) string {
	parts := strings.Split(path, ".")
	var current interface{} = data

	for _, p := range parts {
		if m, ok := current.(map[string]interface{}); ok {
			current = m[p]
		} else {
			return ""
		}
	}

	return fmt.Sprintf("%v", current)
}
