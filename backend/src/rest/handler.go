package rest

import (
	"../dblayer"
	"../models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/charge"
	"github.com/stripe/stripe-go/customer"
	"log"
	"net/http"
	"strconv"

	"github.com/stripe/stripe-go"
)

type HandlerInterface interface {
	GetMainPage(c *gin.Context)
	GetProducts(c *gin.Context)
	GetPromos(c *gin.Context)
	AddUser(c *gin.Context)
	SignIn(c *gin.Context)
	SignOut(c *gin.Context)
	GetOrders(c *gin.Context)
	Charge(c *gin.Context)
}

type Handler struct {
	db dblayer.DBLayer
}

func NewHandler() (HandlerInterface, error) {
	return NewHandlerWithParams("mysql", "root:1234@/gomusic?parseTime=true")
}

func NewHandlerWithParams(dbtype, conn string) (HandlerInterface, error) {
	db, err := dblayer.NewORM(dbtype, conn)
	if err != nil {
		return nil, err
	}
	return &Handler{
		db: db,
	}, nil
}

func (h *Handler) GetMainPage(c *gin.Context) {
	log.Println("Main page....")
	c.String(http.StatusOK, "Main page for secure API!!")
}

func (h *Handler) GetProducts(c *gin.Context) {
	if h.db == nil {
		return
	}

	products, err := h.db.GetAllProducts()

	if err != nil {
		//첫 번째 인자는 HTTP 상태 코드, 두 번째는 응답의 바디
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, products)
}

func (h *Handler) GetPromos(c *gin.Context) {
	if h.db == nil {
		return
	}

	promos, err := h.db.GetPromos()

	if err != nil {
		//첫 번째 인자는 HTTP 상태 코드, 두 번째는 응답의 바디
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, promos)
}

func (h *Handler) SignIn(c *gin.Context) {
	if h.db == nil {
		return
	}

	var customer models.Customer
	err := c.ShouldBindJSON(&customer)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customer, err = h.db.SignInUser(customer.Email, customer.Pass)
	if err != nil {
		//잘못된 패스워드인 경우 forbidden http 에러 반환
		if err == dblayer.ErrINVALIDPASSWORD {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, customer)
}

func (h *Handler) AddUser(c *gin.Context) {
	if h.db == nil {
		return
	}

	var customer models.Customer
	err := c.ShouldBindJSON(&customer)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Println("customer")
		return
	}

	customer, err = h.db.AddUser(customer)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println("db adduser")
		return
	}

	c.JSON(http.StatusOK, customer)
}

func (h *Handler) SignOut(c *gin.Context) {
	if h.db == nil {
		return
	}

	p := c.Param("id")
	id, err := strconv.Atoi(p)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.db.SignOutUserById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

func (h *Handler) GetOrders(c *gin.Context) {
	if h.db == nil {
		return
	}

	p := c.Param("id")

	id, err := strconv.Atoi(p)

	if err != nil {
		fmt.Println("1")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orders, err := h.db.GetCustomerOrderByID(id)

	if err != nil {
		fmt.Println("2")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

func (h *Handler) Charge(c *gin.Context) {
	if h.db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server database error"})
		return
	}

	request := struct {
		models.Order
		Remember    bool   `json:"rememberCard"`
		UseExisting bool   `json:"useExisting"`
		Token       string `json:"token"`
	}{}

	err := c.ShouldBindJSON(&request)
	log.Printf("request: %+v  \n",request)
	//파싱 중 에러 발생 시 보고 후 반환
	if err != nil {
		fmt.Println("json parsing error")
		c.JSON(http.StatusBadRequest, request)
		return
	}

	//input stripe secret key
	stripe.Key = "secret key"

	chargeP := &stripe.ChargeParams{
		Amount:      stripe.Int64(int64(request.Price)),
		Currency:    stripe.String("usd"),
		Description: stripe.String("GoMusic Charge..."),
	}

	//params:=&stripe.PaymentIntentParams{
	//	Amount: stripe.Int64(1000),
	//	Currency: stripe.String("usd"),
	//	PaymentMethodTypes: stripe.StringSlice([]string{
	//		"card",
	//	}),
	//	ReceiptEmail: stripe.String("zlarbals@example.com"),
	//}
	//
	//pi,_:=paymentintent.New(params)
	//
	//fmt.Println(pi)

	stripeCustomerID := ""

	fmt.Println(request.UseExisting)

	if request.UseExisting{
		log.Println("Getting credit card id...")
		stripeCustomerID, err = h.db.GetCreditCardID(request.CustomerID)
		if err != nil {
			log.Println(err)
			fmt.Println("1")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		cp := &stripe.CustomerParams{}
		cp.SetSource(request.Token)
		customer, err := customer.New(cp)
		if err != nil {
			fmt.Println("2")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		stripeCustomerID = customer.ID
		if request.Remember {
			err = h.db.SaveCreditCardForCustomer(request.CustomerID, stripeCustomerID)
			if err != nil {
				fmt.Println("3")
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
	}
	chargeP.Customer = stripe.String(stripeCustomerID)
	_, err = charge.New(chargeP)
	if err != nil {
		fmt.Println("4")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = h.db.AddOrder(request.Order)
	if err != nil {
		fmt.Println("5")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
