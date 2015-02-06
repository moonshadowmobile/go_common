package util

import (
	"sync"
	"flag"
	"os"
	"os/signal"
	"log"
	"strings"
	"io/ioutil"
	"fmt"
	"runtime"
	"syscall"
	"crypto/rand"
)

func UUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatalln("uuid error: ", err.Error())
		return ""
	}
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
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
otherwise use os.Hostname(). */
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
		"Using '%s' (the contents of /home/cast/.hostname) to resolve config file.\n",
		s)

	return s, nil
}

/* We cache this so we don't have to make repeated syscalls elsewhere in the code. */
var MAX_PROCS int = 0
func GetMaxProcs() int {
	if MAX_PROCS == 0 {
		MAX_PROCS = runtime.NumCPU()
		runtime.GOMAXPROCS(MAX_PROCS)
	}

	return MAX_PROCS
}

// Startup message utils
var startup_msgs []string
func RegStartup(msgs ...string) {
	for _, m := range msgs {
		startup_msgs = append(startup_msgs, m)
	}
}

func PrintStartup() {
	for _, msg := range startup_msgs {
		fmt.Printf("STARTUP: %s\n", msg)
	}
}

// Flag utils
var FlagConfigFilename *string
var flagVersion *bool
func RegBaseFlags() {
	hn, err := GetHostname()
	if err != nil {
		log.Fatalf(err.Error())
	}
	FlagConfigFilename = flag.String("config", hn,
		"Accepts a filename (no extension). For a file named 'foo.json', use foo")
	flagVersion = flag.Bool("v", false, "Prints the current binary version")
}
func HandleExitableFlags(appVersion *string) {
	flag.Parse()
	if *flagVersion {
		fmt.Println(*appVersion)
		os.Exit(0)
	}
}

// Signal utils
type SignalHandler func() string
var handlers map[string][]SignalHandler
func ListenForSignals() {
	handlers = make(map[string][]SignalHandler)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGUSR1)
	handlerLock := new(sync.Mutex)
	go listenForSIGUSR1(sigChan, handlerLock)
}
func AddSignalHandler(sigType string, h SignalHandler) {
	handlers[sigType] = append(handlers[sigType], h)
}
func listenForSIGUSR1(sigChan chan os.Signal, handlerLock *sync.Mutex) {
	<-sigChan
	memstats := &runtime.MemStats{}
	runtime.ReadMemStats(memstats)

	fmt.Printf("\n*** Current process state ***\n")
	fmt.Printf("Goroutine count: %d\n\n", runtime.NumGoroutine())

	fmt.Println("*** General statistics ***")
	fmt.Printf("System reported bytes: %d\n", memstats.Sys)
	fmt.Printf("Number of mallocs: %d\n", memstats.Mallocs)
	fmt.Printf("Number of frees: %d\n", memstats.Frees)
	fmt.Printf("Allocated heap bytes: %d\n", memstats.HeapAlloc)
	fmt.Printf("Bytes in use: %d\n", memstats.StackInuse)
	fmt.Printf("Bytes in use (system reported): %d\n\n", memstats.StackSys)

	fmt.Printf("*** Garbage collector statistics ***\n")
	fmt.Printf("Next garbage collection (bytes): %d\n", memstats.NextGC)
	fmt.Printf("Last garbage collection time (ns): %d\n", memstats.LastGC)
	fmt.Printf("Number of garbage collections: %d\n", memstats.NumGC)
	fmt.Printf("Time spent in garbage collection (ns): %d\n\n",
		memstats.PauseTotalNs)

	handlerLock.Lock()
	sh := handlers["SIGUSR1"]
	for i := 0; i < len(sh); i++ {
		fmt.Printf("%s\n", sh[i]())
	}
	handlerLock.Unlock()

	newSigChan := make(chan os.Signal, 1)
	signal.Notify(newSigChan, syscall.SIGUSR1)

	go listenForSIGUSR1(newSigChan, handlerLock)
}
