package utils

import (
  "fmt"
  "os"
  "strconv"
)

func GetDatabaseMaxIdleConnections() (dbMaxIdleConnectionsReturned int) {
  defer func() {
    fmt.Println("Database Max Idle Connection: ", dbMaxIdleConnectionsReturned)
  }()

  defaultDbMaxConnections := 20
  dbMaxIdleConnectionsString := os.Getenv("DB_MAX_IDLE_CONNECTIONS")
  dbMaxIdleConnections, errParseNum := strconv.ParseInt(dbMaxIdleConnectionsString, 10, 64)
  if errParseNum != nil {
    dbMaxIdleConnectionsReturned = defaultDbMaxConnections
    return
  }
  if dbMaxIdleConnections < 10 || dbMaxIdleConnections > 50 {
    dbMaxIdleConnectionsReturned = defaultDbMaxConnections
    return
  }
  dbMaxIdleConnectionsReturned = int(dbMaxIdleConnections)
  return
}
