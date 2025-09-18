package repository

import "fmt"

func (r *Repository) UpdateNotificationStatus(id int, newStatus string) error {
	query := `UPDATE notifications
	SET status = $1
	WHERE id = $2	`

	result, err := r.db.Master.Exec(query, newStatus, id)
	if err != nil {
		return fmt.Errorf("could not update notification status: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not update notification status: %w", err)
	}

	if affected == 0 {
		return ErrNoSuchNotification
	}

	return nil
}
