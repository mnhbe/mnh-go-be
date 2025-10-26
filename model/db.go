package model

import (
  "fmt"
  "net/http"
  "os"
  "strconv"
  "time"

  "github.com/dgrijalva/jwt-go"
  "github.com/jinzhu/gorm"
  "github.com/qor/admin"
  "github.com/qor/qor"
  "github.com/qor/roles"

  "mnhbe/global"
  "mnhbe/utils"
)

var DBInstance *gorm.DB

func Initialize(fullQORDashboardInitialize bool) {
  dbLogModeString := os.Getenv("DB_LOG_MODE")
  dbLogMode, errParseLogMode := strconv.ParseBool(dbLogModeString)
  if errParseLogMode != nil {
    dbLogMode = true
  }
  databaseName := "Postgres (Production)"
  var errConnectMainDB error

  DBInstance, errConnectMainDB = gorm.Open("postgres", global.Config.PostgresConnectionString)
  if errConnectMainDB != nil {
    fmt.Println("Connect Main DB error:", errConnectMainDB)
    panic(errConnectMainDB)
  }
  DBInstance.DB().SetMaxIdleConns(utils.GetDatabaseMaxIdleConnections())
  DBInstance.DB().SetConnMaxIdleTime(time.Minute * 30)
  //DBInstance.DB().SetMaxOpenConns(0) // Not limit
  if DBInstance != nil && DBInstance.DB().Ping() == nil {
    DBInstance.LogMode(dbLogMode)
    fmt.Println("Yay! " + databaseName + " Database [Main] Connected!")
  } else {
    fmt.Println(databaseName + " Database [Main] Connection ERROR!")
  }
}

func CloseDBInstance() {
  if DBInstance != nil {
    _ = DBInstance.Close()
  }
}

func QorUserGetAuth(userID uint) string {
  var user User
  DBInstance.First(&user, "id = ?", userID)
  if user.ID == 0 {
    return ""
  }
  // Create the token
  claims := jwt.MapClaims{}
  if claims["id"] == nil {
    claims["id"] = user.ID
  }
  claims["type"] = "QOR"
  claims["exp"] = time.Now().Add(time.Minute * 3 * 60).Unix()
  at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
  token, err := at.SignedString([]byte(global.Config.JwtSecret))
  if err != nil {
    return ""
  }

  return token
}

func GenerateTokenMiddleware(Admin *admin.Admin) {
  Admin.GetRouter().Use(&admin.Middleware{
    Name: "generate-token",
    Handler: func(context *admin.Context, middleware *admin.Middleware) {
      writer := context.Writer
      // generate token
      currentUser := context.CurrentUser.(*User)
      tokenString := QorUserGetAuth(currentUser.ID)

      expire := time.Now().Add(24 * time.Hour) // Expires in 24 hours
      cookieToken := http.Cookie{Name: "_token", Path: "/admin", Value: tokenString, Expires: expire, MaxAge: 86400}
      http.SetCookie(writer, &cookieToken)
      middleware.Next(context)
    },
  })
}

func DBAutoMigration() {
  //DBInstance.AutoMigrate(&User{})
}

func SetupRoles() {
  // Register roles
  roles.Register(AdminSystem.ToString(), func(req *http.Request, currentUser interface{}) bool {
    return currentUser != nil && currentUser.(*User).UserRoleID == uint(AdminSystem)
  })
  roles.Register(AdminRegion.ToString(), func(req *http.Request, currentUser interface{}) bool {
    return currentUser != nil && currentUser.(*User).UserRoleID == uint(AdminRegion)
  })
  roles.Register(AdminCompany.ToString(), func(req *http.Request, currentUser interface{}) bool {
    return currentUser != nil && currentUser.(*User).UserRoleID == uint(AdminCompany)
  })
}

type DatabaseManager struct {
  MainDB *gorm.DB
}

func NewDatabase() DatabaseManager {
  return DatabaseManager{
    MainDB: DBInstance,
  }
}

func CloseDatabase() {
  if DBInstance != nil {
    _ = DBInstance.Close()
  }
}

type CheckerFunc func(*gorm.DB, *qor.Context) *gorm.DB

type CheckerSaveFunc func(interface{}, *qor.Context) bool
