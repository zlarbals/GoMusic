package rest

import(
	"github.com/gin-gonic/gin"
)

func RunAPI(address string) error{
	//Gin 엔진
	r:=gin.Default()

	//핸들러
	h,_:=NewHandler()

	//상품 목록
	r.GET("/products",h.GetProducts)

	//프로모션 목록
	r.GET("/promos", h.GetPromos)

	//사용자 로그인
	r.POST("/users/signin", h.SignIn)

	//사용자 추가
	r.POST("/users", h.AddUser)

	//사용자 로그아웃
	r.POST("/user/:id/signout", h.SignOut)

	//주문 내역
	r.GET("/user/:id/orders", h.GetOrders)

	//신용카드 결제
	r.POST("/users/charge", h.Charge)

	return r.Run(address)
}
