// Package pulseaudio is a pure-Go (no libpulse) implementation of the PulseAudio native protocol.
//
// Package pulseaudio is a pure-Go (no libpulse) implementation of the PulseAudio native protocol.

// This library is a fork of https://github.com/mafik/pulseaudio
// The original library deliberately tries to hide pulseaudio internals and doesn't expose them.

// For my usecase I needed the exact opposite, access to pulseaudio internals.

package pulseaudio

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/user"
	"path"
	"path/filepath"
)

const version = 32

type packetResponse struct {
	buff *bytes.Buffer
	err  error
}

type packet struct {
	requestBytes []byte
	responseChan chan<- packetResponse
}

type Error struct {
	Cmd  string
	Code uint32
}

func (err *Error) Error() string {
	return fmt.Sprintf("PulseAudio error: %s -> %s", err.Cmd, errors[err.Code])
}

// Client maintains a connection to the PulseAudio server.
type Client struct {
	conn        net.Conn
	clientIndex int
	packets     chan packet
	updates     chan struct{}
	connected   bool
}

// NewClient establishes a connection to the PulseAudio server.
func NewClient(addressArr ...string) (*Client, error) {
	if len(addressArr) < 1 {
		rtp, err := RuntimePath("native")
		if err != nil {
			return nil, err
		}
		addressArr = []string{rtp}
	}

	conn, err := net.Dial("unix", addressArr[0])
	if err != nil {
		return nil, err
	}

	c := &Client{
		conn:      conn,
		packets:   make(chan packet),
		updates:   make(chan struct{}, 1),
		connected: true,
	}

	go c.processPackets()

	err = c.auth()
	if err != nil {
		c.Close()
		return nil, err
	}

	err = c.setName()
	if err != nil {
		c.Close()
		return nil, err
	}

	return c, nil
}

const frameSizeMaxAllow = 1024 * 1024 * 16

func (c *Client) processPackets() {
	recv := make(chan *bytes.Buffer)
	go func(recv chan<- *bytes.Buffer) {
		var err error
		for {
			var b bytes.Buffer
			if _, err = io.CopyN(&b, c.conn, 4); err != nil {
				break
			}
			n := binary.BigEndian.Uint32(b.Bytes())
			if n > frameSizeMaxAllow {
				err = fmt.Errorf("Response size %d is too long (only %d allowed)", n, frameSizeMaxAllow)
				break
			}
			b.Grow(int(n) + 20)
			if _, err = io.CopyN(&b, c.conn, int64(n)+16); err != nil {
				break
			}
			b.Next(20) // skip the header
			recv <- &b
		}
		close(recv)
	}(recv)

	pending := make(map[uint32]packet)
	tag := uint32(0)
	var err error
loop:
	for {
		select {
		case p, ok := <-c.packets: // Outgoing request
			if !ok {
				// Client was closed
				break loop
			}
			// Find an unused tag
			for {
				_, exists := pending[tag]
				if !exists {
					break
				}
				tag++
				if tag == 0xffffffff { // reserved for subscription events
					tag = 0
				}
			}
			if len(p.requestBytes) < 26 {
				p.responseChan <- packetResponse{
					buff: nil,
					err:  fmt.Errorf("request too short. Needs at least 26 bytes"),
				}
				continue
			}
			binary.BigEndian.PutUint32(p.requestBytes, uint32(len(p.requestBytes))-20)
			binary.BigEndian.PutUint32(p.requestBytes[26:], tag) // fix tag
			_, err = c.conn.Write(p.requestBytes)
			if err != nil {
				p.responseChan <- packetResponse{
					buff: nil,
					err:  fmt.Errorf("couldn't send request: %s", err),
				}
			} else {
				pending[tag] = p
			}
		case buff, ok := <-recv: // Incoming request
			if !ok {
				// Client was closed
				break loop
			}
			var tag uint32
			var rsp command
			err = bread(buff, uint32Tag, &rsp, uint32Tag, &tag)
			if err != nil {
				// We've got a weird request from PulseAudio - that should never happen.
				// We could ignore it and continue but it may hide some errors so let's panic.
				panic(err)
			}
			if rsp == commandSubscribeEvent && tag == 0xffffffff {
				select {
				case c.updates <- struct{}{}:
				default:
				}
				continue
			}
			p, ok := pending[tag]
			if !ok {
				// Another case, similar to the one above.
				// We could ignore it and continue but it may hide errors so let's panic.
				panic(fmt.Sprintf("No pending requests for tag %d (%s)", tag, rsp))
			}
			delete(pending, tag)
			if rsp == commandError {
				var code uint32
				bread(buff, uint32Tag, &code)
				cmd := command(binary.BigEndian.Uint32(p.requestBytes[21:]))
				p.responseChan <- packetResponse{
					buff: nil,
					err:  &Error{Cmd: cmd.String(), Code: code},
				}
				continue
			}
			if rsp == commandReply {
				p.responseChan <- packetResponse{
					buff: buff,
					err:  nil,
				}
				continue
			}
			p.responseChan <- packetResponse{
				buff: nil,
				err:  fmt.Errorf("expected Reply or Error but got: %s", rsp),
			}
		}
	}
	// end of packet processing loop, e.g. disconnected
	c.connected = false
	for _, p := range pending {

		p.responseChan <- packetResponse{
			buff: nil,
			err:  fmt.Errorf("PulseAudio client was closed"),
		}
	}
}

