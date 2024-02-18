package ip_parse

import (
	"net"

	"github.com/gofiber/fiber/v2"
)

func New() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		cloudflareRealIPHeader := c.Get("CF-Connecting-IP")

		if len(cloudflareRealIPHeader) > 0 {
			c.Locals("remote_ip", net.ParseIP(cloudflareRealIPHeader))
			return c.Next()
		}

		for _, ip := range c.IPs() {
			remoteIP := net.ParseIP(ip)

			if remoteIP.IsPrivate() {
				continue
			}

			c.Locals("remote_ip", net.ParseIP(ip))
		}

		if c.Locals("remote_ip") == nil {
			c.Locals("remote_ip", net.ParseIP(c.IP()))
		}

		return c.Next()
	}
}
