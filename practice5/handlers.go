package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func GetJobsHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	company := r.URL.Query().Get("company")
	limitStr := r.URL.Query().Get("limit")
	afterIDStr := r.URL.Query().Get("after_id")

	limit := 10
	if limitStr != "" {
		v, _ := strconv.Atoi(limitStr)
		if v > 0 {
			limit = v
		}
	}

	afterID := 0
	if afterIDStr != "" {
		v, _ := strconv.Atoi(afterIDStr)
		afterID = v
	}

	args := []interface{}{}
	query := "SELECT id, title, company, salary, created_at FROM jobs"
	where := []string{}

	if company != "" {
		args = append(args, company)
		where = append(where, fmt.Sprintf("company = $%d", len(args)))
	}

	if afterID > 0 {
		var createdAt time.Time
		err := db.QueryRow("SELECT created_at FROM jobs WHERE id = $1", afterID).Scan(&createdAt)
		if err == sql.ErrNoRows {
			http.Error(w, "invalid after_id", 400)
			return
		} else if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		args = append(args, createdAt, afterID)
		where = append(where,
			fmt.Sprintf("(created_at < $%d OR (created_at = $%d AND id < $%d))",
				len(args)-1, len(args)-1, len(args)))
	}

	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}

	query += " ORDER BY created_at DESC, id DESC"
	args = append(args, limit)
	query += fmt.Sprintf(" LIMIT $%d", len(args))

	rows, err := db.Query(query, args...)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var jobs []Job
	for rows.Next() {
		var j Job
		if err := rows.Scan(&j.ID, &j.Title, &j.Company, &j.Salary, &j.CreatedAt); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		jobs = append(jobs, j)
	}

	var nextAfterID *int
	if len(jobs) > 0 {
		id := jobs[len(jobs)-1].ID
		nextAfterID = &id
	}

	resp := JobsResponse{Items: jobs, NextAfterID: nextAfterID}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Query-Time", fmt.Sprintf("%.4fs", time.Since(start).Seconds()))
	json.NewEncoder(w).Encode(resp)
}
