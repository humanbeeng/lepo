package main

import (
	"bytes"
	"encoding/binary"
	"io"
)

type Command byte

const (
	CmdReg Command = iota
	CmdDereg
	CmdPing
	CmdStrat
)

type Strategy byte

const (
	RoundRobin Strategy = iota
	HashedURL
)

type CommandRegister struct {
	Addr []byte
}

type CommandDeRegister struct {
	Addr []byte
}

type Status byte

const (
	Ok Status = iota
	Error
)

func (c *CommandRegister) Bytes() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, CmdReg)

	addrLen := int32(len(c.Addr))
	binary.Write(buf, binary.LittleEndian, addrLen)
	binary.Write(buf, binary.LittleEndian, c.Addr)
	return buf.Bytes()
}

func (c *CommandDeRegister) Bytes() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, CmdReg)

	addrLen := int32(len(c.Addr))
	binary.Write(buf, binary.LittleEndian, addrLen)
	binary.Write(buf, binary.LittleEndian, c.Addr)
	return buf.Bytes()
}

func ParseAdminCommand(r io.Reader) (Command, error) {
	var cmd Command
	err := binary.Read(r, binary.LittleEndian, &cmd)
	if err != nil {
		return CmdReg, err
	}
	return cmd, nil
}

func ParseRegisterCommand(r io.Reader) CommandRegister {
	var reg CommandRegister
	var addrLen int32
	binary.Read(r, binary.LittleEndian, &addrLen)
	reg.Addr = make([]byte, addrLen)
	binary.Read(r, binary.LittleEndian, &reg.Addr)
	return reg
}

func ParseDeRegisterCommand(r io.Reader) CommandDeRegister {
	var reg CommandDeRegister
	var addrLen int32
	binary.Read(r, binary.LittleEndian, &addrLen)
	reg.Addr = make([]byte, addrLen)
	binary.Read(r, binary.LittleEndian, &reg.Addr)
	return reg
}
