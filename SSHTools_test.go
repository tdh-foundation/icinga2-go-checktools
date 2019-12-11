package icinga2_go_checktools

import (
	"testing"
)

var ict *SSHTools

func Test_NewSSHTools_ArubaController1(t *testing.T) {
	var err error

	ict, err = NewSSHTools("10.10.3.9", "icinga", "", "~/.ssh/id_rsa", 22)
	if err != nil {
		t.Errorf("Error establishing SSH connection to 10.10.3.9 => %s", err)
	}
}

func Test_SendSSHhasPTY_ArubaController1(t *testing.T) {
	err := ict.SendSSHhasPTY([]string{"no paging\r", "show user-table role TDH_Guest\r"}, `(?m)^\(.*\)\s#$`)
	if err != nil {
		t.Errorf("Error sending SSH Command: %s", err)
		return
	}
	if ict.stdOutBuffer == "" {
		t.Errorf("Error command return no response\n")
	} else {
		t.Logf("%s", ict.stdOutBuffer)
	}
}

func Test_SendSSHhasPTY_ArubaController2(t *testing.T) {
	err := ict.SendSSHhasPTY([]string{"no paging\r", "show user-table role TDH_WEB_AUTH\r"}, `(?m)^\(.*\)\s#$`)
	if err != nil {
		t.Errorf("Error sending SSH Command: %s", err)
		return
	}
	if ict.stdOutBuffer == "" {
		t.Errorf("Error command return no response\n")
	} else {
		t.Logf("%s", ict.stdOutBuffer)
	}
}

func Test_NewSSHConnection_CISCO2960(t *testing.T) {
	var err error

	ict, err = NewSSHTools("10.10.100.31", "icinga", "", "~/.ssh/id_rsa", 22)
	if err != nil {
		t.Errorf("Error establishing SSH connection to 10.10.100.31 => %s", err)
	}
}

func Test_SendSSH_CISCO2960(t *testing.T) {
	err := ict.SendSSH("show interface status")
	if err != nil {
		t.Errorf("Error sending SSH Command: %s", err)
		return
	}
	if ict.stdOutBuffer == "" {
		t.Errorf("Error command return no response\n")
	} else {
		t.Logf("%s", ict.stdOutBuffer)
	}
}
