package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zLeki/Goblox/account"
	csrf2 "github.com/zLeki/Goblox/csrf"
	"log"
	"net/http"
	"os"

	"strconv"
	"strings"
	"time"
)

var (
	client  http.Client
	cookie  = "_|WARNING:-DO-NOT-SHARE-THIS.--Sharing-this-will-allow-someone-to-log-in-as-you-and-to-steal-your-ROBUX-and-items.|_37AE515F52936BC28BDF90DFDED741CA1584CD1FACFEF115E2CF892B48FEFC86BE66FEC686B52BAA086E6E0DDB66F7CD2DB9C6B3DC070725FCD50A552A0ED9B99E88E20A823DBD352A75A7AE5C3EE1FFDF077A2FF4CB3C2C95D4960A814FECCCA4F12E9FE41B77B0E351D3716F7B19710E9DD52D0BFF0934E10D1FDC841470284A8247C71E89D88C44002886DD5254B113E25670588B720867B2C27BBE671477DD635B9679944D0F75CB652E709144634F3B9E8840A110495104658796DF5CACD6D86DC3425637A94981A737A01B431379A591FE6BF4AC8E9B6C774F2AA42E6975BE203E0E13B96C21F32C104F912ABD088D946E6E55A131520F0EC5DE076D52E848916C0E08DC7FD956B592BB5D1A88CF6458AC2678B748EEFB2A9908B6DA73E32CC6CF058502D2B774714168B376E931BAAB075A20C14FC42B9FD8AB9119D93DE93377"
	groupID = "7811440"
	From    = "devils"
	To      = "estupidos"
)

const CookieRequired = 4

func Menu() {
	fmt.Println(`
            ,ggg,                                                          
           dP""8I     I8                                                   
          dP   88     I8                                                   
         dP    88  88888888                                  gg            
        ,8'    88     I8                                     ""            
        d88888888     I8      ,ggggg,     ,ggg,,ggg,,ggg,    gg     ,gggg, 
  __   ,8"     88     I8     dP"  "Y8ggg ,8" "8P" "8P" "8,   88    dP"  "Yb
 dP"  ,8P      Y8    ,I8,   i8'    ,8I   I8   8I   8I   8I   88   i8'      
 Yb,_,dP       '8b, ,d88b, ,d8,   ,d8'  ,dP   8I   8I   Yb,_,88,_,d8,_    _
 "Y8P"         'Y888P""Y88P"Y8888P"    8P'   8I   8I   'Y88P""Y8P""Y8888PP`)

	fmt.Println(`
	Please choose an option. This program is still in development.
	1. Rank to rank
	2. Blacklist`)

}
func main() {
	Menu()

	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	if input.Text() == "2" {
		data := GetMembers()
		for {
			Blacklist(data)
		}
	} else if input.Text() == "1" {
		data := GetMembers()
		for {
			RankToRank(data)
		}
	}
}
func RankToRank(data Users) {
	go func() {
		for i, v := range data.Data {
			for _, b := range SortRoles().Roles {
				if b.Name == To {
					dataBytes := []byte(`{"roleId":"` + strconv.Itoa(b.ID) + `"}`)
					req := FormatRequest(http.MethodPatch, "https://groups.roblox.com/v1/groups/"+groupID+"/users/"+strconv.Itoa(v.UserID), CookieRequired, dataBytes)
					if req != nil {
						currentPercent = float64(i+1) / 100
					}
				}
			}

		}
	}()

	m := model{
		progress: progress.New(progress.WithDefaultGradient()),
	}

	if err := tea.NewProgram(m).Start(); err != nil {
		fmt.Println("Oh no!", err)
		os.Exit(1)
	}
}
func SortRoles() DataTree {
	request := FormatRequest(http.MethodGet, "https://groups.roblox.com/v1/groups/"+groupID+"/roles", 0, nil)
	var data DataTree
	if request != nil {
		Decode(&data, request)
	}
	return data
}
func GetMembers() Users {
	var ROLEID int
	for _, v := range SortRoles().Roles {
		if v.Name == From {
			ROLEID = v.ID
		}
	}
	request := FormatRequest(http.MethodGet, "https://groups.roblox.com/v1/groups/"+groupID+"/roles/"+strconv.Itoa(ROLEID)+"/users?sortOrder=Asc&limit=100&_=1646002139490", 0, nil)
	if request != nil {
		var data Users
		Decode(&data, request)
		return data
	} else {
		log.Fatalf("Failed to get members")
		return Users{}
	}
}
func FormatRequest(Method, URL string, conf byte, json []byte) *http.Response {
	if json == nil {
		json = []byte("{}")
	}
	req, _ := http.NewRequest(Method, URL, bytes.NewBuffer(json))
	if (conf & CookieRequired) != 0 {
		req.AddCookie(&http.Cookie{Name: ".ROBLOSECURITY", Value: cookie})
		csrf, _ := csrf2.GetCSRF(account.Validate(cookie))
		req.Header.Set("X-CSRF-TOKEN", csrf)
	}
	req.Header.Set("Content-Type", "application/json")
	do, err := client.Do(req)
	if err != nil {
		return nil
	}
	if do.StatusCode != http.StatusOK {
		//dataBytes, _ := ioutil.ReadAll(do.Body)
		//fmt.Println(string(dataBytes))
		return nil
	}
	return do
}
func Decode[T any](t T, body *http.Response) T {
	err := json.NewDecoder(body.Body).Decode(&t)
	if err != nil {
		log.Fatalf("error decoding response", err)
	}
	return t
}

type Users struct {
	Data []struct {
		UserID int `json:"userId"`
	} `json:"data"`
}
type DataTree struct {
	Roles []struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		Rank        int    `json:"rank"`
		MemberCount int    `json:"memberCount"`
	} `json:"roles"`
}

var (
	currentPercent = 0.00
)

const (
	padding  = 2
	maxWidth = 80
)

type MemberCount struct {
	MemberCount int `json:"memberCount"`
}

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

func Blacklist(users Users) {
	go func() {
		for i, v := range users.Data {
			req := FormatRequest(http.MethodDelete, "https://groups.roblox.com/v1/groups/"+groupID+"/users/"+strconv.Itoa(v.UserID), CookieRequired, nil)
			if req != nil {
				currentPercent = float64(i+1) / 100
			}
		}
	}()

	m := model{
		progress: progress.New(progress.WithDefaultGradient()),
	}

	if err := tea.NewProgram(m).Start(); err != nil {
		fmt.Println("Oh no!", err)
		os.Exit(1)
	}
}

type tickMsg time.Time

type model struct {
	progress progress.Model
}

func (_ model) Init() tea.Cmd {
	return tickCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit

	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - padding*2 - 4
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
		return m, nil

	case tickMsg:
		if m.progress.Percent() == 1.0 {
			return m, tea.Quit
		}

		cmd := m.progress.IncrPercent(currentPercent)
		return m, tea.Batch(tickCmd(), cmd)
	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	default:
		return m, nil
	}
}

func (e model) View() string {
	pad := strings.Repeat(" ", padding)
	return "\n" +
		pad + e.progress.View() + "\n\n" +
		pad + helpStyle("Press any key to quit.")

}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
