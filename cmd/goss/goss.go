package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aelsabbahy/goss"
	"github.com/aelsabbahy/goss/outputs"
	"github.com/urfave/cli"
)

var version string

func main() {
	startTime := time.Now()
	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Version = version
	app.Name = "goss"
	app.Usage = "Quick and Easy server validation"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "gossfile, g",
			Value:  "./goss.yaml",
			Usage:  "Goss file to read from / write to",
			EnvVar: "GOSS_FILE",
		},
		cli.StringFlag{
			Name:   "vars",
			Usage:  "json/yaml file containing variables for template",
			EnvVar: "GOSS_VARS",
		},
		cli.StringFlag{
			Name:   "vars-inline",
			Usage:  "json/yaml string containing variables for template (overwrites vars)",
			EnvVar: "GOSS_VARS_INLINE",
		},
		cli.StringFlag{
			Name:  "package",
			Usage: "Package type to use [rpm, deb, apk, pacman]",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "validate",
			Aliases: []string{"v"},
			Usage:   "Validate system",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "format, f",
					Value:  "rspecish",
					Usage:  fmt.Sprintf("Format to output in, valid options: %s", outputs.Outputers()),
					EnvVar: "GOSS_FMT",
				},
				cli.StringSliceFlag{
					Name:   "format-options, o",
					Usage:  fmt.Sprintf("Extra options passed to the formatter, valid options: %s", outputs.FormatOptions()),
					EnvVar: "GOSS_FMT_OPTIONS",
				},
				cli.BoolFlag{
					Name:   "color",
					Usage:  "Force color on",
					EnvVar: "GOSS_COLOR",
				},
				cli.BoolFlag{
					Name:   "no-color",
					Usage:  "Force color off",
					EnvVar: "GOSS_NOCOLOR",
				},
				cli.DurationFlag{
					Name:   "sleep,s",
					Usage:  "Time to sleep between retries, only active when -r is set",
					Value:  1 * time.Second,
					EnvVar: "GOSS_SLEEP",
				},
				cli.DurationFlag{
					Name:   "retry-timeout,r",
					Usage:  "Retry on failure so long as elapsed + sleep time is less than this",
					Value:  0,
					EnvVar: "GOSS_RETRY_TIMEOUT",
				},
				cli.IntFlag{
					Name:   "max-concurrent",
					Usage:  "Max number of tests to run concurrently",
					Value:  50,
					EnvVar: "GOSS_MAX_CONCURRENT",
				},
			},
			Action: func(c *cli.Context) error {
				goss.Validate(c, startTime)
				return nil
			},
		},
		{
			Name:    "serve",
			Aliases: []string{"s"},
			Usage:   "Serve a health endpoint",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "format, f",
					Value:  "rspecish",
					Usage:  fmt.Sprintf("Format to output in, valid options: %s", outputs.Outputers()),
					EnvVar: "GOSS_FMT",
				},
				cli.StringSliceFlag{
					Name:   "format-options, o",
					Usage:  fmt.Sprintf("Extra options passed to the formatter, valid options: %s", outputs.FormatOptions()),
					EnvVar: "GOSS_FMT_OPTIONS",
				},
				cli.DurationFlag{
					Name:   "cache,c",
					Usage:  "Time to cache the results",
					Value:  5 * time.Second,
					EnvVar: "GOSS_CACHE",
				},
				cli.StringFlag{
					Name:   "listen-addr,l",
					Value:  ":8080",
					Usage:  "Address to listen on [ip]:port",
					EnvVar: "GOSS_LISTEN",
				},
				cli.StringFlag{
					Name:   "endpoint,e",
					Value:  "/healthz",
					Usage:  "Endpoint to expose",
					EnvVar: "GOSS_ENDPOINT",
				},
				cli.IntFlag{
					Name:   "max-concurrent",
					Usage:  "Max number of tests to run concurrently",
					Value:  50,
					EnvVar: "GOSS_MAX_CONCURRENT",
				},
			},
			Action: func(c *cli.Context) error {
				goss.Serve(c)
				return nil
			},
		},
		{
			Name:    "render",
			Aliases: []string{"r"},
			Usage:   "render gossfile after imports",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "debug, d",
					Usage: fmt.Sprintf("Print debugging info when rendering"),
				},
			},
			Action: func(c *cli.Context) error {
				fmt.Print(goss.RenderJSON(c))
				return nil
			},
		},
		{
			Name:    "autoadd",
			Aliases: []string{"aa"},
			Usage:   "automatically add all matching resource to the test suite",
			Action: func(c *cli.Context) error {
				goss.AutoAddResources(c.GlobalString("gossfile"), c.Args(), c)
				return nil
			},
		},
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "add a resource to the test suite",
			Flags: []cli.Flag{
				cli.StringSliceFlag{
					Name:  "exclude-attr",
					Usage: "Exclude the following attributes when adding a new resource",
				},
			},
			Subcommands: []cli.Command{
				{
					Name:  "package",
					Usage: "add new package",
					Action: func(c *cli.Context) error {
						goss.AddResources(c.GlobalString("gossfile"), "Package", c.Args(), c)
						return nil
					},
				},
				{
					Name:  "file",
					Usage: "add new file",
					Action: func(c *cli.Context) error {
						goss.AddResources(c.GlobalString("gossfile"), "File", c.Args(), c)
						return nil
					},
				},
				{
					Name:  "addr",
					Usage: "add new remote address:port - ex: google.com:80",
					Flags: []cli.Flag{
						cli.DurationFlag{
							Name:  "timeout",
							Value: 500 * time.Millisecond,
						},
					},
					Action: func(c *cli.Context) error {
						goss.AddResources(c.GlobalString("gossfile"), "Addr", c.Args(), c)
						return nil
					},
				},
				{
					Name:  "port",
					Usage: "add new listening [protocol]:port - ex: 80 or udp:123",
					Action: func(c *cli.Context) error {
						goss.AddResources(c.GlobalString("gossfile"), "Port", c.Args(), c)
						return nil
					},
				},
				{
					Name:  "service",
					Usage: "add new service",
					Action: func(c *cli.Context) error {
						goss.AddResources(c.GlobalString("gossfile"), "Service", c.Args(), c)
						return nil
					},
				},
				{
					Name:  "user",
					Usage: "add new user",
					Action: func(c *cli.Context) error {
						goss.AddResources(c.GlobalString("gossfile"), "User", c.Args(), c)
						return nil
					},
				},
				{
					Name:  "group",
					Usage: "add new group",
					Action: func(c *cli.Context) error {
						goss.AddResources(c.GlobalString("gossfile"), "Group", c.Args(), c)
						return nil
					},
				},
				{
					Name:  "command",
					Usage: "add new command",
					Flags: []cli.Flag{
						cli.DurationFlag{
							Name:  "timeout",
							Value: 10 * time.Second,
						},
					},
					Action: func(c *cli.Context) error {
						goss.AddResources(c.GlobalString("gossfile"), "Command", c.Args(), c)
						return nil
					},
				},
				{
					Name:  "dns",
					Usage: "add new dns",
					Flags: []cli.Flag{
						cli.DurationFlag{
							Name:  "timeout",
							Value: 500 * time.Millisecond,
						},
						cli.StringFlag{
							Name:  "server",
							Usage: "The IP address of a DNS server to query",
						},
					},
					Action: func(c *cli.Context) error {
						goss.AddResources(c.GlobalString("gossfile"), "DNS", c.Args(), c)
						return nil
					},
				},
				{
					Name:  "process",
					Usage: "add new process name",
					Action: func(c *cli.Context) error {
						goss.AddResources(c.GlobalString("gossfile"), "Process", c.Args(), c)
						return nil
					},
				},
				{
					Name:  "http",
					Usage: "add new http",
					Flags: []cli.Flag{
						cli.BoolFlag{
							Name: "insecure, k",
						},
						cli.BoolFlag{
							Name: "no-follow-redirects, r",
						},
						cli.DurationFlag{
							Name:  "timeout",
							Value: 5 * time.Second,
						},
						cli.StringFlag{
							Name:  "username, u",
							Usage: "Username for basic auth",
						},
						cli.StringFlag{
							Name:  "password, p",
							Usage: "Password for basic auth",
						},
					},
					Action: func(c *cli.Context) error {
						goss.AddResources(c.GlobalString("gossfile"), "HTTP", c.Args(), c)
						return nil
					},
				},
				{
					Name:  "goss",
					Usage: "add new goss file, it will be imported from this one",
					Action: func(c *cli.Context) error {
						goss.AddResources(c.GlobalString("gossfile"), "Gossfile", c.Args(), c)
						return nil
					},
				},
				{
					Name:  "kernel-param",
					Usage: "add new goss kernel param",
					Action: func(c *cli.Context) error {
						goss.AddResources(c.GlobalString("gossfile"), "KernelParam", c.Args(), c)
						return nil
					},
				},
				{
					Name:  "mount",
					Usage: "add new mount",
					Action: func(c *cli.Context) error {
						goss.AddResources(c.GlobalString("gossfile"), "Mount", c.Args(), c)
						return nil
					},
				},
				{
					Name:  "interface",
					Usage: "add new interface",
					Action: func(c *cli.Context) error {
						goss.AddResources(c.GlobalString("gossfile"), "Interface", c.Args(), c)
						return nil
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
