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
