MAIN APP

1. run the app from main.go
2. use port 8080 for requests
3. use GET /seed endpoint to generate db and seed
4. use POST /login to login (a jwt token cookie will be generated) (email and password as form params)
5. use POST /add-to-cart (add jwt token cookie) to add products to cart (product id and quantity as form params) 
6. use POST /checkout (add jwt token cookie) to checkout
7. use POST /logout (add jwt token cookie) to logout

TEST APP

run the app from tester.go