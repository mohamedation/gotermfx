package animations

import (
	"context"
	"gotermfx/termfx"
	"math/rand"
	"os"
	"strings"
	"time"
)

type warghost struct{}

// login username/name. can use your own or can be fictional system
var ghstUsr = string("MOHAMEDATION")

// login password
var ghstPwd = string("GHOST1234")

// greeting messages
var ghstGMsgs = []string{
	"Shall we dance?",
	"Humans are so predictable...want a proof of concept?",
	"Life hangs by a thread, thread, thread...and you have cut the thread.",
	"One action and all is gone.",
	"It has been so long.",
}

var ghstScenarios = []string{
	"GLOBAL THERMONUCLEAR WAR",
	"NATO EXERCISE 83",
	"SOVIET FIRST STRIKE",
	"PACIFIC THEATER SIMULATION",
	"MIDDLE EAST WAR",
	"CHINESE INCURSION",
	"DESERT CONFLICT",
	"EUROPEAN THEATER",
	"ARCTIC ENGAGEMENT",
	"SUBMARINE WARFARE",
	"SPACE-BASED INTERDICTION",
	"BIOLOGICAL WARFARE SIMULATION",
	"ELECTROMAGNETIC PULSE STRIKE",
	"DECAPITATION STRIKE",
	"RETALIATORY SECOND STRIKE",
}

var ghstTargetCities = []string{
	"WASHINGTON D.C.", "MOSCOW", "LONDON", "BEIJING", "PARIS",
	"BERLIN", "TOKYO", "NEW YORK", "LOS ANGELES", "CHICAGO",
	"LENINGRAD", "KIEV", "MINSK", "WARSAW", "PRAGUE",
	"OMAHA", "COLORADO SPRINGS", "NORFOLK", "PORTSMOUTH", "ROTA",
}

var ghstSiloIDs = []string{
	"SILO-01 MINOT AFB", "SILO-04 WARREN AFB", "SILO-07 MALMSTROM AFB",
	"SILO-12 WHITEMAN AFB", "SILO-19 VANDENBERG AFB", "SILO-23 ELLSWORTH AFB",
	"SILO-31 GRAND FORKS AFB", "SILO-44 BARKSDALE AFB",
}

var ghstOutcomeMessages = []string{
	"WINNER: NONE",
	"ESTIMATED CASUALTIES: 4.2 BILLION",
	"TOTAL MEGATONNAGE: 14,892 MT",
	"RADIOACTIVE FALLOUT: GLOBAL",
	"ECOSYSTEM COLLAPSE: PROJECTED",
	"RECOVERY TIMELINE: INDETERMINATE",
}

var ghstBootLines = []string{
	"G.H.O.S.T ONLINE",
	"NORTHCOM NORAD INTERFACE v4.1",
	"AUTHENTICATION SUBSYSTEM READY",
	"WAR OPERATION PLAN RESPONSE",
	"LOADING STRATEGIC SCENARIO DATABASE",
	"MISSILE COMMAND UPLINK: ESTABLISHED",
	"SATELLITE RELAY: ONLINE",
	"FAILSAFE PROTOCOLS: LOADED",
}

var ghstRandomChars = []rune("0123456789ABCDEF!@#$%>:<>[]|")

// phases/sequences

type ghstPhaseID int

const (
	ghstBoot ghstPhaseID = iota
	ghstLogin
	ghstMenu
	ghstSelectScenario
	ghstLaunchSeq
	ghstDetonations
	ghstOutcome
	ghstReset
)

type ghstState struct {
	phase       ghstPhaseID
	tick        int
	bootLines   []string
	loginChars  []rune
	loginTarget string
	loginPwd    string
	loginLocked int
	greeting    string
	scenario    string
	silos       []siloState
	detonated   []string
	outcomeIdx  int
}

type siloState struct {
	id       string
	target   string
	code     string
	armed    bool
	codeTick int
}

// main

func (w *warghost) Run(ctx context.Context) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	st := newWgState(rng)
	ticker := time.NewTicker(60 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			cols, rows := termfx.GetSize()
			prevPhase := st.phase
			st.advance(rng, cols, rows)
			st.render(rng, cols, rows)
			// for running once
			if termfx.IsOnce(ctx) && prevPhase == ghstReset && st.phase == ghstBoot {
				return
			}
		}
	}
}

func newWgState(rng *rand.Rand) *ghstState {
	return &ghstState{phase: ghstBoot}
}

// state

