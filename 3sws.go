package main

import (
	"encoding/base64"
	"fmt"
	tray "github.com/getlantern/systray"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

var icon = "AAABAAEAICAAAAEAIACoEAAAFgAAACgAAAAgAAAAQAAAAAEAIAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACdnZ0AbW1tAFpaWgBSUlIATExMAEdHRwBBQUEAKysrAAAAAAAAAAAAAAAAAAAAAAABAQEBAQEBAgEBAQUBAQEJAAAADgAAABABAQELAQEBBwEBAQMBAQEBAQEBAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAH19fQBAQEAAJycnABoaGgAQEBAACAgIAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEBAQIBAQEGAQEBDgEBARguLi4kPj4+KgEBAR0BAQETAQEBCQEBAQMBAQEBAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAb29vACoqKgALCwsAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABAQEBBQEBARFeXl4sUlJSRsLCwors7OzLW1tbUlpaWjo2NjYbAQEBCQEBAQMBAQEBAAAAAQAAAAEAAAABAAAAAQAAAAAAAAAAAAAAAAAAAABoaGgAHR0dAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQEBAQQBAQENAQEBJtvb26br6+ve+Pj49f7+/v7u7u7n7Ozs1qGhoVMBAQEXAQEBCAEBAQMBAQEDAAAABQAAAAQAAAADAAAAAgEBAQEAAAAAAAAAAFxcXAAEBAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAACAAAABgEBARUyMjI94eHhzP/////5+fn48vLy8v7+/v75+fn3qKioeC0tLSYBAQEPAQEBCgEBAQwAAAARAAAAEAAAAAsAAAAHAQEBBAEBAQIBAQEBOzs7AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABAAAAAgAAAAUAAAALAAAAH8bGxoL29vbx/////5OTk6lLS0uP39/f3/7+/v7j4+PGj4+PRQEBAR4BAQEbLS0tJZ6enkaTk5NBAAAAIgAAABcBAQEOAQEBBgEBAQIAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAIAAAAGAAAADAAAABUAAAAozc3Nj/b29vH/////kpKSp2dnZ5Le3t7d/f39/ebm5syNjY1TZ2dnQXR0dE+NjY1t6+vr0+zs7NaMjIxpfX19TW9vbzEBAQESAQEBCAAAAAAAAAAAAAAAAAAAAAEBAQEBAQEBAgAAAAIAAAADAQEBBwAAABIAAAAjAAAALwAAADobGxtS4uLi0//////29vb16enp6P39/f36+vr4o6OjhltbW1Tb29uz5+fn2ubm5uL+/v7+/v7+/uHh4d7s7Ozf3d3drgEBASoBAQESAAAAAAAAAAABAQEBAQEBAwEBAQYBAQEIAAAACAAAAAoBAQEVS0tLNdPT05ra2tq32tratm9vb2vT09O06Ojo3vj4+Pb+/v7+7Ozs5+np6daQkJBoampqYObm5tv///////////7+/v7+/v7////////////q6urdAQEBQQEBAR0BAQEAAQEBAQEBAQQBAQEKAQEBFAEBARoAAAAZAAAAHQAAADK5ubmM+fn59///////////urq6qV5eXmlra2tstLS0md/f38hra2txenp6WiAgIEOIiIh15eXl4//////29vb4ycnJy8rKysz29vb3/////+np6eZ/f39qAAAALQEBAQEBAQEDAQEBCwEBAR6SkpJLxMTEfEhISEwAAABJOTk5ZOHh4dr////////////////y8vLxdnZ2hwAAAGVVVVVxwMDApp+fn4ABAQFKjIyMW+rq6tb+/v7+/v7+/snJycsiIiJxIiIics/Pz9H+/v7+/v7+/vHx8eIAAABAAQEBAgEBAQgBAQEajIyMVerq6tb7+/v56urq5dvb29bm5ubn/f39/v/////////////////////u7u7w09PT1eLi4uH09PT09/f385qamoR5eXll6+vr2v7+/v7+/v7/ysrKzCIiInIAAABy0NDQ0v7+/v7+/v7+8fHx4gAAAEABAQEDAQEBDVRUVCzm5ubF/v7+/v///////////v7+///////////////////////////////////////9/f3+////////////////9PT07YaGhmyDg4N64eHh3//////29vb3z8/P0dDQ0NL5+fn6/////+Tk5OFxcXFmAAAALAEBAQQBAQENPDw8LNDQ0K3////////////////////////////////////////////////////////////////////////////////s7OziSEhIVnJycl/s7Ozh///////////+/v7+/v7+/v//////////7+/v5TAwMEMBAQEeAQEBBAEBAQ0BAQEofHx8c/r6+vv//////////////////////////////////////////////////////////////////////f39/r+/v7MBAQFFX19fQ9vb27Hq6ure6enp5f7+/v7+/v7+5OTk4e/v7+Xd3d2vKioqKwEBARIBAQEHAQEBFgEBATNzc3N7+fn5+v/////////////////////+/v7+7u7u79vb29nl5eXl/v7+/v/////////////////////9/f3+ubm5sgAAAEgBAQEtAQEBMgEBAUKCgoJp8vLy4/Pz8+Nzc3NiMTExQSsrKyoBAQERAQEBBwEBARMnJycyQUFBXdLS0sz/////////////////////+fn5+r+/v7A6OjpdAAAASgAAAFKUlJSM9vb29v/////////////////////r6+vra2trdSEhIUABAQEnAQEBHgAAACMAAAAwAAAALwAAACABAQEWAQEBDQEBAQYBAQECAQEBKdTU1Kvx8fHu//////////////////////7+/v7S0tLJPj4+UQEBASkBAQEcAQEBIwAAAD+rq6ue/v7+///////////////////////7+/v619fXuZmZmVAAAAAaAAAADwAAABAAAAAPAAAACgEBAQYBAQEEAQEBAQEBAQEAAAA8////////////////////////////////+Pj4+qqqqpcBAQE0AQEBFQEBAQwBAQERAAAAJy4uLmH7+/v7///////////////////////////8/Pz7ycnJigAAACAAAAALAAAABgAAAAQAAAADAAAAAgEBAQEAAAAAAAAAAAAAAD/////////////////////////////////39/f5o6OjkAAAADEBAQEUAQEBCwEBAQ8AAAAlGhoaXff39/j///////////////////////////v7+/rHx8eLAAAAIAAAAAoAAAADAAAAAQAAAAEAAAAAAAAAAAAAAAAAAAAAAAAAMe/v79z4+Pj4/v7+//////////////////39/f7IyMi7Hx8fRwAAACIAAAAXAAAAHQAAADiNjY2I/v7+///////////////////////6+vr67u7u37y8vGsAAAAYAAAABwAAAAIAAAABAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAYUlJSQ6mpqYzl5eXk//////////////////////Dw8PGvr6+YLCwsTAAAADwuLi5FfHx8dOvr6+n/////////////////////9PT09r29vaZdXV1OAAAAIQAAAAwAAAADAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAoAAAAcAAAAPI2NjYv8/Pz9//////////////////////X19fbR0dHKw8PDtMzMzMPt7e3u//////////////////////39/f7GxsbAAAAASwAAACMAAAAOAAAABAAAAAEAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAyMjIAAQEBBQEBAQ8AAAAqenp6dfn5+fr////////////////////////////////9/f3+/////////////////////////////////f39/ru7u60AAAA5AAAAFQAAAAcAAAACAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAADQ0NAElJSQABAQEEAQEBDSoqKivFxcWf////////////////////////////////////////////////////////////////////////////////5eXl1zIyMjwBAQETAQEBBQAAAAEAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAYGBgAT09PAAEBAQMBAQENPT09K+Pj48P////////////////////////////////////////////////////////////////////////////////19fXrkpKSSgEBARMBAQEFAQEBAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABsbGwBRUVEAAQEBAgEBAQkBAQEdsbGxcff39/H/////9PT09Onp6erx8fH0///////////////////////////4+Pj65OTk5+/v7+/+/v7++/v7+c3NzZxjY2MuAQEBDQEBAQQAAAABAAAAAAAAAAAAAAAAAAAAAAAAAAABAQEAHR0dAFJSUgACAgIBAgICBAICAg0BAQElpKSkYtHR0ZuOjo5qYmJiYI2NjYTm5ubk////////////////9/f3+JiYmJBbW1tge3t7YsPDw43Nzc2MODg4MgEBARMBAQEGAQEBAgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAYGBgAfHx8AVFRUAAICAgACAgIBAgICBQICAg0BAQEZUFBQJQEBASAAAAAkAAAAO7+/v538/Pz8///////////S0tK5Ly8vQwAAACYAAAAgLy8vIkZGRiABAQERAQEBBwEBAQIAAAABAAAAAAAAAAAAAAAAAAAAAAAAAAABAQEACwsLACIiIgBWVlYAAQEBAAEBAQACAgIBAgICBAEBAQcBAQEKAQEBCgAAAA0AAAAckZGRT+Xl5cTt7e3c7+/v2pOTk1oAAAAgAQEBDgEBAQoBAQEKAQEBCQEBAQUBAQECAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAABwcHAAwMDAAUFBQAKSkpAFtbWwAAAAAAAAAAAAEBAQABAQEBAQEBAgEBAQMBAQEDAAAABAAAAAoAAAAYKioqLCUlJTU5OTkxAAAAGwEBAQsBAQEFAQEBAwAAAAMBAQECAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAExMTACAgIAAkJCQAJycnAC0tLQBBQUEAbW1tAAAAAAAAAAAAAAAAAAEBAQABAQEAAQEBAQEBAQEAAAABAAAAAwAAAAcAAAANAAAAEAAAAA4AAAAIAQEBAwEBAQEAAAABAAAAAQAAAAEAAAAAAAAAAAAAAAAAAAAAEhISAFBQUABjY2MAaGhoAGpqagBsbGwAcHBwAH5+fgCdnZ0A//AD///wAf//4AAP/8AAA/+AAAD/AAAA/gAAAOAAAADAAAAAgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAwAAAA8AAAAfAAAAPwAAAH8AAAB/AAAA/wAAAP8AAAD/AAAB/4AAAf/AAAP/4AAP//gAH/8="

func main() {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Failed to determine current working directory. Exiting.")
		os.Exit(1)
	}
	// Normalize current working directory
	cwd = filepath.ToSlash(cwd)
	// Set gin mode to release (better performance)
	gin.SetMode(gin.ReleaseMode)
	// Creat new router instance
	r := gin.New()
	// Serve the current working directory with dir listing enabled
	r.Use(static.Serve("/", static.LocalFile(cwd, true)))
	// Convert the base64 encoded .ico to a bytes slices
	bin, err := base64.StdEncoding.DecodeString(icon)
	if err != nil {
		os.Exit(1)
	}
	// Initialize tray instance with onReady & onExit handlers
	tray.Run(
		// onReady handler, setup the tray and configure the server
		func() {
			// Add menu items and set tray icon
			tray.SetIcon(bin)
			status := tray.AddMenuItem("Server starting...", "")
			status.Disable()
			quit := tray.AddMenuItem("Quit", "")

			// Goroutine for handling tray menu clicks
			go func() {
				for {
					select {
					// Trigger onExit handler when quit option clicked
					case <-quit.ClickedCh:
						tray.Quit()
						return
					}
				}
			}()

			// Create a TCP listener on the first available port
			var listener net.Listener
			port := 8080
			for {
				listener, err = net.Listen("tcp", fmt.Sprintf(":%d", port))
				if port > 9000 {
					fmt.Println(
						"No available port in the range 8080-9000. Exiting.",
					)
					os.Exit(1)
				} else if err != nil {
					port++
					continue
				} else {
					break
				}
			}

			fmt.Printf("Server started on port %d.\n", port)
			// Show port number in the tray menu
			status.SetTitle(fmt.Sprintf("Running on port %d", port))
			// Open a browser window with the url of the running instance
			openInBrowser(fmt.Sprintf("http://localhost:%d", port))

			// Run the server on the first available port
			err := http.Serve(listener, r)
			if err != nil {
				fmt.Println("Failed to start 3SWS. Exiting.")
				os.Exit(1)
			}

		},
		// onExit handler, just exit the program gracefully
		func() {
			os.Exit(0)
		},
	)
}

func openInBrowser(url string) error {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	return err
}
