# Frekwent API


## GET Routes
- No auth routes:
  - `http://localhost:8080/products`- gets all products
    - `http://localhost:8080/products?trending=true` - get all trending products
    - `http://localhost:8080/products?new=true` - get all products sorted from newest to oldest
    - `http://localhost:8080/products?page=2` - get all products from page 2 of the pagination
    - `http://localhost:8080/products?page=2&new=true` - get all products from page 2 of the pagination sorted from newest to oldest
  - `http://localhost:8080/products/<id>` - get product by ID
  - `http://localhost:8080/list/optout/<email>` - opt out of emails
  - `http://localhost:8080/products/category/<category>` - get all products by category
- Auth routes:
  - `http://localhost:8080/iriguchi/coupon` - get coupons
  - `http://localhost:8080/iriguchi/coupon/code/<coupon code>` - get coupons by code
  - `http://localhost:8080/iriguchi/ip/get/<ip>` - get ip by ip address
  - `http://localhost:8080/iriguchi/ip` - get all ip addresses
  - `http://localhost:8080/iriguchi/purchases` - get all purchases
  - `http://localhost:8080/iriguchi/purchases/<id>` - get purchases by ID
  - `http://localhost:8080/iriguchi/customer` - get all customers
  - `http://localhost:8080/iriguchi/customer/name?firstName=<firstName>&lastName=<lastName>` - get customers by first name and last name
  - `http://localhost:8080/iriguchi/auth/logout` - logout
  - `http://localhost:8080/iriguchi/email` - get all emails
  - `http://localhost:8080/iriguchi/email/get/<email address>` - get emails by email address
  - `http://localhost:8080/iriguchi/customer/optin` - get all customers that have opted in to receiving emails
  - `http://localhost:8080/iriguchi/email/status/<status>` - get all emails by status
    - `http://localhost:8080/iriguchi/email/status/success`
    - `http://localhost:8080/iriguchi/email/status/pending`
  - All the get all routes above have two optional query params, `page`(pretty much pagination), `new`(sorted from newest to oldest)
    - Ex. `http://localhost:8080/iriguchi/email?page=8&new=false`
  - `http://localhost:8080/iriguchi/purchases/transactions/<state>` - get all transactions
    - `http://localhost:8080/iriguchi/purchases/transactions/ca` - California
    - `http://localhost:8080/iriguchi/purchases/transactions/all` - all states
  - `http://localhost:8080/iriguchi/items/name?name=<name>` - gets a product by name
    - `http://localhost:8080/iriguchi/items/name?name=test product`
---
  
## POST Routes
- No auth route:
  - `http://localhost:8080/iriguchi/auth/login` - login
```
{
    "email": "admin@admin.com",
    "password": "password"
}
```
- Auth routes:
- `http://localhost:8080/iriguchi/items` - create product
``` 
{
    "name": "test product5",
    "images": [],
    "price": "10.01",
    "quantity": 20,
    "description": "desc...",
    "ingredients": [],
    "category": "faceWash"
}
```
- `http://localhost:8080/iriguchi/coupon` - create coupon
``` 
{
    "code": "testcode",
    "percentage": 20
}
```
- `http://localhost:8080/products/buy` - purchase product
``` 
{
    "firstName": "John",
    "lastName": "Doe",
    "email": "jdoe@gmail.com",
    "streetAddress": "1st st.",
    "optionalAddress": "Apt 1",
    "city": "Los Angeles",
    "state": "CA",
    "zipCode": "90043",
    "purchasedItems": [
        {
            "id": "61c32a0a61d12a5f03b73fc7",
            "name": "test product updated",
            "images": [],
            "price": "22",
            "quantity": 25,
            "description": "new description",
            "ingredients": [
                "soap",
                "water"
            ]
        }
    ],
    "finalPrice": 120,
    "tax": 20,
    "infoEmailOptIn": true
}
```
- `http://localhost:8080/iriguchi/coupon/send/couponPromotion` - email a coupon to a customer
``` 
{
    "email": "jdoe@gmail.com",
    "subject": "test coupon",
    "content": "test coupon",
    "couponCode": "testcode"
}
```
- `http://localhost:8080/iriguchi/coupon/send/customerInteraction` - email a customer
``` 
{
    "email": "jdoe@gmail.com",
    "subject": "test customer interaction",
    "content": "test customer interaction"
}
```
- `http://localhost:8080/iriguchi/coupon/mass/send` - mass email coupon codes to all opted in customers
``` 
{
    "couponCode": "testcode"
} 
```
---

## PUT Routes(All Routes Are Authenticated)
- `http://localhost:8080/iriguchi/items/name/<id>` - update product name
``` 
{
    "name": "new product updated"
}
```
- `http://localhost:8080/iriguchi/items/price/<id>` - update product price
``` 
{
    "price": "25.90"
}
```
- `http://localhost:8080/iriguchi/items/quantity/<id>` - update product quantity
``` 
{
    "quantity": 28
}
```
- `http://localhost:8080/iriguchi/items/description/<id>` - update product description
``` 
{
    "description": "new product description"
}
```
- `http://localhost:8080/iriguchi/items/ingredients/<id>` - update product ingredients
``` 
{
    "ingredients": ["soap", "water", "sugar"]
}
```
- `http://localhost:8080/iriguchi/purchases/shipped/<id>` - update purchase's shipped status
``` 
{
    "shipped": true,
    "trackingId": "dkjkdlkldskdkldsjfkldsfkls"
}
```
- `http://localhost:8080/iriguchi/purchases/delivered/<id>` - update purchase's delivered status
``` 
{
    "delivered": true
}
```
- `http://localhost:8080/iriguchi/purchases/tracking/<id>` - update purchase's tracking ID 
``` 
{
    "trackingId": "lssiodjidkdkdd"
}
```
- `http://localhost:8080/iriguchi/purchases/address/<id>` - update purchase's address
``` 
{
    "streetAddress": "2nd st.",
    "optionalAddress": "Apt 2",
    "city": "Los Angeles",
    "state": "CA",
    "zipCode": "90044"
}
```
- `http://localhost:8080/iriguchi/items/category/<id>` - update category
``` 
{
    "category": "faceWash"
} 
```
---

## DELETE Routes(All Routes Are Authenticated)
- `http://localhost:8080/iriguchi/items/delete/<id>` - delete product by ID
- `http://localhost:8080/iriguchi/coupon/code/<code>` - delete coupon by code
---

## Prometheus Command
- `docker run -p 9090:9090 -v $(pwd)/prometheus.yml:/etc/prometheus/prometheus.yml prom/prometheus`
---

## Docker Build
- `docker build -t freq .`
---