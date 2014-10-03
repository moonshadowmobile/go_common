package go_common

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"
)

func SetMaxProcs() int {
	procs := runtime.NumCPU()
	runtime.GOMAXPROCS(procs)

	return procs
}

func PrintStartupMsg(max_procs int, description string, version string,
	config_filename string) {

	fmt.Printf("STARTUP: %s\n", description)
	fmt.Printf("STARTUP: Version: %s\n", version)
	fmt.Printf("STARTUP: PID: %d\n", os.Getpid())
	fmt.Printf("STARTUP: Using a maximum of %d logical CPUs. \n", max_procs)
	fmt.Printf("STARTUP: Configuration file: %s\n", config_filename)
}

/* Extracts the middle directory structure for the run executable. This can be
concatenated with os.Getwd() to get the absolute path of the project root. */
func GetBasePathForExecutable(exec string) string {
	wd, wd_err := os.Getwd()
	if wd_err != nil {
		log.Fatalf("ERROR: Unable to get current working directory.")
	}

	// Remove leading . in the ./ used to run the executable
	a := strings.TrimPrefix(os.Args[0], ".")
	// Remove the executable name
	a = strings.TrimSuffix(a, exec)

	// Jump out of bin into project root
	a = wd + a + "../"

	return a
}

/* Use the contents of /home/cast/.hostname file if it exists,
otherwise use os.Hostname(). This is to circumvent Joyent uuid hostnames. */
func GetHostname() (string, error) {
	hn_bytes, err := ioutil.ReadFile("/home/cast/.hostname")
	if err != nil {
		hn, err := os.Hostname()
		if err != nil {
			return "", err
		}

		return hn, nil
	}

	// Split off the newline that can confuse the parser
	lines := strings.Split(string(hn_bytes), "\n")
	s := lines[0]

	fmt.Printf(
		"STARTUP: Using '%s' (the contents of /home/cast/.hostname) to resolve config file.\n",
		s)

	return s, nil
}
