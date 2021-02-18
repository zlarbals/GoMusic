package rest

import(
	"github.com/gin-gonic/gin"
)

func RunAPI(address string) error{
	//Gin 엔진
	r:=gin.Default()

	//핸들러
	h,_:=NewHandler()

	//메인 페이지
	r.GET("/",h.GetMainPage)

	//상품 목록
	r.GET("/products",h.GetProducts)

	//프로모션 목록
	r.GET("/promos", h.GetPromos)

	userGroup := r.Group("/user")
	{
		//사용자 로그아웃
		userGroup.POST("/:id/signout",h.SignOut)
		//주문 내역
		userGroup.GET("/:id/orders",h.GetOrders)
	}

	usersGroup := r.Group("/users")
	{
		//사용자 로그인
		usersGroup.POST("/signin",h.SignIn)
		//사용자 추가
		usersGroup.POST("",h.AddUser)
		//신용카드 결제
		usersGroup.POST("/charge",h.Charge)
	}

	//서버 시작
	return r.Run(address)
}
