package ftp

import (
  /*"fmt"*/
  "testing"
)

// TODO: Mock out network access.  Tests currently require an FTP server
// running on localhost.

func TestDial(t *testing.T) {
  _, err := Dial("")
  if err == nil {
    t.Error("Dial() should return an error on blank 'host'!")
  }
  _, err = Dial("localhost:21")
  if err != nil {
    t.Error(err)
  }
}

func TestLogin(t *testing.T) {
  conn, err := Dial("localhost:21")
  if err != nil {
    t.Fatal(err)
  }
  err = conn.Login("anonymous", "anonymous@")
  if err != nil {
    t.Error(err)
  }
  err = conn.Login("", "anonymous@")
  if err == nil {
    t.Error("Login() should return an error on blank 'user'!")
  }
  err = conn.Login("anonymous", "")
  if err == nil {
    t.Error("Login() should return an error on blank 'password'!")
  }
}

func TestCmd(t *testing.T) {
  conn, err := Dial("localhost:ftp")
  if err != nil {
    t.Fatal(err)
  }
  _, _, err = conn.Cmd("USER", "anonymous")
  if err != nil {
    t.Error(err)
  }
}

func TestExtractDataPort(t *testing.T) {
  test_string := "227 Entering Passive Mode (127,0,0,1,205,238)."
  port, err := extractDataPort(test_string)
  if err != nil {
    t.Error(err)
  }
  if port != 52718 {
    t.Error("Failed port calculation! Expected 52718, got", port)
  }
}

func TestDownload(t *testing.T) {
  conn, err := Dial("localhost:21")
  if err != nil {
    t.Fatal(err)
  }
  err = conn.Login("anonymous", "anonymous@")
  if err != nil {
    t.Error(err)
  }
  err = conn.DownloadFile("/test_file.txt", "./download_test.txt", BINARY)
  if err != nil {
    t.Error(err)
  }
}

func TestUpload(t *testing.T) {
  conn, err := Dial("localhost:21")
  if err != nil {
    t.Fatal(err)
  }
  err = conn.Login("anonymous", "anonymous@")
  if err != nil {
    t.Error(err)
  }
  err = conn.UploadFile("upload_file.txt", "/up/uploaded.txt", BINARY)
  if err != nil {
    t.Error(err)
  }
}
