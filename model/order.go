package model

import (
	"database/sql"
	"time"

	"github.com/pkg/errors"
)

// Order is the main struct (now the only) to request for
type Order struct {
	ID          int64  `json:"id,omitempty" db:"id,omitempty"`
	Name        string `json:"name" db:"name,omitempty"`
	StuID       string `json:"stu_id" db:"stu_id,omitempty"`
	Date        string `json:"date,omitempty" db:"date,omitempty"`
	Comment     string `json:"comment" db:"comment"`
	ServiceType string `json:"service_type" db:"service_type"`
	// Below are info to modify after order created
	Rating     int64     `json:"rating,omitempty" db:"rating"`
	OperaterID int64     `json:"operater_id,omitempty" db:"operater_id"`
	SecretID   string    `json:"secret_id,omitempty" db:"secret_id"`
	DoneFlag   bool      `json:"done_flag" db:"done_flag"`
	CreateTime time.Time `json:"create_time" db:"create_time"`
	UpdateTime time.Time `json:"update_time" db:"update_time"`
}

type OrderGuarded struct {
	ID          int64  `json:"id,omitempty" db:"id,omitempty"`
	Name        string `json:"name" db:"name,omitempty"`
	Date        string `json:"date,omitempty" db:"date,omitempty"`
	ServiceType string `json:"service_type" db:"service_type"`
	// Below are info to modify after order created
	Rating     int64     `json:"rating,omitempty" db:"rating"`
	DoneFlag   bool      `json:"done_flag" db:"done_flag"`
	CreateTime time.Time `json:"create_time" db:"create_time"`
}

func (o *Order) Insert(db Storager) (err error) {
	sqlStr := "INSERT INTO orders (name, stu_id, date, comment, service_type, secret_id, create_time, update_time)VALUES(?,?,?,?,?,?,?,?)"
	_, err = db.Queryx(sqlStr, o.Name, o.StuID, o.Date, o.Comment, o.ServiceType, o.SecretID, time.Now().Local(), time.Now().Local())
	if err != nil {
		err = errors.Wrap(err, "insert error")
		return
	}
	return
}

func OrderByStuID(db Storager, stuID string) (o Order, err error) {
	sqlStr := "SELECT * FROM orders where stu_id = ? AND done_flag = FALSE ORDER BY CREATE_TIME DESC LIMIT 1"
	err = db.QueryRowx(sqlStr, stuID).StructScan(&o)
	if err != nil && err != sql.ErrNoRows {
		err = errors.Wrap(err, "order by stu_id error")
		return
	}
	if err == sql.ErrNoRows {
		err = nil
		o.ID = -1
	}
	return
}

func OrderByID(db Storager, ID int) (err error) {
	return
}

func UpdateOrderDoneFlagBySecret(db Storager, secretID string) (err error) {
	sqlStr := "UPDATE orders SET done_flag = TRUE, update_time = ? WHERE secret_id = ?"
	_, err = db.Exec(sqlStr, time.Now().Local(), secretID)

	if err != nil && err != sql.ErrNoRows {
		err = errors.Wrap(err, "order by stu_id error")
		return
	}
	if err == sql.ErrNoRows {
		err = nil
	}
	return
}

func OrderBySecret(db Storager, secretID string) (o Order, err error) {
	sqlStr := "SELECT * FROM orders where secret_id = ?"
	err = db.QueryRowx(sqlStr, secretID).StructScan(&o)
	if err != nil && err != sql.ErrNoRows {
		err = errors.Wrap(err, "order by stu_id error")
		return
	}
	if err == sql.ErrNoRows {
		err = nil
		o.ID = -1
	}
	return
}

func OrderPager(db Storager, pg int, perpg int) (ol []Order, err error) {
	sqlStr := "SELECT * FROM orders ORDER BY id desc LIMIT ? OFFSET ?"
	rows, er := db.Queryx(sqlStr, perpg, perpg*(pg-1))
	if er != nil {
		err = errors.Wrap(er, "order by page error")
		return
	}
	defer rows.Close()
	o := Order{}
	for rows.Next() {
		err = rows.StructScan(&o)
		if err != nil {
			err = errors.Wrap(err, "order by page, fetching error")
			return
		}
		ol = append(ol, o)
	}
	return
}

func OrderCount(db Storager) (cnt int, err error) {
	sqlStr := "SELECT COUNT(*) FROM orders"
	err = db.QueryRowx(sqlStr).Scan(&cnt)
	if err != nil {
		err = errors.Wrap(err, "get order count error")
		return
	}
	return
}

func OrderGuardedPager(db Storager, pg int, perpg int) (ol []OrderGuarded, err error) {
	sqlStr := "SELECT * FROM orders ORDER BY id desc LIMIT ? OFFSET ?"
	rows, er := db.Unsafe().Queryx(sqlStr, perpg, perpg*(pg-1))
	if er != nil {
		err = errors.Wrap(er, "order by page error")
		return
	}
	defer rows.Close()
	o := OrderGuarded{}
	for rows.Next() {
		err = rows.StructScan(&o)
		if err != nil {
			err = errors.Wrap(err, "order by page, fetching error")
			return
		}
		ol = append(ol, o)
	}
	return
}
