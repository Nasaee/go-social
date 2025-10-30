package db

import (
	"context"
	"database/sql"
	"log"
	"time"
)

// New สร้าง connection ไปยังฐานข้อมูล Postgres และตั้งค่าพารามิเตอร์ของ connection pool
func New(addr string, maxOpenConns, maxIdleConns int, maxIdleTime string) (*sql.DB, error) {
	db, err := sql.Open("postgres", addr)
	if err != nil {
		return nil, err
	}

	// จำกัดจำนวน connection สูงสุดที่เปิดพร้อมกันได้
	db.SetMaxOpenConns(maxOpenConns)

	// จำกัดจำนวน connection ที่สามารถว่าง (idle) อยู่ใน pool ได้พร้อมกัน
	db.SetMaxIdleConns(maxIdleConns)

	// แปลงค่าระยะเวลาจาก string → time.Duration เช่น "15m", "1h"
	duration, err := time.ParseDuration(maxIdleTime)
	if err != nil {
		// ถ้าแปลงไม่สำเร็จ เช่น format ผิด จะคืน error กลับไป
		return nil, err
	}

	// ตั้งค่าระยะเวลาสูงสุดที่ connection ว่างจะถูกเก็บไว้ก่อนปิด
	db.SetConnMaxIdleTime(duration)

	// สร้าง context พร้อม timeout 5 วินาที เพื่อทดสอบการเชื่อมต่อ
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // คืน resource หลังจากใช้งานเสร็จ

	// PingContext คือการ "ทดสอบ" ว่าฐานข้อมูลตอบสนองหรือไม่
	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}

	log.Println("database connection pool established 🌍")
	return db, nil
}
