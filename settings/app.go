package settings

import (
    "os"
    "fmt"
    "regexp"
    "stt_back/flags"
)

const RabbitServerPassword = "fXotHxuY"
const RabbitServerLogin = "stt_back"
const RabbitServer = "core.140140.ru"
const RabbitServerPort = "5672"
const rabbitServerVirtualhost = "/microservices"
const rabbitServerVirtualhostTest = "/microservices-test"

const MicroserviceAuthKey = "e452ff4c-4b05-4835-bd1f-91b9c66f2ae1"

func IsDev() bool {
    var matchDev = regexp.MustCompile("^/tmp/go-build")
    return matchDev.Match([]byte(os.Args[0])) || *flags.IsDev
}

func GetVirtualHost() string {

    for _, param := range os.Args {

        if param == "dev" {
            fmt.Println("RPC use DEV environment")
            return rabbitServerVirtualhostTest
        }
    }

    if IsDev() {
        fmt.Println("RPC use DEV environment")
        return rabbitServerVirtualhostTest
    }

    fmt.Println("RPC use PRODUCTION environment")
    return rabbitServerVirtualhost
}