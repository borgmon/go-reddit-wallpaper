package main

import "github.com/robfig/cron/v3"

func clearAndSetCron(text string) (cron.EntryID, error) {
	clearAllCronJobs()
	return cronJob.AddFunc(text, func() {
		go Start()
	})
}

func clearAllCronJobs() {
	jobs := cronJob.Entries()
	for _, job := range jobs {
		cronJob.Remove(job.ID)
	}
}

func newCron() *cron.Cron {
	return cron.New()
}

func parseCron(text string) (cron.Schedule, error) {
	return cron.ParseStandard(text)
}
