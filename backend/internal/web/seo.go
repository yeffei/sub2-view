//go:build embed

package web

import (
	"bytes"
	"encoding/json"
	"html"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

const (
	defaultSEOSiteName    = "山枢庭"
	defaultSEOTagline     = "统一入口，安静流转。"
	defaultSEODescription = "山枢庭 SST 提供统一的 AI API 入口，将接入、计量、账册与稳定供给整理在同一套长期可维护的服务中。"
)

type seoPage struct {
	Title       string
	Description string
	Path        string
	Robots      string
}

type routeSEO struct {
	Title       string
	Description string
	Canonical   string
	Robots      string
	SiteName    string
	Public      bool
	CacheKey    string
}

var publicSEOPages = map[string]seoPage{
	"/home": {
		Title:       defaultSEOTagline,
		Description: "山枢庭 SST 是统一的 AI API 入口，将接入、计量、账册与稳定供给整理在一套安静、长期可维护的服务中。",
		Path:        "/home",
		Robots:      "index,follow",
	},
	"/pricing": {
		Title:       "模型价格",
		Description: "查看山枢庭 SST 的 Codex 与 Claude Code 模型价格、倍率与计量口径，理解 API 调用成本和开通边界。",
		Path:        "/pricing",
		Robots:      "index,follow",
	},
	"/docs": {
		Title:       "接入文档",
		Description: "山枢庭 SST API 接入文档，说明 OpenAI 兼容 base_url、API Key、模型列表、流式输出和常见错误排查。",
		Path:        "/docs",
		Robots:      "index,follow",
	},
	"/docs/openai-compatible-api": {
		Title:       "OpenAI 兼容 API 接入",
		Description: "山枢庭 SST 保持 OpenAI 兼容调用方式，便于现有 SDK、CLI 与服务端应用接入统一入口。",
		Path:        "/docs/openai-compatible-api",
		Robots:      "index,follow",
	},
	"/docs/base-url": {
		Title:       "base_url 与 endpoint 配置",
		Description: "配置 SST base_url 时应保留 /v1 前缀，避免把业务路径、网关路径和模型路径混在一起。",
		Path:        "/docs/base-url",
		Robots:      "index,follow",
	},
	"/docs/api-key": {
		Title:       "API Key 创建与使用",
		Description: "了解 SST API Key 的创建、保存、调用、轮换和权限边界，避免在公开代码或客户端泄露凭据。",
		Path:        "/docs/api-key",
		Robots:      "index,follow",
	},
	"/docs/streaming": {
		Title:       "流式输出与 SSE 调用",
		Description: "SST 支持 OpenAI 兼容的流式输出方式，适用于聊天窗口、终端客户端和需要逐段展示的应用。",
		Path:        "/docs/streaming",
		Robots:      "index,follow",
	},
	"/docs/codex": {
		Title:       "Codex 客户端接入",
		Description: "使用 SST 为 Codex 类客户端配置统一 API 入口、Key、模型与账册，便于集中计量和维护。",
		Path:        "/docs/codex",
		Robots:      "index,follow",
	},
	"/docs/claude-code": {
		Title:       "Claude Code 接入说明",
		Description: "通过 SST 统一管理 Claude Code 类客户端的入口、Key、模型权限、倍率和调用账册。",
		Path:        "/docs/claude-code",
		Robots:      "index,follow",
	},
	"/faq": {
		Title:       "常见问题",
		Description: "山枢庭 SST 常见问题入口，说明价格、开通、控制台和实际配置应以公开价格页与登录后系统状态为准。",
		Path:        "/faq",
		Robots:      "index,follow",
	},
	"/privacy": {
		Title:       "隐私政策",
		Description: "山枢庭 SST 隐私政策，说明账户、账册、请求、安全信息和第三方服务边界的数据处理原则。",
		Path:        "/privacy",
		Robots:      "index,follow",
	},
	"/terms": {
		Title:       "服务条款",
		Description: "山枢庭 SST 服务条款，说明账户、API Key、计费、订单、第三方上游、可用性和终止条件等规则。",
		Path:        "/terms",
		Robots:      "index,follow",
	},
}

var (
	managedMetaRe = regexp.MustCompile(`(?is)\s*<meta\s+(?:name|property)=["'](?:description|robots|twitter:card|twitter:title|twitter:description|twitter:image|og:site_name|og:title|og:description|og:type|og:url|og:image)["'][^>]*>`)
	canonicalRe   = regexp.MustCompile(`(?is)\s*<link\s+rel=["']canonical["'][^>]*>`)
	titleRe       = regexp.MustCompile(`(?is)<title>.*?</title>`)
)

func resolveRouteSEO(req *http.Request, settingsJSON []byte) routeSEO {
	path := normalizeSEOPath(req.URL.Path)
	origin := requestOrigin(req)
	siteName := resolveSEOSiteName(settingsJSON)

	page, ok := publicSEOPages[path]
	if !ok {
		page = seoPage{
			Title:       siteName,
			Description: defaultSEODescription,
			Path:        path,
			Robots:      "noindex,nofollow",
		}
	}

	title := buildSEOTitle(siteName, page.Title)
	canonicalPath := page.Path
	if canonicalPath == "" {
		canonicalPath = path
	}

	return routeSEO{
		Title:       title,
		Description: page.Description,
		Canonical:   origin + canonicalPath,
		Robots:      page.Robots,
		SiteName:    siteName,
		Public:      ok,
		CacheKey:    origin + "|" + path,
	}
}

func injectRouteSEO(htmlBytes []byte, seo routeSEO) []byte {
	result := managedMetaRe.ReplaceAll(htmlBytes, nil)
	result = canonicalRe.ReplaceAll(result, nil)
	title := []byte("<title>" + html.EscapeString(seo.Title) + "</title>")
	if titleRe.Match(result) {
		result = titleRe.ReplaceAll(result, title)
	} else {
		result = bytes.Replace(result, []byte("</head>"), append(title, []byte("</head>")...), 1)
	}

	headClose := []byte("</head>")
	block := []byte(buildSEOHeadBlock(seo))
	return bytes.Replace(result, headClose, append(block, headClose...), 1)
}

func buildSEOHeadBlock(seo routeSEO) string {
	var b strings.Builder
	b.WriteString(`<meta name="description" content="`)
	b.WriteString(html.EscapeString(seo.Description))
	b.WriteString(`">`)
	b.WriteString(`<meta name="robots" content="`)
	b.WriteString(html.EscapeString(seo.Robots))
	b.WriteString(`">`)
	b.WriteString(`<link rel="canonical" href="`)
	b.WriteString(html.EscapeString(seo.Canonical))
	b.WriteString(`">`)
	b.WriteString(`<meta property="og:site_name" content="`)
	b.WriteString(html.EscapeString(seo.SiteName))
	b.WriteString(`">`)
	b.WriteString(`<meta property="og:title" content="`)
	b.WriteString(html.EscapeString(seo.Title))
	b.WriteString(`">`)
	b.WriteString(`<meta property="og:description" content="`)
	b.WriteString(html.EscapeString(seo.Description))
	b.WriteString(`">`)
	b.WriteString(`<meta property="og:type" content="website">`)
	b.WriteString(`<meta property="og:url" content="`)
	b.WriteString(html.EscapeString(seo.Canonical))
	b.WriteString(`">`)
	b.WriteString(`<meta property="og:image" content="`)
	b.WriteString(html.EscapeString(requestURLWithPath(seo.Canonical, "/logo.png")))
	b.WriteString(`">`)
	b.WriteString(`<meta name="twitter:card" content="summary_large_image">`)
	b.WriteString(`<meta name="twitter:title" content="`)
	b.WriteString(html.EscapeString(seo.Title))
	b.WriteString(`">`)
	b.WriteString(`<meta name="twitter:description" content="`)
	b.WriteString(html.EscapeString(seo.Description))
	b.WriteString(`">`)
	b.WriteString(`<meta name="twitter:image" content="`)
	b.WriteString(html.EscapeString(requestURLWithPath(seo.Canonical, "/logo.png")))
	b.WriteString(`">`)

	if seo.Public {
		b.WriteString(buildStructuredDataScript(seo))
	}

	return b.String()
}

func buildStructuredDataScript(seo routeSEO) string {
	graph := map[string]any{
		"@context": "https://schema.org",
		"@graph": []map[string]any{
			{
				"@type":       "Organization",
				"name":        seo.SiteName,
				"url":         requestURLWithPath(seo.Canonical, "/home"),
				"description": defaultSEODescription,
				"logo":        requestURLWithPath(seo.Canonical, "/logo.png"),
			},
			{
				"@type":       "WebSite",
				"name":        seo.SiteName,
				"url":         requestURLWithPath(seo.Canonical, "/home"),
				"description": defaultSEODescription,
			},
		},
	}
	payload, err := json.Marshal(graph)
	if err != nil {
		return ""
	}
	return `<script type="application/ld+json" nonce="` + NonceHTMLPlaceholder + `">` + string(payload) + `</script>`
}

func buildSEOTitle(siteName, pageTitle string) string {
	if strings.TrimSpace(pageTitle) == "" || pageTitle == siteName {
		return siteName
	}
	if pageTitle == defaultSEOTagline {
		return siteName + " - " + defaultSEOTagline
	}
	return pageTitle + " - " + siteName
}

func resolveSEOSiteName(settingsJSON []byte) string {
	var cfg struct {
		SiteName string `json:"site_name"`
	}
	if err := json.Unmarshal(settingsJSON, &cfg); err == nil {
		if trimmed := strings.TrimSpace(cfg.SiteName); trimmed != "" && trimmed != "Sub2API" {
			return trimmed
		}
	}
	return defaultSEOSiteName
}

func normalizeSEOPath(path string) string {
	trimmed := strings.TrimSpace(path)
	if trimmed == "" || trimmed == "/" || trimmed == "/index.html" {
		return "/home"
	}
	if i := strings.IndexAny(trimmed, "?#"); i >= 0 {
		trimmed = trimmed[:i]
	}
	if len(trimmed) > 1 {
		trimmed = strings.TrimRight(trimmed, "/")
	}
	if !strings.HasPrefix(trimmed, "/") {
		return "/" + trimmed
	}
	return trimmed
}

func requestOrigin(req *http.Request) string {
	scheme := firstHeaderValue(req.Header.Get("X-Forwarded-Proto"))
	if scheme == "" {
		if req.TLS != nil {
			scheme = "https"
		} else {
			scheme = "http"
		}
	}
	if scheme != "http" && scheme != "https" {
		scheme = "https"
	}
	host := firstHeaderValue(req.Header.Get("X-Forwarded-Host"))
	if host == "" {
		host = req.Host
	}
	if host == "" {
		host = "localhost"
	}
	return scheme + "://" + host
}

func firstHeaderValue(value string) string {
	if i := strings.Index(value, ","); i >= 0 {
		value = value[:i]
	}
	return strings.TrimSpace(value)
}

func requestURLWithPath(rawURL, path string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}
	parsed.Path = path
	parsed.RawQuery = ""
	parsed.Fragment = ""
	return parsed.String()
}

