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

func TestLogout(t *testing.T) {
  conn, err := Dial("localhost:21")
  if err != nil {
    t.Fatal(err)
  }
  err = conn.Login("anonymous", "anonymous@")
  if err != nil {
    t.Error(err)
  }
  err = conn.Logout()
  if err != nil {
    t.Error(err)
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

func TestCheckResponseCode(t *testing.T) {
  // Success
  err := checkResponseCode(2, 230)
  if err != nil {
    t.Error("Should return nil if response code matches! ExpectCode: 2  ActualCode: 230")
  }
  err = checkResponseCode(23, 230)
  if err != nil {
    t.Error("Should return nil if response code matches! ExpectCode: 23  ActualCode: 230")
  }
  err = checkResponseCode(230, 230)
  if err != nil {
    t.Error("Should return nil if response code matches! ExpectCode: 230  ActualCode: 230")
  }

  // Failure
  err = checkResponseCode(2, 500)
  if err == nil {
    t.Error("Should return error if response code doesn't match (ExpectCode: 2  ActualCode: 500)!")
  }
  err = checkResponseCode(23, 500)
  if err == nil {
    t.Error("Should return error if response code doesn't match (ExpectCode: 23  ActualCode: 500)!")
  }
  err = checkResponseCode(230, 500)
  if err == nil {
    t.Error("Should return error if response code doesn't match (ExpectCode: 230  ActualCode: 500)!")
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
