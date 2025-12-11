package routes

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	httpSwagger "github.com/swaggo/http-swagger"

	"go-transport-hub/utils"
)

// Prevent abnormal shutdown while panic
func recoverHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %+v", err)
				log.Print(string(debug.Stack()))
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

// Put params in context for sharing them between handlers
func wrapHandler(next http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := context.WithValue(r.Context(), "params", ps)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

type TokenRes struct {
	Exp int `json:"exp"`
}

// RouterConfig function
func RouterConfig() http.Handler {
	router := httprouter.New()
	router.PanicHandler = panicHandler
	//indexHandlers := alice.New(recoverHandler)

	//indexHandlersWithLogger = alice.New(loggingHandler, recoverHandler)
	//recoverHandler := alice.New(loggingHandler, recoverHandler)
	recoverHandler := alice.New(recoverHandler)

	setPingRoutes(router, recoverHandler)
	login(router, recoverHandler)
	employeeRole(router, recoverHandler)
	employeeV1(router, recoverHandler)
	driverHub(router, recoverHandler)
	vehicleHub(router, recoverHandler)
	vendorHub(router, recoverHandler)
	customerHub(router, recoverHandler)
	branchHub(router, recoverHandler)
	loadingUnloadingHub(router, recoverHandler)
	tripManagementHub(router, recoverHandler)
	managePod(router, recoverHandler)
	lrReceipt(router, recoverHandler)
	dashboardAPI(router, recoverHandler)
	tripcompleteXLS(router, recoverHandler)
	attendanceEmp(router, recoverHandler)
	paymentAndInvoice(router, recoverHandler)
	invoiceHub(router, recoverHandler)
	// notificationRoutes(router, recoverHandler) // Temporarily disabled

	//commonAPIs
	router.GET("/view/v1/image", wrapHandler(recoverHandler.ThenFunc(ViweImage)))
	router.GET("/v1/:orgId/prerequisite", wrapHandler(recoverHandler.ThenFunc(createPrerequisite)))

	router.Handler("GET", "/swagger", httpSwagger.WrapHandler)
	router.Handler("GET", "/swagger/:one", httpSwagger.WrapHandler)
	router.Handler("GET", "/swagger/:one/:two", httpSwagger.WrapHandler)
	router.Handler("GET", "/swagger/:one/:two/:three", httpSwagger.WrapHandler)
	router.Handler("GET", "/swagger/:one/:two/:three/:four", httpSwagger.WrapHandler)
	router.Handler("GET", "/swagger/:one/:two/:three/:four/:five", httpSwagger.WrapHandler)
	return router
}

func panicHandler(w http.ResponseWriter, r *http.Request, c interface{}) {
	fmt.Printf("(thub)Recovering from panic-Reason: %+v \n", c.(error))
	debug.PrintStack()
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(c.(error).Error()))
}

type CustomWriter struct {
	http.ResponseWriter
	StatusCode int
}

func loggingHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		t := time.Now().In(utils.TimeLoc())
		cw := &CustomWriter{w, 0}

		next.ServeHTTP(cw, r)

		elapsed := time.Since(t)
		ip := utils.GetClientIP(r)
		log.Printf("(IP: %s) %d [%s] %s %s", ip, cw.StatusCode, r.Method, r.URL.String(), elapsed)
	}

	return http.HandlerFunc(fn)
}
