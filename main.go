package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	if err := mainErr(); err != nil {
		fmt.Printf("%v", err)
		os.Exit(-1)
	}
}

func mainErr() error {
	var err error

	result := make([]string, 10, 10)
	safeWaitGroup := NewSafeWaitGroup()
	for i := 0; i < len(result); i += 1 {
		i := i // capture i
		safeWaitGroup.Run(func() {
			result[i] = strings.Repeat("*", i)
		})
	}
	safeWaitGroup.Wait()
	fmt.Printf("%v", result)

	err = writeFile1("./test1.txt", "test1")
	if err != nil {
		return err
	}

	err = writeFile2("./test2.txt", "test2")
	if err != nil {
		return err
	}

	{
		file := NewFileResource("./test3.txt", os.O_CREATE|os.O_WRONLY, 0600)

		err = writeFile3(file, "test3")
		if err != nil {
			return err
		}

		err = writeFile3(TempFileResource, "whatever")
		if err != nil {
			return err
		}

		err = TempFileResource(func(fd1 *os.File) error {
			return TempFileResource(func(fd2 *os.File) error {
				_, err := fd1.Write([]byte("hi!"))
				if err != nil {
					return err
				}
				_, err = fd2.Write([]byte("hi!"))
				return err
			})
		})
		if err != nil {
			return err
		}
	}


	{
		db := NewDBResource("sqlite3", "./demo.sqlite")

		err = db.Use(func(db *sql.DB) error {
			err := initDB(db)
			if err != nil {
				return err
			}

			result1, err := helloSql_Cool(db, "MessageBird")
			if err != nil {
				return err
			}
			println(result1)

			result2, err := helloSql_NotCoolAtAll(db, "Twilio")
			if err != nil {
				return err
			}
			println(result2)

			return nil
		})
		if err != nil {
			return err
		}
	}

	return nil
}

