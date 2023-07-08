package imaiot

import "testing"

func TestMQTTBrokerConnection(t *testing.T) {
	err := ConnectIOT()
	if err != nil {
		t.Error("connect is failed:", err)
	}
	t.Log("success connect mqtt broker")
}
