package iwouldlove

import (
  "fmt"
  "runtime"
  "reflect"
  "testing"
)

var errorLocations = make(map[int]struct{})
var currentBefore func()
const colorRed = "\x1b[31;1m"
const colorGreen = "\x1b[32;1m"
const colorNormal = "\x1b[0m"

func before(beforeFunction func()) {
  currentBefore = beforeFunction
}

func it(description string, itFunction func()) {
  if currentBefore != nil {
    currentBefore()
  }

  itFunction()
  _, _, line, _ := runtime.Caller(1)

  if _, errors := errorLocations[line]; errors {
    fmt.Printf("\t%s%s%s\n", colorRed, description, colorNormal)
  } else {
    fmt.Printf("\t%s%s%s\n", colorGreen, description, colorNormal)
  }
}

func idLove (t *testing.T) func(interface{}, string, interface{}) {
  return func(actual interface{}, matcher string, expected interface{}) {
    _, file, line, _ := runtime.Caller(1)
    var errorMessage string

    switch matcher {
    case "to equal":
      if !reflect.DeepEqual(actual, expected) {
        errorMessage = fmt.Sprintf("Expected %v to equal %v", actual, expected)
      }
    case "to not equal":
      if reflect.DeepEqual(actual, expected) {
        errorMessage = fmt.Sprintf("Expected %v to not equal %v", actual, expected)
      }
    default:
      errorMessage = "Given matcher \"" + matcher + "\" not known"
    }

    if errorMessage != "" {
      _, _, itLocation, _ := runtime.Caller(3)
      errorLocations[itLocation] = struct{}{}
      t.Errorf("\n%sError%s in %s:%d\n%s", colorRed, colorNormal, file, line, errorMessage)
    }
  }
}

func Init(t *testing.T) (func(interface {}, string, interface {}), func(string, func()), func(func())) {
  return idLove(t), it, before
}