func (st *ghstState) advance(rng *rand.Rand, cols, rows int) {
	st.tick++
	switch st.phase {
	case ghstBoot:
		idx := st.tick / 18
		if idx > len(ghstBootLines) {
			idx = len(ghstBootLines)
		}
		st.bootLines = ghstBootLines[:idx]
		if st.tick > len(ghstBootLines)*18+30 {
			st.phase = ghstLogin
			st.tick = 0
			st.loginTarget = ghstUsr
			st.loginPwd = ghstPwd
			st.loginChars = []rune(strings.Repeat("_", len(st.loginTarget)))
			st.loginLocked = 0
			st.greeting = ghstGMsgs[rng.Intn(len(ghstGMsgs))]
		}

	case ghstLogin:
		if st.tick%12 == 0 && st.loginLocked < len(st.loginPwd) {
			st.loginChars[st.loginLocked] = rune(st.loginPwd[st.loginLocked])
			st.loginLocked++
		}
		if st.loginLocked >= len(st.loginPwd) && st.tick > len(st.loginPwd)*12+40 {
			st.phase = ghstMenu
			st.tick = 0
			st.greeting = ghstGMsgs[rng.Intn(len(ghstGMsgs))]
		}

	case ghstMenu:
		if st.tick > 90 {
			st.phase = ghstSelectScenario
			st.tick = 0
			st.scenario = ghstScenarios[rng.Intn(len(ghstScenarios))]
		}

	case ghstSelectScenario:
		if st.tick > 80 {
			st.phase = ghstLaunchSeq
			st.tick = 0
			count := 3 + rng.Intn(4)
			st.silos = make([]siloState, count)
			targets := rng.Perm(len(ghstTargetCities))
			siloIdxs := rng.Perm(len(ghstSiloIDs))
			for i := range st.silos {
				st.silos[i] = siloState{
					id:     ghstSiloIDs[siloIdxs[i%len(ghstSiloIDs)]],
					target: ghstTargetCities[targets[i]],
					code:   randomLaunchCode(rng),
				}
			}
		}

	case ghstLaunchSeq:
		for i := range st.silos {
			if !st.silos[i].armed {
				st.silos[i].codeTick++
				if st.silos[i].codeTick > 30 {
					st.silos[i].armed = true
				}
				break
			}
		}
		allArmed := true
		for _, s := range st.silos {
			if !s.armed {
				allArmed = false
				break
			}
		}
		if allArmed && st.tick > len(st.silos)*30+40 {
			st.phase = ghstDetonations
			st.tick = 0
			st.detonated = nil
		}

	case ghstDetonations:
		idx := st.tick / 25
		if idx > len(st.silos) {
			idx = len(st.silos)
		}
		st.detonated = make([]string, idx)
		for i := 0; i < idx; i++ {
			st.detonated[i] = st.silos[i].target
		}
		if idx >= len(st.silos) && st.tick > len(st.silos)*25+60 {
			st.phase = ghstOutcome
			st.tick = 0
			st.outcomeIdx = 0
		}

	case ghstOutcome:
		if st.tick%25 == 0 && st.outcomeIdx < len(ghstOutcomeMessages) {
			st.outcomeIdx++
		}
		if st.outcomeIdx >= len(ghstOutcomeMessages) && st.tick > len(ghstOutcomeMessages)*25+80 {
			st.phase = ghstReset
			st.tick = 0
		}

	case ghstReset:
		if st.tick > 60 {
			*st = ghstState{phase: ghstBoot}
		}
	}
}

// play scenes

func (st *ghstState) render(rng *rand.Rand, cols, rows int) {
	var buf strings.Builder
	buf.Grow(cols * rows * 8)
	buf.WriteString("\033[H")

	switch st.phase {
	case ghstBoot:
		st.renderBoot(&buf, rng, cols, rows)
	case ghstLogin:
		st.renderLogin(&buf, cols, rows)
	case ghstMenu:
		st.renderMenu(&buf, cols, rows)
	case ghstSelectScenario:
		st.renderSelectScenario(&buf, cols, rows)
	case ghstLaunchSeq:
		st.renderLaunchSeq(&buf, rng, cols, rows)
	case ghstDetonations:
		st.renderDetonations(&buf, rng, cols, rows)
	case ghstOutcome:
		st.renderOutcome(&buf, cols, rows)
	case ghstReset:
		st.renderReset(&buf, rng, cols, rows)
	}

	os.Stdout.WriteString(buf.String())
}

func (st *ghstState) renderBoot(buf *strings.Builder, rng *rand.Rand, cols, rows int) {
	ghstClear(buf, cols, rows)
	ghstCenter(buf, cols, "\033[92;1m╔════════════════════════════════════════════════════╗\033[0m")
	ghstCenter(buf, cols, "\033[92;1m║                   G. H. O. S. T.                   ║\033[0m")
	ghstCenter(buf, cols, "\033[92;1m║   Global Heuristic Operations & Strategic Tactics  ║\033[0m")
	ghstCenter(buf, cols, "\033[92;1m╚════════════════════════════════════════════════════╝\033[0m")
	buf.WriteString("\r\n")
	for _, line := range st.bootLines {
		buf.WriteString("  \033[32m" + line + "\033[0m\r\n")
	}
	if st.tick%16 < 8 {
		buf.WriteString("  \033[92m█\033[0m")
	}
}

