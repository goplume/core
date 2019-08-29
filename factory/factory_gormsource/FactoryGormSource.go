package factory_gormsource

import (
    "fmt"
    "github.com/goplume/core/configuration"
    "github.com/goplume/core/utils/logger"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/postgres"
    "time"
)

//Factory for GORM DataSource
// Configuration
//    root path = "services." + serviceName + ".gorm.datasource."
//
// attrs:
//    host: host
//    username: user name
//    password: password
//    dbname: database name
//    dialect: database dialect (postgress)
//    enable: true | false - enabled or disable datasource
//    max-idle-connections: max idle connection
//    max-open-connections: max open connection
//    connection-max-lifetime-sec: connection max lifetime in seconds
//    log-mode: true | false

type FactoryGormSource struct {
    Log *logger.Logger
}

func (this *FactoryGormSource) InitFactory() {
}

func (this *FactoryGormSource) CreateGormSource(
    serviceName string,
) (dataSource *gorm.DB) {
    config := configuration.NewServiceConfiguration(serviceName, "gorm.datasource.%s", this.Log)
    this.Log.RLog.Info("Read configuration from context " + config.ConfigurationContext)

    DB_HOST := config.GetString("host")
    DB_USERNAME := config.GetString("username")
    DB_DBNAME := config.GetString("dbname")
    DB_DIALECT := config.GetString("dialect")
    DB_DBPASSWD := config.GetString("password")
    //CONNECTION_URL := config.GetString("connectionUrl")
    MaxIdleConnections := config.GetInt("max-idle-connections")               // 10
    MaxOpenConnections := config.GetInt("max-open-connections")               // 100
    ConnectionsMaxLifetimeSec := config.GetInt("connection-max-lifetime-sec") // 100
    LogMode := config.GetBool("log-mode")                                     // 100

    connectionString := fmt.Sprintf(
        "host=%s user=%s dbname=%s sslmode=disable password=%s",
        DB_HOST, DB_USERNAME, DB_DBNAME, DB_DBPASSWD,
    )

    connectionStringForLogging := fmt.Sprintf(
        "host=%s user=%s dbname=%s sslmode=disable password=%s",
        DB_HOST, DB_USERNAME, DB_DBNAME, config.MaskSensitiveValue(DB_DBPASSWD),
    )

    this.Log.RLog.Info(connectionStringForLogging)
    dataSource, err := gorm.Open(DB_DIALECT, connectionString)

    // Enable Logger, show detailed log
    dataSource.LogMode(LogMode)

    dataSource.SetLogger(logger.NewWithLogger(this.Log.RLog))

    if err != nil {
        this.Log.RLog.Error("Failed to connect database")
        this.Log.RLog.Error(err)
        return
    }

    // Get generic database object sql.DB to use its functions
    db := dataSource.DB()

    // SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
    db.SetMaxIdleConns(MaxIdleConnections)

    // SetMaxOpenConns sets the maximum number of open connections to the database.
    db.SetMaxOpenConns(MaxOpenConnections)

    // SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
    db.SetConnMaxLifetime(time.Duration(ConnectionsMaxLifetimeSec) * time.Second)

    // Disable Logger, don't show any log even errors
    //db.LogMode(false)

    // Debug a single operation, show detailed log for this operation
    //db.Debug().Where("name = ?", "jinzhu").First(&User{})

    //defer vars.dataSource.Close()

    return dataSource
}
