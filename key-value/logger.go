package main

import (
	"bufio"
	"fmt"
	"os"
)

type TransactionLogger interface {
	WriteDelete(key string)
	WritePut(key, value string)
	Err() <-chan error
	ReadEvents() (<-chan EventLog, <-chan error)
	Run()
	Close()
}

type EventType byte

const (
	_                     = iota
	EventDelete EventType = iota
	EventPut
)

type EventLog struct {
	Seq   uint64
	Type  EventType
	Key   string
	Value string
}

type FileTransactionLogger struct {
	events       chan<- EventLog // Writer-only channel for sending events
	errors       <-chan error    // Read-only channel for receiving errors
	lastSequence uint64          // The last used event sequence number
	file         *os.File        // The location of the transaction log
}

func (l *FileTransactionLogger) WritePut(key, value string) {
	l.events <- EventLog{Type: EventPut, Key: key, Value: value}
}

func (l *FileTransactionLogger) WriteDelete(key string) {
	l.events <- EventLog{Type: EventDelete, Key: key}
}

func (l *FileTransactionLogger) Err() <-chan error {
	return l.errors
}

func (l *FileTransactionLogger) Run() {
	events := make(chan EventLog, 16)
	l.events = events

	errors := make(chan error, 1)
	l.errors = errors

	go func() {
		defer l.file.Close()

		for e := range events {
			l.lastSequence++

			_, err := fmt.Fprintf(
				l.file,
				"%d\t%d\t%s\t%s\n",
				l.lastSequence, e.Type, e.Key, e.Value)

			if err != nil {
				errors <- err
				return
			}
		}
	}()
}

func (l *FileTransactionLogger) Close() {
	close(l.events)
}

func (l *FileTransactionLogger) ReadEvents() (<-chan EventLog, <-chan error) {
	scanner := bufio.NewScanner(l.file)
	outEvent := make(chan EventLog)
	outError := make(chan error, 1)

	go func() {
		var e EventLog

		defer close(outEvent)
		defer close(outError)

		for scanner.Scan() {
			line := scanner.Text()

			if _, err := fmt.Sscanf(line, "%d\t%d\t%s\t%s",
				&e.Seq, &e.Type, &e.Key, &e.Value); err != nil {

				outError <- fmt.Errorf("input parse error: %w", err)
				return
			}

			if l.lastSequence >= e.Seq {
				outError <- fmt.Errorf("transaction numbers out of sequence")
			}

			l.lastSequence = e.Seq

			outEvent <- e
		}

		if err := scanner.Err(); err != nil {
			outError <- fmt.Errorf("transaction log read failure: %w", err)
			return
		}
	}()

	return outEvent, outError
}

func NewFileTransactionLogger(path string) (TransactionLogger, error) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		return nil, fmt.Errorf("cannot open transaction log file: %w", err)
	}

	return &FileTransactionLogger{file: f}, nil
}
