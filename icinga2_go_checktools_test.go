package icinga2_go_checktools

import (
	"testing"
)

var ict *SSHTools

func Test_NewSSHTools1(t *testing.T) {
	var err error

	ict, err = NewSSHTools("10.10.3.9", "icinga", "", "~/.ssh/id_rsa", 22)
	if err != nil {
		t.Errorf("Error establishing SSH connection to 10.10.3.9 => %s", err)
	}
}

func Test_SendSSHhasPTY_ArubaController1(t *testing.T) {
	ict.SendSSHhasPTY([]string{"no paging\r", "show user-table role TDH_Guest\r"}, `(?m)^\(.*\)\s#$`)
	t.Logf("%s", ict.stdOutBuffer)
}

func Test_SendSSHhasPTY_ArubaController2(t *testing.T) {
	ict.SendSSHhasPTY([]string{"no paging\r", "show user-table role TDH_WEB_AUTH\r"}, `(?m)^\(.*\)\s#$`)
	t.Logf("%s", ict.stdOutBuffer)
}

func Test_NewSSHConnection2(t *testing.T) {
	var err error

	ict, err = NewSSHTools("10.10.100.31", "icinga", "", "~/.ssh/id_rsa", 22)
	if err != nil {
		t.Errorf("Error establishing SSH connection to 10.10.100.31 => %s", err)
	}
}

func Test_SendSSH_CISCO2960(t *testing.T) {
	ict.SendSSH("show interface status")
	t.Logf("%s", ict.stdOutBuffer)
}
