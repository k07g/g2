package main

import (
	"net/http"
	"os"
	"testing"
	"time"
)

func TestRun_Shutdown(t *testing.T) {
	srv := newServer("0") // ポート 0 で空きポートを自動割り当て
	quit := make(chan os.Signal, 1)

	done := make(chan error, 1)
	go func() {
		done <- run(srv, quit)
	}()

	// サーバー起動を待つ
	time.Sleep(50 * time.Millisecond)

	// シグナルを送信してシャットダウンをトリガー
	quit <- os.Interrupt

	select {
	case err := <-done:
		if err != nil {
			t.Errorf("run returned error: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Error("shutdown timed out")
	}
}

func TestRun_RequestHandled(t *testing.T) {
	srv := newServer("18080")
	quit := make(chan os.Signal, 1)

	go run(srv, quit)

	// サーバー起動を待つ
	time.Sleep(50 * time.Millisecond)

	resp, err := http.Get("http://localhost:18080/tasks")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	quit <- os.Interrupt
}
