package main

import "errors"

type storage struct {
	token string
	host  string
}

func (s *storage) SetToken(token string, chat int64) {
	s.token = token
}

func (s *storage) GetToken(_ int64) (string, error) {
	return s.token, nil
}

func (s *storage) SetHost(host string, chat int64) {
	s.host = host
}

func (s *storage) GetHost(chat int64) (string, error) {
	return s.host, nil
}

func (s *storage) SetActivity(activity string, chat int64) {
}

func (s *storage) GetActivity(chat int64) (string, error) {
	return "", errors.New("Activity is not implemented")
}

func (s *storage) ResetData(chat int64) error {
	return nil
}
