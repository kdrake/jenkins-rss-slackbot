package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Job struct {
	Name  string
	URL   string
	Color string
}

type Assembly struct {
	Jobs []Job
}

// UpdateAssemblyList update actual assembly list
func UpdateAssemblyList(assemblyURL string, jobPrefix []string) {
	r, err := http.Get(assemblyURL)
	if err != nil {
		log.Panic(err)
	}
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Panic(err)
	}

	var assembly Assembly
	if err := json.Unmarshal(body, &assembly); err != nil {
		log.Panic(err)
	}

	var jobIds []interface{}
	for _, job := range assembly.Jobs {
		goodJob := false
		for _, prefix := range jobPrefix {
			if strings.HasPrefix(job.Name, prefix) {
				goodJob = true
				break
			}
		}

		if goodJob {
			row := db.QueryRow("select id from assembly where name = ?", job.Name)
			var id int64
			err = row.Scan(&id)
			if err == sql.ErrNoRows {
				result, err := db.Exec("INSERT INTO assembly(name, url, color) VALUES(?, ?, ?)", job.Name, job.URL, job.Color)
				if err != nil {
					log.Panic(err)
				}

				id, err := result.LastInsertId()
				if err != nil {
					log.Panic(err)
				}
				jobIds = append(jobIds, id)
			} else if err != nil {
				log.Panic(err)
			} else {
				jobIds = append(jobIds, id)
			}
		}
	}

	if jobIds != nil {
		query := fmt.Sprintf("UPDATE assembly SET active=0 WHERE id NOT IN(%s)",
			strings.Join(strings.Split(strings.Repeat("?", len(jobIds)), ""), ","))
		stmt, _ := db.Prepare(query)
		_, err := stmt.Exec(jobIds...)
		if err != nil {
			log.Panic(err)
		}
	}
}

// GetActualAssembly return list of actual assembly
func GetActualAssembly() []string {
	rows, err := db.Query("select url from assembly where active=1")
	if err != nil {
		log.Panic(err)
	}
	defer rows.Close()

	var rssUrls []string
	for rows.Next() {
		var url string
		err = rows.Scan(&url)
		if err != nil {
			log.Panic(err)
		}
		rssUrls = append(rssUrls, url)
	}
	if err = rows.Err(); err != nil {
		log.Panic(err)
	}

	return rssUrls
}
