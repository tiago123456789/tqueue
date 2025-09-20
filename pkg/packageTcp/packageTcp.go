package packagetcp

import (
	"bytes"
	"errors"
	"strings"

	"github.com/tiago123456789/tqueue/pkg/instruction"
)

const INCOMPLETE = -1
const COMPLETE = 0

func ParseMessage(previousBuf []byte, buf []byte) ([]string, int) {
	items := []string{}
	if string(previousBuf) != "" && bytes.Contains(buf, []byte("\n")) {
		partMessage := string(previousBuf) + string(buf)
		startPostion := 0
		for i := 0; i < len(partMessage); i++ {
			if partMessage[i] == '\n' {
				items = append(items, string(partMessage[startPostion:i]))
				startPostion = i + 1
			}
		}
		if len(items) == 0 {
			items = append(items, partMessage)
		}
		partMessage = ""
		return items, COMPLETE
	} else if bytes.Contains(buf, []byte("\n")) {
		startPostion := 0
		for i := 0; i < len(buf); i++ {
			if buf[i] == '\n' {
				items = append(items, string(buf[startPostion:i]))
				startPostion = i + 1
			}
		}
		return items, COMPLETE
	} else {
		return nil, INCOMPLETE
	}
}

func ParseResponse(response string) error {
	if strings.Split(string(response), "\n")[0] == instruction.RESPONSE_AUTHENTICATED {
		return nil
	}

	if strings.Split(string(response), "\n")[0] == instruction.RESPONSE_OK {
		return nil
	}

	if strings.Split(string(response), "\n")[0] == instruction.RESPONSE_NOT_AUTHENTICATED {
		return errors.New("Client authentication failed")
	}

	if strings.Split(string(response), "\n")[0] != instruction.RESPONSE_OK {
		return errors.New("Operation failed")
	}

	return nil
}
