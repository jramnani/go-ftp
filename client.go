package ftp

import (
  "bufio"
  "fmt"
  "io"
  "net"
  "os"
  "regexp"
  "strconv"
  "strings"
)

// Knows the control connection where commands are sent to the server.
type Connection struct {
  control io.ReadWriteCloser
  hostname string
}

var CRLF = "\r\n"
var ASCII = "A"
var BINARY = "I"
var IMAGE = "I" //Synonymous with "Binary"

// Dials up a remote FTP server.
// host should be in the form of address:port e.g. myserver:21 or myserver:ftp
// Returns a pointer to a Connection
func Dial(host string) (*Connection, os.Error) {
  if host == "" {
    return nil, os.NewError("FTP Connection Error: Host can not be blank!")
  }
  if !hasPort(host) {
    return nil, os.NewError("FTP Connection Error: Host must have a port! e.g. host:21")
  }
  conn, err := net.Dial("tcp", "", host)
  if err != nil {
    return nil, err
  }
  // Upon connect, most servers respond with a welcome message.
  // The welcome message contains a status code, just like any other command.
  // TODO: Handle servers with no welcome message.
  welcomeMsg := make([]byte, 1024)
  _, err = conn.Read(welcomeMsg)
  if err != nil {
    return nil, os.NewError("Couldn't read the server's initital connection information. Error: " + err.String())
  }
  code, err := strconv.Atoui(string(welcomeMsg[0:3]))
  err = checkResponseCode(2, code)
  if err != nil {
    return nil, os.NewError("Couldn't read the server's Welcome Message. Error: " + err.String())
  }
  // This doesn't work for IPv6 addresses.
  hostParts := strings.Split(host, ":", 2)
  return &Connection{conn, hostParts[0]}, nil
}

// Executes an FTP command.
// Sends the command to the server.
// Returns the response code, the response text from the server, and any errors.
// The response code will be zero if an error is encountered.
// The response string will be the empty string if an error is encountered.
func (c *Connection) Cmd(command string, arg string) (code uint, response string, err os.Error) {
  // Format command to be sent to the server.
  formattedCommand := command + " " + arg + CRLF

  // Send command to the server.
  _, err = c.control.Write([]byte(formattedCommand))
  if err != nil {
    return 0, "", err
  }

  // Process the response.
  reader := bufio.NewReader(c.control)
  regex := regexp.MustCompile("[0-9][0-9][0-9] ")
  for {
    ln, err := reader.ReadString('\n')
    if err != nil {
      return 0, "", err
    }

    response += ln
    if regex.MatchString(ln) {
      break
    }
  }
  code, err = strconv.Atoui(response[0:3])
  if err != nil {
    return 0, response, err
  }
  return code, response, err
}

// Log into a FTP server using username and password.
func (c *Connection) Login(user string, password string) os.Error {
  if user == "" {
    return os.NewError("FTP Connection Error: User can not be blank!")
  }
  if password == "" {
    return os.NewError("FTP Connection Error: Password can not be blank!")
  }
  // TODO: Check the server's response codes.
  _, _, err := c.Cmd("USER", user)
  _, _, err = c.Cmd("PASS", password)
  if err != nil {
    return err
  }
  return nil
}

// Download a file from a remote server.
// Assumes only passive FTP connections for now.
func (c *Connection) DownloadFile(src, dest, mode string) os.Error {
  // Use PASV to set up the data port.
  pasvCode, pasvLine, err := c.Cmd("PASV", "")
  if err != nil {
    return err
  }
  pasvErr := checkResponseCode(2, pasvCode)
  if pasvErr != nil {
    msg := fmt.Sprintf("Cannot set PASV. Error: %v", pasvErr)
    return os.NewError(msg)
  }
  dataPort, err := extractDataPort(pasvLine)
  /*_, err = extractDataPort(pasvLine)*/
  if err != nil {
    return err
  }

  // Set the TYPE (ASCII or Binary)
  typeCode, typeLine, err := c.Cmd("TYPE", mode)
  if err != nil {
    return err
  }
  typeErr := checkResponseCode(2, typeCode)
  if typeErr != nil {
    msg := fmt.Sprintf("Cannot set TYPE. Error: '%v'. Line: '%v'", typeErr, typeLine)
    return os.NewError(msg)
  }

  // Can't use Cmd() for RETR because it doesn't return until *after* you've
  // downloaded the requested file.
  command := []byte("RETR " + src + CRLF)
  _, err = c.control.Write(command)
  if err != nil {
    return err
  }

  // Open connection to remote data port.
  remoteConnectString := c.hostname + ":" + fmt.Sprintf("%d", dataPort)
  download_conn, err := net.Dial("tcp", "", remoteConnectString)
  defer download_conn.Close()
  if err != nil {
    msg := fmt.Sprintf("Couldn't connect to server's remote data port. Error: %v", err)
    return os.NewError(msg)
  }

  // Set up the destination file
  var filePerms uint32 = 0664
  destFile, err := os.Open(dest, os.O_CREAT|os.O_WRONLY, filePerms)
  defer destFile.Close()
  if err != nil {
    msg := fmt.Sprintf("Cannot open destination file, '%s'. %v", dest, err)
    return os.NewError(msg)
  }

  // Buffer for downloading and writing to file
  bufLen := 1024
  buf := make([]byte, bufLen)

  // Read from the server and write the contents to a file
  for {
    bytesRead, readErr := download_conn.Read(buf)
    if bytesRead > 0 {
      _, err := destFile.Write(buf[0:bytesRead])
      if err != nil {
        msg := fmt.Sprintf("Coudn't write to file, '%s'. Error: %v", dest, err)
        return os.NewError(msg)
      }
    }
    if readErr == os.EOF {
      break
    }
    if readErr != nil {
      return readErr
    }
  }
  return nil
}

// Given an prefix, does the response code match the expected code?
func checkResponseCode(expectCode, code uint) os.Error {
  if 1 <= expectCode && expectCode < 10 && code/100 != expectCode ||
    10 <= expectCode && expectCode < 100 && code/10 != expectCode ||
      100 <= expectCode && expectCode < 1000 && code != expectCode {
        msg := fmt.Sprintf("Bad response from server. Expected: %d, Got: %d", expectCode, code)
        return os.NewError(msg)
      }
  return nil
}

// Interrogate a server response for the remote port on which to connect.
// Returns the port to be used as the data port for transfers.
func extractDataPort(line string) (port uint, err os.Error) {
  // We only care about the last two octets
  portPattern:= "[0-9]+,[0-9]+,[0-9]+,[0-9]+,([0-9]+,[0-9]+)"
  re, err := regexp.Compile(portPattern)
  if err != nil {
    return 0, err
  }
  match := re.FindStringSubmatch(line)
  if len(match) == 0 {
    msg := "Cannot find data port in server output: " + line
    return 0, os.NewError(msg)
  }
  octets := strings.Split(match[1], ",", 2)
  firstOctet, _ := strconv.Atoui(octets[0])
  secondOctet, _ := strconv.Atoui(octets[1])
  port = firstOctet * 256 + secondOctet

  return port, nil
}

// Reused from src/pkg/http/client.go
// Given a string of the form "host", "host:port", or "[ipv6::address]:port",
// return true if the string includes a port.
func hasPort(s string) bool { return strings.LastIndex(s, ":") > strings.LastIndex(s, "]") }

