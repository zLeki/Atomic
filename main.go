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
	cookie  = ""
	groupID = ""
	From    = ""
	To      = ""
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

	req := FormatRequest(http.MethodPatch, "https://groups.roblox.com/v1/groups/"+groupID+"/status", CookieRequired, []byte(`{"message":"get nuked lol!!"}`))
	if req == nil {
		fmt.Println("Failed to post a shout, maybe theres captcha? Continuing...")
	}
	Check := func(data Users) {
		if len(data.Data) == 0 {
			fmt.Println("Finished")
			time.Sleep(time.Second * 3)
			os.Exit(0)
		}
	}
	switch input.Text() {
	case "2":
		for {
			data := GetMembers()
			Check(data)
			Blacklist(data)
		}
	case "1":
		for {
			data := GetMembers()
			Check(data)
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

	if err := tea.NewProgram(&m).Start(); err != nil {
		fmt.Println("Oh no!", err)
		os.Exit(1)
	}
}
func SortRoles() DataTree {
	request := FormatRequest(http.MethodGet, "https://groups.roblox.com/v1/groups/"+groupID+"/roles", 0, nil)
	var data DataTree
	if request != nil {
		err := json.NewDecoder(request.Body).Decode(&data)
		if err != nil {
			log.Fatalf("error decoding response %d", err)
		}
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
		err := json.NewDecoder(request.Body).Decode(&data)
		if err != nil {
			log.Fatalf("error decoding response %d", err)
		}
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

	if err := tea.NewProgram(&m).Start(); err != nil {
		fmt.Println("Oh no!", err)
		os.Exit(1)
	}
}

type tickMsg time.Time

type model struct {
	progress progress.Model
}

func (_ *model) Init() tea.Cmd {
	return tickCmd()
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (e *model) View() string {
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
