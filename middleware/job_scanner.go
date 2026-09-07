// middleware/job_scanner.go (100行以下)
package middleware

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func writeScraperLog(line string) {
	_ = os.MkdirAll("logs", 0755)
	f, err := os.OpenFile("logs/scraper.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil { return }
	defer f.Close()
	lvl := "INFO"
	if strings.Contains(line, "[STDERR") || strings.Contains(line, "ERR") { lvl = "ERROR" }
	ts := time.Now().Format("2006-01-02 15:04:05")
	_, _ = f.WriteString(fmt.Sprintf("[%s] [%s] %s\n", ts, lvl, line))
}

var progressRegex = regexp.MustCompile(`^PROGRESS:\s*(\d+)/(\d+)\s*\|\s*(.*)$`)

func (j *JobOrchestrator) scanStdoutProgress(jobID string, r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		j.appendLog(jobID, line)

		matches := progressRegex.FindStringSubmatch(line)
		if len(matches) == 4 {
			cur, _ := strconv.Atoi(matches[1])
			tot, _ := strconv.Atoi(matches[2])
			msg := matches[3]

			var pct float64
			if tot > 0 {
				pct = (float64(cur) / float64(tot)) * 100.0
				if pct > 100.0 {
					pct = 100.0
				}
			}

			j.mu.Lock()
			if p, ok := j.jobs[jobID]; ok {
				p.Current = cur
				p.Total = tot
				p.Percentage = pct
				p.Message = msg
				snapshot := *p
				j.mu.Unlock()
				j.emitEvent("job:progress", &snapshot)
			} else {
				j.mu.Unlock()
			}
		}
	}
	if err := scanner.Err(); err != nil {
		j.appendLog(jobID, "[SCAN_ERR] "+err.Error())
	}
}

func (j *JobOrchestrator) scanStderr(jobID string, r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		j.appendLog(jobID, "[STDERR] "+line)
	}
	if err := scanner.Err(); err != nil {
		j.appendLog(jobID, "[STDERR_SCAN_ERR] "+err.Error())
	}
}

func (j *JobOrchestrator) appendLog(jobID, line string) {
	go writeScraperLog(line)
	j.mu.Lock()
	if p, ok := j.jobs[jobID]; ok {
		p.Logs = append(p.Logs, line)
		if len(p.Logs) > j.maxLogs {
			p.Logs = p.Logs[len(p.Logs)-j.maxLogs:]
		}
	}
	j.mu.Unlock()
	j.emitEvent("job:log", map[string]string{"id": jobID, "line": line})
}

func (j *JobOrchestrator) emitEvent(eventName string, data interface{}) {
	if j.emitter != nil {
		j.emitter(eventName, data)
	}
}
