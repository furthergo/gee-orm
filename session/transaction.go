package session

import "github.com/futhergo/gee-orm/log"

func (s *Session)Begin() (err error)  {
	log.Info("Transaction Begin")
	if s.tx, err = s.db.Begin(); err != nil {
		log.Error(err)
		return
	}
	return
}

func (s *Session)Commit() (err error) {
	log.Info("Transaction Commit")
	if err = s.tx.Commit(); err != nil {
		log.Error(err)
		return
	}
	return
}

func (s *Session)Rollback() (err error) {
	log.Info("Transaction Rollback")
	if err = s.tx.Rollback(); err != nil {
		log.Error(err)
		return
	}
	return
}
