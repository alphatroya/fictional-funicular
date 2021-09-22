package main

import "errors"

type Storage struct {
	token string
	host  string
}

func (s *Storage) SetToken(token string, chat int64) {
	s.token = token
}

func (s *Storage) GetToken(_ int64) (string, error) {
	return s.token, nil
}

func (s *Storage) SetHost(host string, chat int64) {
	s.host = host
}

func (s *Storage) GetHost(chat int64) (string, error) {
	return s.host, nil
}

func (s *Storage) SetActivity(activity string, chat int64) {
}

func (s *Storage) GetActivity(chat int64) (string, error) {
	return "", errors.New("Not implemented activity")
}

func (s *Storage) ResetData(chat int64) error {
	return nil
}
