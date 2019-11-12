package qemu

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/digitalocean/go-qemu/qmp"
	"github.com/digitalocean/go-qemu/qmp/raw"
	"github.com/shirou/gopsutil/process"
)

func (q *Qemu) init() error {
	if ok, err := q.findProcess(q.name); err != nil {
		return fmt.Errorf("Failed to find process: err='%s'", err.Error())
	} else if ok {
		c, err := q.proc.CmdlineSlice()
		if err != nil {
			return fmt.Errorf("Failed to get command line: err='%s'", err.Error())
		}

		if err := q.parseArgs(c); err != nil {
			return fmt.Errorf("Failed to parse arguments of command line: err='%s'", err.Error())
		}

		if err := q.initQMP(); err != nil {
			return fmt.Errorf("Failed to initialize QMP socket: err='%s'", err.Error())
		}
	}

	return nil
}

// Want test!!!
func (q *Qemu) findProcess(contain string) (bool, error) {
	ps, err := process.Processes()
	if err != nil {
		return false, fmt.Errorf("Failed to get process list")
	}

	for _, p := range ps {
		c, _ := p.Cmdline() // エラーが発生する場合が考えられない
		if strings.HasPrefix(c, "qemu") {
			a := ParseQemuArgs(strings.Split(c, " "))
			_, kwds, ok := a.GetTopParsedOptionValue("-name")
			if ok && kwds["guest"] == q.name {
				q.proc = p

				return true, nil
			}
		}
	}

	return false, nil
}

func (q *Qemu) initQMP() error {
	b := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 5) // CIがこけてしまうため長くしているが、かなり長時間待つようになっているため普通の処理においてタイムアウトする危険性がある
	err := backoff.Retry(func() (err error) {
		q.qmp, err = qmp.NewSocketMonitor("unix", q.qmpPath, 3*time.Second)
		return
	}, b)
	if err != nil {
		return fmt.Errorf("Failed to open QMP socket: err='%s'", err.Error())
	}

	if err := q.qmp.Connect(); err != nil {
		return fmt.Errorf("Failed to connect QMP socket: err='%s'", err.Error())
	}

	q.m = raw.NewMonitor(q.qmp)

	return nil
}

func (q *Qemu) parseArgs(args []string) error {
	q.args = ParseQemuArgs(args)

	_, kwds, ok := q.args.GetTopParsedOptionValue("-mon")
	if ok {
		if v, ok := kwds["chardev"]; ok {
			_, kwds, _ := q.args.GetParsedOptionValueById("-chardev", v)
			q.qmpPath = kwds["path"]
		}
	}

	return nil
}

func GetNewListenPort(begin uint16) uint16 {
	for p := begin; p <= uint16(65535); p++ {
		l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", p))
		if err == nil {
			defer l.Close()

			return p
		}
	}

	return 0
}

// QemuArgs is QEMU arguments parser
// TODO: 毎回文字列探索するので遅い
type QemuArgs struct {
	args []string
}

func ParseQemuArgs(args []string) *QemuArgs {
	return &QemuArgs{
		args: args,
	}
}

func (q QemuArgs) GetArgs() []string {
	return q.args
}

func (q QemuArgs) GetOptionValues(key string) []string {
	values := make([]string, 0)

	for i, a := range q.args {
		if a == key {
			values = append(values, q.args[i+1])
		}
	}

	return values
}

func (q QemuArgs) ParseOptionValue(value string) ([]string, map[string]string) {
	args := make([]string, 0)
	kwds := make(map[string]string)

	for _, v := range strings.Split(value, ",") {
		kv := strings.Split(v, "=")
		if len(kv) == 1 {
			args = append(args, v)
		} else {
			kwds[kv[0]] = kv[1]
		}
	}

	return args, kwds
}

func (q QemuArgs) GetTopParsedOptionValue(key string) ([]string, map[string]string, bool) {
	values := q.GetOptionValues(key)
	if len(values) == 0 {
		return nil, nil, false
	}

	args, kwds := q.ParseOptionValue(values[0])
	return args, kwds, true
}

func (q QemuArgs) GetParsedOptionValueById(key, id string) ([]string, map[string]string, bool) {
	values := q.GetOptionValues(key)

	for _, value := range values {
		args, kwds := q.ParseOptionValue(value)

		if kwds["id"] == id {
			return args, kwds, true
		}
	}

	return nil, nil, false
}
