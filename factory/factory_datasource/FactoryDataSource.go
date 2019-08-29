package factory_datasource

import (
    "database/sql"
    "github.com/goplume/core/configuration"
    "github.com/goplume/core/utils/logger"
    "time"
)

//Factory for DataSource from database
// Configuration
//    root path = "services." + serviceName + ".data-source."
//
// attrs:
//    connection-url:  -db_user/db_password@(description = (address=(protocol=tcp)(host=db-host)(port
//    driver: goracle, postgress, etc
//    enable: true | false - enabled or disable datasource
//    max-idle-connections: max idle connection
//    max-open-connections: max open connection
//    connection-max-lifetime-sec: connection max lifetime in seconds
type FactoryDataSource struct {
    Log *logger.Logger
}

func (this FactoryDataSource) InitFactory() {
}

func (this FactoryDataSource) CreateDataSource(
    serviceName string,
) *sql.DB {
    config := configuration.NewServiceConfiguration(serviceName, "data-source.%s", this.Log)
    this.Log.RLog.Info("Read configuration from context " + config.ConfigurationContext)

    CONNECTION_URL := config.GetString("connection-url")
    DRIVER := config.GetString("driver")
    DATAOURCE_ENABLE := config.GetBool("enable")
    MaxIdleConnections := config.GetInt("max-idle-connections")               // 10
    MaxOpenConnections := config.GetInt("max-open-connections")               // 100
    ConnectionsMaxLifetimeSec := config.GetInt("connection-max-lifetime-sec") // 100

    if DATAOURCE_ENABLE == false {
        return nil
    }

    db, err := sql.Open(DRIVER, CONNECTION_URL)
    //this.DataSource = db
    if err != nil {
        this.Log.RLog.Error(err)
        return nil
    }

    // Test connection to database
    if err = db.Ping(); err != nil {
        this.Log.RLog.Error(err)
    } else {
        this.Log.RLog.Info("Database connected successfully")
    }

    // SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
    db.SetMaxIdleConns(MaxIdleConnections)

    // SetMaxOpenConns sets the maximum number of open connections to the database.
    db.SetMaxOpenConns(MaxOpenConnections)

    // SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
    db.SetConnMaxLifetime(time.Duration(ConnectionsMaxLifetimeSec) * time.Second)

    return db
}
