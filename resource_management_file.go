package main

import (
	"io/ioutil"
	"os"
)


const NewFileFlag = os.O_CREATE | os.O_WRONLY
const OwnerRWOnly os.FileMode = 0600


func writeFile1(path, content string) error {
	file, err := os.OpenFile(path, NewFileFlag, OwnerRWOnly)   //  1
	if err != nil { // boilerplate                             //
		return err                                             //
	}                                                          //
	// we should know that we should call file.Close()         //
	defer file.Close() // error is not handled                 //  4 (!)
	                                                           //
	_, err = file.Write([]byte(content))                       //  2
	if err != nil {                                            //
		return err                                             //
	}                                                          //
	// write more content                                      //
	_, err = file.Write([]byte(content))                       //  3
	                                                           //
	return err                                                 //  5
}


func writeFile1_WithErrorHandling(path, content string) error {
	file, err := os.OpenFile(path, NewFileFlag, OwnerRWOnly)   //  1
	if err != nil {                                            //
		return err                                             //
	}                                                          //
	defer func() {                                             //
		closeErr := file.Close()                               //  4
		if err == nil {                                        //
			err = closeErr                                     //  5
		}                                                      //
	}()                                                        //
	                                                           //
	_, err = file.Write([]byte(content))                       //  2
	if err != nil {                                            //
		return err                                             //
	}                                                          //
	// write more content                                      //
	_, err = file.Write([]byte(content))                       //  3
	                                                           //
	return err                                                 //  6   returns err from (3) or closeErr from (4~5)
}

// What is bad about resource management here?
/*
	* You have to remember to close resource     ->   You can forget
    * You have to know how to do it              ->   You can make a mistake
 */

// The solution
/*
	Extract resource management to a separate function:
		* inverse control
		* make intuitive nesting
		* handle error by design
 */


func writeFile2(path, content string) error {
	return NewFileResource(path, NewFileFlag, OwnerRWOnly)(
		func(file *os.File) error {
			_, err := file.Write([]byte(content))
			if err != nil {
				return err
			}
			// more content to the God of content
			_, err = file.Write([]byte(content))
			return err
		},
	)
}


func writeFile3(fr FileResource, content string) error {
	return fr(func(file *os.File) error {
		_, err := file.Write([]byte(content))
		if err != nil {
			return err
		}
		// more content to the God of content
		_, err = file.Write([]byte(content))
		return err
	})
}



////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////





type FileResourceCallback = func (fd *os.File) error
type FileResource = func(callback FileResourceCallback) error

// func NewFileResource(path string, flags int, perm os.FileMode, callback FileResourceCallback) error {

func NewFileResource(path string, flags int, perm os.FileMode) FileResource {

	return func(callback FileResourceCallback) error {

		file, err := os.OpenFile(path, flags, perm)
		if err != nil {
			return err
		}

		err = callback(file)

		if err != nil {
			// try to close, but return user's error anyway
			// or maybe combine in one error
			_ = file.Close()
			return err
		} else {
			return file.Close()
		}
	}
}

var TempFileResource FileResource =
	func(callback FileResourceCallback) error {

		file, err := ioutil.TempFile("", "")
		if err != nil {
			return err
		}
		defer file.Close()
		defer os.Remove(file.Name())

		return callback(file)
	}