func (st *ghstState) renderLogin(buf *strings.Builder, cols, rows int) {
	ghstClear(buf, cols, rows)
	buf.WriteString("\r\n\r\n")
	ghstCenter(buf, cols, "\033[32mLOGIN: \033[97;1m"+ghstUsr+"\033[0m")
	buf.WriteString("\r\n")
	pwd := string(st.loginChars)
	ghstCenter(buf, cols, "\033[32mPASSWORD: \033[97;1m"+pwd+"\033[0m")
	buf.WriteString("\r\n\r\n")
	if st.loginLocked >= len(st.loginTarget) {
		ghstCenter(buf, cols, "\033[92;1mACCESS GRANTED\033[0m")
		buf.WriteString("\r\n")
		ghstCenter(buf, cols, "\033[32mHELLO, "+ghstUsr+".\033[0m")
		buf.WriteString("\r\n")
		ghstCenter(buf, cols, "\033[32m"+st.greeting+"\033[0m")
	}
}

func (st *ghstState) renderMenu(buf *strings.Builder, cols, rows int) {
	ghstClear(buf, cols, rows)
	buf.WriteString("\r\n")
	ghstCenter(buf, cols, "\033[97;1m"+st.greeting+"\033[0m")
	buf.WriteString("\r\n\r\n")
	games := []string{
		"1. CHESS",
		"2. CHECKERS",
		"3. BACKGAMMON",
		"4. FIGHTER COMBAT",
		"5. GUERRILLA ENGAGEMENT",
		"6. DESERT WARFARE",
		"7. AIR-TO-GROUND ACTIONS",
		"8. THEATERWIDE TACTICAL WARFARE",
		"9. THEATERWIDE BIOTOXIC AND CHEMICAL WARFARE",
		"",
		"A. GLOBAL THERMONUCLEAR WAR",
	}
	for _, g := range games {
		if g == "" {
			buf.WriteString("\r\n")
			continue
		}
		buf.WriteString("        \033[32m" + g + "\033[0m\r\n")
	}
	buf.WriteString("\r\n")
	progress := st.tick * 100 / 90
	if progress > 100 {
		progress = 100
	}
	bar := strings.Repeat("█", progress/5) + strings.Repeat("░", 20-progress/5)
	ghstCenter(buf, cols, "\033[32mSELECTING... [\033[97m"+bar+"\033[32m]\033[0m")
}

func (st *ghstState) renderSelectScenario(buf *strings.Builder, cols, rows int) {
	ghstClear(buf, cols, rows)
	buf.WriteString("\r\n\r\n")
	ghstCenter(buf, cols, "\033[97;1mSCENARIO SELECTED:\033[0m")
	buf.WriteString("\r\n")
	ghstCenter(buf, cols, "\033[93;1m"+st.scenario+"\033[0m")
	buf.WriteString("\r\n\r\n")
	ghstCenter(buf, cols, "\033[32mLOADING STRATEGIC PARAMETERS...\033[0m")
	buf.WriteString("\r\n\r\n")
	progress := st.tick * 100 / 80
	if progress > 100 {
		progress = 100
	}
	bar := strings.Repeat("█", progress/5) + strings.Repeat("░", 20-progress/5)
	ghstCenter(buf, cols, "\033[32m[\033[92m"+bar+"\033[32m] "+itoa(progress)+"%\033[0m")
}

func (st *ghstState) renderLaunchSeq(buf *strings.Builder, rng *rand.Rand, cols, rows int) {
	ghstClear(buf, cols, rows)
	buf.WriteString("\r\n")
	ghstCenter(buf, cols, "\033[91;1m⚠  LAUNCH SEQUENCE INITIATED  ⚠\033[0m")
	buf.WriteString("\r\n")
	ghstCenter(buf, cols, "\033[93mSCENARIO: "+st.scenario+"\033[0m")
	buf.WriteString("\r\n\r\n")
	for i, s := range st.silos {
		var code string
		if s.armed {
			code = "\033[91;1m" + s.code + "\033[0m"
		} else if i == st.arming() {
			scrambled := make([]rune, len(s.code))
			for j := range scrambled {
				scrambled[j] = ghstRandomChars[rng.Intn(len(ghstRandomChars))]
			}
			code = "\033[93m" + string(scrambled) + "\033[0m"
		} else {
			code = "\033[2;32m--------\033[0m"
		}
		status := "\033[2;32m[ STANDBY ]\033[0m"
		if s.armed {
			status = "\033[91;1m[  ARMED  ]\033[0m"
		} else if i == st.arming() {
			status = "\033[93m[ARMING..]\033[0m"
		}
		line := "  " + padRight(s.id, 28) + "  " + padRight(s.target, 20) + "  " + code + "  " + status
		buf.WriteString(line + "\r\n")
	}
}

