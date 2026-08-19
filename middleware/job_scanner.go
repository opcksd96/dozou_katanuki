// middleware/job_scanner.go (100行以下)
package middleware

import (
	"bufio"
	"io"
	"regexp"
	"strconv"
)

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
}

func (j *JobOrchestrator) scanStderr(jobID string, r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		j.appendLog(jobID, "[STDERR] "+line)
	}
}

func (j *JobOrchestrator) appendLog(jobID, line string) {
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