func (c *Client) request(cmd command, args ...interface{}) (*bytes.Buffer, error) {
	var b bytes.Buffer
	args = append([]interface{}{uint32(0), // dummy length -- we'll overwrite at the end when we know our final length
		uint32(0xffffffff),   // channel
		uint32(0), uint32(0), // offset high & low
		uint32(0),              // flags
		uint32Tag, uint32(cmd), // command
		uint32Tag, uint32(0), // tag
	}, args...)
	err := bwrite(&b, args...)
	if err != nil {
		return nil, err
	}
	if b.Len() > frameSizeMaxAllow {
		return nil, fmt.Errorf("Request size %d is too long (only %d allowed)", b.Len(), frameSizeMaxAllow)
	}
	responseChan := make(chan packetResponse)

	err = c.addPacket(packet{
		requestBytes: b.Bytes(),
		responseChan: responseChan,
	})
	if err != nil {
		return nil, err
	}

	response := <-responseChan
	return response.buff, response.err
}

func (c *Client) addPacket(data packet) (err error) {
	defer func() {
		if recover() != nil {
			err = fmt.Errorf("connection closed")
		}
	}()
	c.packets <- data
	return nil
}

func (c *Client) auth() error {
	const protocolVersionMask = 0x0000FFFF
	cookiePath, err := cookiePath()
	if err != nil {
		return err
	}
	cookie, err := ioutil.ReadFile(cookiePath)
	if err != nil {
		return err
	}
	const cookieLength = 256
	if len(cookie) != cookieLength {
		return fmt.Errorf("pulse audio client cookie has incorrect length %d: Expected %d (path %#v)",
			len(cookie), cookieLength, cookiePath)
	}
	b, err := c.request(commandAuth,
		uint32Tag, uint32(version),
		arbitraryTag, uint32(len(cookie)), cookie)
	if err != nil {
		return err
	}
	var serverVersion uint32
	err = bread(b, uint32Tag, &serverVersion)
	if err != nil {
		return err
	}
	serverVersion &= protocolVersionMask
	if serverVersion < version {
		return fmt.Errorf("pulseAudio server supports version %d but minimum required is %d", serverVersion, version)
	}
	return nil
}

func (c *Client) setName() error {
	props := map[string]string{
		"application.name":           path.Base(os.Args[0]),
		"application.process.id":     fmt.Sprintf("%d", os.Getpid()),
		"application.process.binary": os.Args[0],
		"application.language":       "en_US.UTF-8",
		"window.x11.display":         os.Getenv("DISPLAY"),
	}
	if current, err := user.Current(); err == nil {
		props["application.process.user"] = current.Username
	}
	if hostname, err := os.Hostname(); err == nil {
		props["application.process.host"] = hostname
	}
	b, err := c.request(commandSetClientName, props)
	if err != nil {
		return err
	}
	var clientIndex uint32
	err = bread(b, uint32Tag, &clientIndex)
	if err != nil {
		return err
	}
	c.clientIndex = int(clientIndex)
	return nil
}

// Close closes the connection to PulseAudio server and makes the Client unusable.
func (c *Client) Close() {
	close(c.packets)
	c.conn.Close()
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// Connected returns a bool specifying if the connection to pulse is alive
func (c *Client) Connected() bool {
	return c != nil && c.connected
}

// RuntimePath resolves a file in the pulse runtime path
// E.g. pass "native" to get the address for pulse' native socket
// Original implementation: https://github.com/pulseaudio/pulseaudio/blob/6c58c69bb6b937c1e758410d3114fc3bc0606fbe/src/pulsecore/core-util.c
// Except we do not support legacy $HOME paths
func RuntimePath(fn string) (string, error) {

	if rtp := os.Getenv("PULSE_RUNTIME_PATH"); rtp != "" {
		return filepath.Join(rtp, fn), nil
	}

	if xdgdir := os.Getenv("XDG_RUNTIME_DIR"); xdgdir != "" {
		if exists(xdgdir) {
			return filepath.Join(xdgdir, "/pulse/", fn), nil
		}
	}

	defaultxdg := fmt.Sprintf("/run/user/%d", os.Getuid())
	if exists(defaultxdg) {
		return filepath.Join(defaultxdg, "/pulse/", fn), nil
	}

	return "", fmt.Errorf("No valid directory for Pulse RuntimePath found")
}

func cookiePath() (string, error) {

	p := filepath.Join(os.Getenv("PULSE_COOKIE"))
	if exists(p) {
		return p, nil
	}

	if confHome := os.Getenv("XDG_CONFIG_HOME"); confHome != "" {
		cookie := filepath.Join(confHome, "/pulse/cookie")
		if exists(cookie) {
			return cookie, nil
		}
	}

	p = filepath.Join(os.Getenv("HOME"), "/.config/pulse/cookie")
	if exists(p) {
		return p, nil
	}

	p = filepath.Join(os.Getenv("HOME"), "/.pulse-cookie")
	if exists(p) {
		return p, nil
	}

	return "", fmt.Errorf("No valid path for Pulse cookie found")
}
