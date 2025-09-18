package service

func (s *Service) UpdateNotificationStatus(id int, newStatus string) error {
	if err := s.cache.Set(id, newStatus); err != nil {
		return err
	}

	if err := s.storage.UpdateNotificationStatus(id, newStatus); err != nil {
		return err
	}

	return nil
}