func buildRobotsTXT(req *http.Request) string {
	origin := requestOrigin(req)
	return "User-agent: *\n" +
		"Allow: /home\n" +
		"Allow: /pricing\n" +
		"Allow: /docs\n" +
		"Allow: /docs/openai-compatible-api\n" +
		"Allow: /docs/base-url\n" +
		"Allow: /docs/api-key\n" +
		"Allow: /docs/streaming\n" +
		"Allow: /docs/codex\n" +
		"Allow: /docs/claude-code\n" +
		"Allow: /faq\n" +
		"Allow: /privacy\n" +
		"Allow: /terms\n" +
		"Disallow: /admin\n" +
		"Disallow: /dashboard\n" +
		"Disallow: /keys\n" +
		"Disallow: /usage\n" +
		"Disallow: /profile\n" +
		"Disallow: /login\n" +
		"Disallow: /register\n" +
		"Disallow: /payment\n" +
		"Disallow: /setup\n" +
		"Sitemap: " + origin + "/sitemap.xml\n"
}

func buildSitemapXML(req *http.Request) string {
	origin := requestOrigin(req)
	paths := []string{
		"/home",
		"/pricing",
		"/docs",
		"/docs/openai-compatible-api",
		"/docs/base-url",
		"/docs/api-key",
		"/docs/streaming",
		"/docs/codex",
		"/docs/claude-code",
		"/faq",
		"/privacy",
		"/terms",
	}
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	b.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`)
	for _, path := range paths {
		b.WriteString(`<url><loc>`)
		b.WriteString(html.EscapeString(origin + path))
		b.WriteString(`</loc></url>`)
	}
	b.WriteString(`</urlset>`)
	return b.String()
}
