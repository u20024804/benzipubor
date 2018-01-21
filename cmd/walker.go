/*
 *     benzipubor
 *     Copyright (C) 2018 bobo liu
 *
 *     This program is free software: you can redistribute it and/or modify
 *     it under the terms of the GNU General Public License as published by
 *     the Free Software Foundation, either version 3 of the License, or
 *     (at your option) any later version.
 *
 *     This program is distributed in the hope that it will be useful,
 *     but WITHOUT ANY WARRANTY; without even the implied warranty of
 *     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *     GNU General Public License for more details.
 *
 *     You should have received a copy of the GNU General Public License
 *     along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"os"
)

func walkDir(path string) ([]string, error) {
	lo := *l
	lo.SetPrefix("[WALK]")

	f, err := os.Open(path)
	if err != nil {
		lo.Println(err)
		return nil, err
	}
	defer f.Close()

	fds, err := f.Readdir(0)
	if err != nil {
		lo.Println(err)
		return nil, err
	}
	fileList := make([]string, 0)
	for _, fd := range fds {
		if fd.IsDir() {
			lo.Fatalln("仅支持一层子目录")
		} else {
			if fileFilter(path) {
				lo.Println("New image file:", path)
				fileList = append(fileList, path)
			} else {
				lo.Println("Unknown file:", path)
				continue
			}
		}
	}

	return fileList, nil
}

func walkRoot(path string) ([]unit, error) {
	lo := *l
	lo.SetPrefix("[WALK]")

	f, err := os.Open(path)
	if err != nil {
		lo.Println(err)
		return nil, err
	}
	defer f.Close()
	os.Chdir(path)

	fds, err := f.Readdir(0)
	if err != nil {
		lo.Println(err)
		return nil, err
	}

	units := make([]unit, 0)
	haveDir := false
	for _, fd := range fds {
		if fd.IsDir() {
			haveDir = true
			if autoMode == MODE_AIO {
				list, err := walkDir(fd.Name())
				if err != nil {
					lo.Println("Cannot access:", fd.Name())
					lo.Fatalln(err)
				}
				units = append(units, unit{Name: fd.Name(), ImageList: list})
			}
			// todo: MODE_SINGLE
		} else {
			if haveDir {
				lo.Fatalln("存在子目录，根目录下的文件将被忽略:", fd.Name())
				continue
			}
		}
	}
	if !haveDir {
		fd, _ := f.Stat()
		list, err := walkDir(path)
		if err != nil {
			lo.Println("Cannot access:", fd.Name())
			lo.Fatalln(err)
		}
		units = append(units, unit{Name: fd.Name(), ImageList: list})
	}
	return units, nil
}

type unit struct {
	Name      string
	ImageList []string
}

func fileFilter(fn string) bool {
	mime := getMime(fn)
	switch mime {
	case "image/jpeg", "image/png", "image/gif":
		return true
	default:
		return false
	}
}
