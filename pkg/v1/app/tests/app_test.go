package tests

import (
	"fmt"
	"testing"

	"github.com/vanilla-os/sdk/pkg/v1/app"
	"github.com/vanilla-os/sdk/pkg/v1/app/types"
)

func TestNewApp(t *testing.T) {
	app, err := app.NewApp(types.AppOptions{
		RDNN:    "com.vanilla-os.batsignal",
		Name:    "BatSignal",
		Version: "1.0.0",
	})
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}
	if app.Sign == "" {
		t.Errorf("Expected app.Sign to be non-empty, got empty")
	}
	fmt.Println("App created")
	fmt.Printf("\tSign: %s\n", app.Sign)
	fmt.Printf("\tRDNN: %s\n", app.RDNN)
	fmt.Printf("\tName: %s\n", app.Name)
	fmt.Printf("\tVersion: %s\n", app.Version)
}
