package msdk

import "fmt"

const StartSignal = "mgc_start:hello_world_lets_hope_noone_ever_sends_this_as_a_log_lol"

// Signal to Magic that tests can start running
func SignalSuccessfulStart() {
	fmt.Println(StartSignal)
}
