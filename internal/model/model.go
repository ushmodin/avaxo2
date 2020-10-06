package model

type DirItem struct {
	Name     string `json:"name"`
	Size     int64  `json:"size"`
	Modified string `json:"modified"`
	IsDir    bool   `json:"isDir"`
	Error    string `json:"error"`
}

type ProcInfo struct {
	Cmd      string   `json:"cmd"`
	Args     []string `json:"args"`
	Exited   bool     `json:"exited"`
	ExitCode int      `json:"exitCode"`
	Out      []byte   `json:"out"`
	Created  string   `json:"created"`
}

type ProcPsItem struct {
	ID      string   `json:"ID"`
	Cmd     string   `json:"cmd"`
	Args    []string `json:"args"`
	Exited  bool     `json:"exited"`
	Created string   `json:"created"`
}

type ExecRq struct {
	Cmd  string   `json:"cmd"`
	Args []string `json:"args"`
}

type ExecRs struct {
	ProcID string `json:"procId"`
}

type ForwardPacketType int

const (
	ForwardOK    ForwardPacketType = 0
	ForwardInit  ForwardPacketType = 1
	ForwardStart ForwardPacketType = 100
	ForwardBytes ForwardPacketType = 200
)

type ForwardPacketBody struct {
	Str   string `json:"s"`
	Bytes []byte `json:"b"`
}

type ForwardPacket struct {
	Type ForwardPacketType `json:"type"`
	Body ForwardPacketBody `json:"body"`
}
