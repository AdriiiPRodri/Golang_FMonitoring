package main

import (
	"crypto/md5"
	"io"
	"os"
	"log"
	"encoding/hex"
	"fmt"
	"time"
	"syscall"
	"path/filepath"
)

const (
	not_symlink = "NULL"
	deleted_symlink = "DELETED"
)

/*
	Error types:
		1X: Input errors
			10: Invalid number of arguments
		2X: File operation errors
			20: Can not find the specified file
		3X: Checksum errors
			30: Can not open the specified file
			31: Error copying the file in the hash interface
*/

func fatal_error(err error, code int) {
	if err != nil {
		log.Fatalf("[%d] ", code, err.Error())
	}
}

func _get_stats(file_path string) (os.FileInfo) {
	stats, err := os.Lstat(file_path); fatal_error(err, 20)

	return stats
}

func get_filename(file_path string) (string) {
	stats := _get_stats(file_path)

	return stats.Name()
}

func get_filepath(file_path string) (string) {
	filepath, err := filepath.Abs(file_path); fatal_error(err, 20)

	return filepath
}

func get_file_size(file_path string) (int64) {
	stats := _get_stats(file_path)

	return stats.Size()
}

func get_permissions(file_path string) (os.FileMode) {
	stats := _get_stats(file_path)

	return stats.Mode()
}

func get_file_modtime(file_path string) (time.Time) {
	stats := _get_stats(file_path)

	return stats.ModTime()
}

func file_is_dir(file_path string) (bool) {
	stats := _get_stats(file_path)

	return stats.IsDir()
}

func file_is_symbolic(file_path string) (bool) {
	stats := _get_stats(file_path)

	return stats.Mode() & os.ModeSymlink != 0
}

func get_file_inode(file_path string) (uint64) {
	stats := _get_stats(file_path)

	return stats.Sys().(*syscall.Stat_t).Ino
}

func get_file_ino(file_path string) (uint64) {
	stats := _get_stats(file_path)

	return stats.Sys().(*syscall.Stat_t).Ino
}

func get_file_hardlinks(file_path string) (uint64) {
	stats := _get_stats(file_path)

	return stats.Sys().(*syscall.Stat_t).Nlink
}

func get_stats(file_path string) (string, int64, os.FileMode, time.Time, bool, uint64, uint64, bool, string) {
	stats := _get_stats(file_path)
	is_symbolic := stats.Mode() & os.ModeSymlink != 0
	var resolve_symbolic string = not_symlink
	var err error
	if is_symbolic {
		resolve_symbolic, err = os.Readlink(file_path)
		if err != nil {
			log.Print("[WARNING] Unable to resolve the symlink", file_path)
		}
		if _, err := os.Stat(file_path); os.IsNotExist(err) {
			log.Print("[WARNING] Symbolic link %s goes to non existent file %s", file_path, resolve_symbolic)
			resolve_symbolic = deleted_symlink
		}
	}

	return get_filepath(file_path), stats.Size(), stats.Mode(), stats.ModTime(), stats.IsDir(), stats.Sys().(*syscall.Stat_t).Ino, stats.Sys().(*syscall.Stat_t).Nlink, is_symbolic, resolve_symbolic
}

func get_checksum(file_path string) (string) {
	// Initialize the file md5 string
	var file_md5 string

	// Open the specified file path
	file, err := os.Open(file_path); fatal_error(err, 20)

	// Close the file when the current functions returns
	defer file.Close()

	// Open hash interface
	hash := md5.New()

	// Copy the file in the hash interface and check for any error
	if _, err := io.Copy(hash, file); err != nil {
		fatal_error(err, 31)
	}

	// Get the 16 bytes hash
	hash_in_bytes := hash.Sum(nil)[:16]

	// Converts the bytes to string and return it
	file_md5 = hex.EncodeToString(hash_in_bytes)
	return file_md5
}

func file_monitoring(sleep_time float32) {
	for true {
		name, size, mode, modtime, isdir, ino, hlinks, slink, resolve_symbolic := get_stats(os.Args[1])

		// Get the file name
		fmt.Println("Name:\t\t", name)

		// Get file size
		fmt.Println("Size:\t\t", size)

		// Get the file mode
		fmt.Println("Mode:\t\t", mode)

		// Get the file modification time
		fmt.Println("ModTime:\t", modtime)

		// Check if is file or a directory
		fmt.Println("IsDir:\t\t", isdir)

		// Get the file inode
		fmt.Println("Inode:\t\t", ino)

		// Get the file hardlinks
		fmt.Println("Hlinks:\t\t", hlinks)

		// Check if it is a symbolic link
		fmt.Println("IsSymlink:\t", slink)

		// Follow symbolic link
		fmt.Println("RSymlink:\t", resolve_symbolic)

		if resolve_symbolic != deleted_symlink {
			// Get the file md5 hash
			fmt.Println("MD5:\t\t", get_checksum(os.Args[1]))
		} else {
			fmt.Println("MD5:\t\t N/A")
		}
		time.Sleep(time.Duration(sleep_time)*time.Second)
	}
}

func main() {
	// Create the log file
	logger, err := os.OpenFile("fim.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("[FATAL] Error opening file: %v", err)
	}
	defer logger.Close()
	log.SetOutput(logger)

	// Check the users args
	if len(os.Args) != 2 {
		fmt.Println("[ERROR] You must specify one file path")
		os.Exit(10)
	}

	go file_monitoring(5.0)

	for true {
		time.Sleep(time.Second*60)
	}
}
