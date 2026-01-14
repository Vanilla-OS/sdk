package types

/*	License: GPLv3
	Authors:
		Mirko Brombin <brombin94@gmail.com>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2026
	Description: Vanilla OS SDK component.
*/

// ConnectionInfo represents information about an connection
type ConnectionInfo struct {
	// LocalAddr is the local address of the connection
	LocalAddr string `json:"local_address"`

	// RemoteAddr is the remote address of the connection
	RemoteAddr string `json:"remote_address"`

	// State is the state of the connection
	State ConnState `json:"state"`
}

// ConnState represents the status of a connection
type ConnState string

var (
	ConnStateEnstablished ConnState = "ESTABLISHED"
	ConnStateSynSent      ConnState = "SYN_SENT"
	ConnStateSynRecv      ConnState = "SYN_RECV"
	ConnStateFinWait1     ConnState = "FIN_WAIT1"
	ConnStateFinWait2     ConnState = "FIN_WAIT2"
	ConnStateTimeWait     ConnState = "TIME_WAIT"
	ConnStateClose        ConnState = "CLOSE"
	ConnStateCloseWait    ConnState = "CLOSE_WAIT"
	ConnStateLastAck      ConnState = "LAST_ACK"
	ConnStateListen       ConnState = "LISTEN"
	ConnStateClosing      ConnState = "CLOSING"
	ConnStateNewSynRecv   ConnState = "NEW_SYN_RECV"
	ConnStateUnknown      ConnState = "UNKNOWN"
)
