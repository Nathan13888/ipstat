package main

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})
	app.Get("/", getIP)
	app.Get("/ip", getIP)
	app.Get("/all", getAllHeaders)
	app.Get("/all.json", getAllHeadersJSON)

	app.Listen(":3000")
}

// SOURCE: https://github.com/Ferluci/fasthttp-realip/blob/master/realip.go

// Header may return multiple IP addresses in the format: "client IP, proxy 1 IP, proxy 2 IP", so we take the the first one.
var xOriginalForwardedForHeader = http.CanonicalHeaderKey("X-Original-Forwarded-For")
var xForwardedForHeader = http.CanonicalHeaderKey("X-Forwarded-For")
var xForwardedHeader = http.CanonicalHeaderKey("X-Forwarded")
var forwardedForHeader = http.CanonicalHeaderKey("Forwarded-For")
var forwardedHeader = http.CanonicalHeaderKey("Forwarded")

// Standard headers used by Amazon EC2, Heroku, and others
var xClientIPHeader = http.CanonicalHeaderKey("X-Client-IP")

// Nginx proxy/FastCGI
var xRealIPHeader = http.CanonicalHeaderKey("X-Real-IP")

// Cloudflare.
// @see https://support.cloudflare.com/hc/en-us/articles/200170986-How-does-Cloudflare-handle-HTTP-Request-headers-
// CF-Connecting-IP - applied to every request to the origin.
var cfConnectingIPHeader = http.CanonicalHeaderKey("CF-Connecting-IP")

// Fastly CDN and Firebase hosting header when forwared to a cloud function
var fastlyClientIPHeader = http.CanonicalHeaderKey("Fastly-Client-Ip")

// Akamai and Cloudflare
var trueClientIPHeader = http.CanonicalHeaderKey("True-Client-Ip")

func getIP(c *fiber.Ctx) error {
	xClientIP := c.Request().Header.Peek(xClientIPHeader)
	if xClientIP != nil {
		return c.SendString(string(xClientIP))
	}

	xOriginalForwardedFor := c.Request().Header.Peek(xOriginalForwardedForHeader)
	if xOriginalForwardedFor != nil {
		return c.SendString(string(xOriginalForwardedFor))
	}

	xForwardedFor := c.Request().Header.Peek(xForwardedForHeader)
	if xForwardedFor != nil {
		return c.SendString(string(xForwardedFor))
	}

	ipHeaders := [...]string{cfConnectingIPHeader, fastlyClientIPHeader, trueClientIPHeader, xRealIPHeader}
	for _, iplHeader := range ipHeaders {
		if clientIP := c.Request().Header.Peek(iplHeader); clientIP != nil {
			return c.SendString(string(clientIP))
		}
	}

	forwardedHeaders := [...]string{xForwardedHeader, forwardedForHeader, forwardedHeader}
	for _, forwardedHeader := range forwardedHeaders {
		if forwarded := c.Request().Header.Peek(forwardedHeader); forwarded != nil {
			return c.SendString(string(forwarded))
		}
	}

	var remoteIP string
	remoteAddr := c.Context().RemoteAddr().String()

	if strings.ContainsRune(remoteAddr, ':') {
		remoteIP, _, _ = net.SplitHostPort(remoteAddr)
	} else {
		remoteIP = remoteAddr
	}
	return c.SendString(remoteIP)
}

func getAllHeaders(c *fiber.Ctx) error {
	headers := c.GetReqHeaders()
	var parsed string
	for k, v := range headers {
		parsed += fmt.Sprintf("%s: %s\n", k, v)
	}
	return c.SendString(parsed)
}

func getAllHeadersJSON(c *fiber.Ctx) error {
	headers := c.GetReqHeaders()
	return c.JSON(headers)
}
