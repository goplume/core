//  Типы были выделенны для того чтобы
//  в списке аргументов избежать несовпадний из за однотипности
package types

type (
    OrderId         string
    ExternalOrderId string // номер заказа, может быть пустым;
    RequestId       string
    PaymentsDayId   string
    Channel         string
    Action          string
)

type Entity interface {
    Id() interface{}
}
