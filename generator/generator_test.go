package gen

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"testing"

	"github.com/pmezard/go-difflib/difflib"
)

func TestTemplateModel_Generate(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	t.Logf("Current test filename: %s", filename)

	testData, _ := filepath.Abs(filepath.Join(filepath.Dir(filename), "testdata"))
	absSource := filepath.Join(testData, "source")
	//absTarget := filepath.Join(testData, "target")
	absTarget, _ := ioutil.TempDir("", "tymlate-test")

	absConfig := filepath.Join(testData, "conf.yml")
	tm, err := NewTemplateModel(absSource, absTarget, absConfig, false)
	if err != nil {
		panic(err)
	}

	err = os.Setenv("TEST_TYMLATE_USER", "Tymlate tester")
	if err != nil {
		panic(err)
	}

	if err := tm.Generate(); err != nil {
		panic(err)
	}

	//diff folder with target
	absExpected := filepath.Join(testData, "target")
	diff, err := dirDiff(absTarget, absExpected)

	if diff {
		t.Errorf("The target output does not match the expected one ")
	}

	err = os.Setenv("TEST_TYMLATE_USER", "")

}

//dirDiff compares two folders (left and right)
//and it returns a bool true if the filenames and structure matches
//and an int that denotes how manny files are different by content.
//As a side effect it prints out the diffs.
//Note: The sole purpose of this function is to test the generated output.
func dirDiff(left, right string) (bool, error) {

	leftFiles, _ := files(left)
	rightFiles, _ := files(right)

	numberOfDiffs := 0

	if len(leftFiles) != len(rightFiles) {
		return true, nil
	}

	for i, leftFile := range leftFiles {
		if filepath.Base(rightFiles[i]) != filepath.Base(leftFile) { //not found on right side
			return true, nil
		}

		leftStr, _ := ioutil.ReadFile(leftFile)
		rightStr, _ := ioutil.ReadFile(rightFiles[i])

		if !bytes.Equal(leftStr, rightStr) {
			fmt.Printf("!! Differs... %s  <---> %s\n", leftFile, rightFiles[i])
			//fmt.Println(leftStr, rightStr))
			ddiff := difflib.ContextDiff{
				A:        difflib.SplitLines(string(leftStr)),
				B:        difflib.SplitLines(string(rightStr)),
				FromFile: "Left",
				ToFile:   "Right",
				Context:  3,
				Eol:      "\n",
			}
			res, _ := difflib.GetContextDiffString(ddiff)
			fmt.Print(strings.Replace(res, "\t", " ", -1))
			numberOfDiffs++
		}
	}

	if numberOfDiffs > 0 {
		return true, nil
	}

	return false, nil
}

func files(src string) ([]string, error) {
	var files []string

	err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})

	if err != nil {
		return nil, err
	}

	sort.Strings(files)
	return files[1:], nil
}
