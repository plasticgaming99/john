package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

var (
	downloader = "curl -O %d -o %d"
)

func main() {
	if os, _ := DetectOS(); os == "other" {
		fmt.Println("Maybe unsupported os?")
		panic(1)
	}
	if len(os.Args) == 1 {
		fmt.Println("no command was input.")
	} else {
		switch os.Args[1] {
		case "install":
			if len(os.Args) <= 1 {
				fmt.Println("Installing latest stable")
				installnvim("stable")
			} else {
				if os.Args[2] == "nightly" {
					fmt.Println("Installing latest nightly")
					installnvim("nightly")
				} else {
					fmt.Printf("Installing version %s", os.Args[2])
					installnvim(os.Args[2])
				}
			}
		case "update":
			fmt.Println("i'll update nvim.")
		case "release":
			if len(os.Args) <= 1 {
				fmt.Println("release needs 1 more options.")
			}
		case "help":
			help()
		default:
			fmt.Printf("invaild command: %s \n", os.Args[1])
		}
	}
}

func help() {
	fmt.Println("john usage:")
	fmt.Println("    john [subcommand] [options]")
	fmt.Println("  john has 4 subcommands:")
	fmt.Println("    install : install nvim ")
	fmt.Println("    update  : update nvim to latest")
	fmt.Println("    release : change stable or nightly")
	fmt.Println("    help    : show help of john")
	fmt.Println("  john has no options")
}

func installnvim(nvimver string) {
	Check4dir("downloads")
	downloaddir := Getjohndir() + "downloads/temp"
	//exec.Command(SplitCmd(fmt.Sprintf(downloader, "https://github.com/neovim/neovim/releases/download/"+nvimver+Gettargetfile(), downloaddir)))
	fmt.Println("curl", "-#L", "https://github.com/neovim/neovim/releases/download/"+nvimver+Gettargetfile(), "-o", downloaddir+Gettargetfile())
	executecmd("curl", "-#L", "https://github.com/neovim/neovim/releases/download/"+nvimver+Gettargetfile(), "-o", downloaddir+Gettargetfile())

	// install nvim
	Check4dir("versions")
	installto := string("")
	if nvimver == "nightly" || nvimver == "stable" {
		Check4dir("rolling")
		installto = "versions/rolling"
	} else {
		Check4dir("tags")
		installto = "versions/tags"
	}
	fmt.Println(installto)

	Check4dir("downloads/temp")

	if os, _ := DetectOS(); os == "windows" {
		executecmd("tar", "-xvf", Getjohndir()+"/downloads/temp"+Gettargetfile(), "-C", Getjohndir()+"/rolling")
	}
}

func Getjohndir() string {
	var instdir = string("")
	homedir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	osname, _ := DetectOS()

	if osname == "windows" {
		instdir = homedir + "/AppData/Local/john/"
		_, err := os.Stat(instdir)
		if err != nil {
			os.Mkdir(instdir, os.ModePerm)
		}
		return instdir
	} else if osname == "linux" {
		instdir = homedir + "/.local/share/john/"
		_, err := os.Stat(instdir)
		if err != nil {
			os.Mkdir(instdir, os.ModePerm)
		}
		return instdir
	} else {
		return "n/a"
	}
}

func Gettargetfile() string {
	os, arch := DetectOS()
	var targfile = "/"
	if os == "windows" {
		targfile += "nvim-win" + strconv.Itoa(arch) + ".zip"
		return targfile
	} else if os == "linux" {
		if arch == 32 {
			fmt.Println("Unsupported arch")
			targfile += "notavaliable"
		} else if arch == 64 {
			targfile += "nvim-linux64.tar.gz"
		}
	} else {
		fmt.Println("Maybe unsupported os?")
		targfile += "unsupported"
	}
	return targfile
}

func DetectOS() (string, int) {
	// Detect OS
	var os = ""
	if runtime.GOOS == "windows" {
		os = "windows"
	} else if runtime.GOOS == "linux" {
		os = "linux"
	} else {
		os = "other"
	}
	// Detect Architecture
	var arch = int(0)
	if runtime.GOARCH == "386" {
		arch = 32
	} else if runtime.GOARCH == "amd64" {
		arch = 64
	} else {
		arch = 0
	}
	// Return simplified value
	return os, arch
}

func Check4dir(relativepathfromjohndir string) {
	_, err := os.Stat(Getjohndir() + relativepathfromjohndir)
	if err != nil {
		os.Mkdir(Getjohndir()+relativepathfromjohndir, os.ModePerm)
	}
	return
}

func SplitCmd(inputcmd string) (string, string) {
	slicedcmd := strings.Split(inputcmd, " ")
	firstcmd := slicedcmd[0]
	argument := strings.Join(slicedcmd[1:], " ")
	return firstcmd, argument
}

func executecmd(cmdname string, args ...string) {
	toexec := exec.Command(cmdname, args...)
	toexec.Stdin = os.Stdin
	toexec.Stdout = os.Stdout
	toexec.Stderr = os.Stderr
	toexec.Start()
	toexec.Wait()
}
