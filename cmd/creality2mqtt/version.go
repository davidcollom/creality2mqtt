package main

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func init() {
	rootCmd.Version = version + " (commit: " + commit + ", date: " + date + ")"
}