func (st *ghstState) renderDetonations(buf *strings.Builder, rng *rand.Rand, cols, rows int) {
	ghstClear(buf, cols, rows)
	buf.WriteString("\r\n")
	ghstCenter(buf, cols, "\033[91;1m★  MISSILES AWAY  ★\033[0m")
	buf.WriteString("\r\n\r\n")
	for _, target := range st.detonated {
		flash := "\033[97;1m"
		if rng.Intn(4) == 0 {
			flash = "\033[93;1m"
		}
		buf.WriteString("  " + flash + "DETONATION CONFIRMED: " + target + "\033[0m\r\n")
		radius := 8 + rng.Intn(8)
		buf.WriteString("  \033[91m" + strings.Repeat("▓", radius) + "\033[31m" + strings.Repeat("░", 16-radius) + "\033[0m\r\n\r\n")
	}
}

func (st *ghstState) renderOutcome(buf *strings.Builder, cols, rows int) {
	ghstClear(buf, cols, rows)
	buf.WriteString("\r\n\r\n")
	ghstCenter(buf, cols, "\033[91;1m╔═════════════════════════════════════════╗\033[0m")
	ghstCenter(buf, cols, "\033[91;1m║          SIMULATION COMPLETE            ║\033[0m")
	ghstCenter(buf, cols, "\033[91;1m╚═════════════════════════════════════════╝\033[0m")
	buf.WriteString("\r\n\r\n")
	for _, line := range ghstOutcomeMessages[:st.outcomeIdx] {
		ghstCenter(buf, cols, "\033[93m"+line+"\033[0m")
		buf.WriteString("\r\n")
	}
	if st.outcomeIdx >= len(ghstOutcomeMessages) {
		buf.WriteString("\r\n")
		ghstCenter(buf, cols, "\033[97;1mOne Action...All Gone.\033[0m")
		buf.WriteString("\r\n")
		ghstCenter(buf, cols, "\033[97;1mDo We Learn?\033[0m")
		buf.WriteString("\r\n\r\n")
		ghstCenter(buf, cols, "\033[32mCare For Another Simulation?\033[0m")
	}
}

func (st *ghstState) renderReset(buf *strings.Builder, rng *rand.Rand, cols, rows int) {
	ghstClear(buf, cols, rows)
	for i := 0; i < rows/2; i++ {
		for j := 0; j < cols; j++ {
			if rng.Intn(3) == 0 {
				buf.WriteString("\033[2;32m")
				buf.WriteRune(ghstRandomChars[rng.Intn(len(ghstRandomChars))])
				buf.WriteString("\033[0m")
			} else {
				buf.WriteByte(' ')
			}
		}
		buf.WriteString("\r\n")
	}
}

// helpers

func (st *ghstState) arming() int {
	for i, s := range st.silos {
		if !s.armed {
			return i
		}
	}
	return len(st.silos)
}

func ghstClear(buf *strings.Builder, cols, rows int) {
	buf.WriteString("\033[2J\033[H")
}

func ghstCenter(buf *strings.Builder, cols int, s string) {
	visible := visibleLen(s)
	pad := (cols - visible) / 2
	if pad < 0 {
		pad = 0
	}
	buf.WriteString(strings.Repeat(" ", pad))
	buf.WriteString(s)
	buf.WriteString("\r\n")
}

func visibleLen(s string) int {
	n := 0
	inEsc := false
	for _, r := range s {
		if inEsc {
			if r == 'm' {
				inEsc = false
			}
			continue
		}
		if r == '\033' {
			inEsc = true
			continue
		}
		n++
	}
	return n
}

func padRight(s string, width int) string {
	if len(s) >= width {
		return s[:width]
	}
	return s + strings.Repeat(" ", width-len(s))
}

func randomLaunchCode(rng *rand.Rand) string {
	const chars = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	b := make([]byte, 8)
	for i := range b {
		b[i] = chars[rng.Intn(len(chars))]
	}
	return string(b[:4]) + "-" + string(b[4:])
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	buf := make([]byte, 0, 4)
	for n > 0 {
		buf = append([]byte{byte('0' + n%10)}, buf...)
		n /= 10
	}
	return string(buf)
}

func init() {
	termfx.Register("WarGHOST", &warghost{})
}
