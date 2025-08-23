package webServer

import (
	"client-runaway-zenoti/internal/config"
	"client-runaway-zenoti/internal/db/models"
	"client-runaway-zenoti/internal/runway"
	"client-runaway-zenoti/internal/tgbot"
	"client-runaway-zenoti/internal/zenoti"
	runwayv2 "client-runaway-zenoti/packages/runwayV2"
	zenotiv1 "client-runaway-zenoti/packages/zenotiV1"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
	"github.com/Lavina-Tech-LLC/lavinagopackage/v2/conf"
	"github.com/gin-gonic/gin"
)

func Listen() {

	router := gin.Default()
	router.Use(ErrorHandler)
	router.Use(corsMiddleware)
	router.Use(recovery)

	router.GET("/", Default)
	router.GET("/health", ok)
	router.GET("/contact/:id/:location", goToContact)
	router.StaticFS("/createNew", gin.Dir(conf.GetPath()+"internal/webserver/webform", false))
	router.StaticFS("/chat", gin.Dir(conf.GetPath()+"internal/webserver/ringcentral/build/", false))
	router.StaticFS("/report", gin.Dir(conf.GetPath()+"internal/webserver/report-app/build/", false))
	router.StaticFS("/app", gin.Dir(conf.GetPath()+"internal/webserver/report-app/build/", false))
	router.POST("/create/:id/:location", create)
	router.POST("/webhook/newOpportunity", newOpportunity)
	router.GET("/locations", locations)

	fmt.Printf("Path for ws is %s\n", conf.GetPath()+"internal/webserver/webform")

	// chat stuff
	setChatRoutes(router)
	setReportRoutes(router)
	setAuthRoutes(router)
	setSettingsRoutes(router)
	setSurveyRoutes(router)
	setRunwayRoutes(router)
	setChatlyRoutes(router)
	setZenotiWebhooksRoutes(router)
	setRoutes(router)

	srv := &http.Server{
		Addr:    config.Confs.Settings.SrvAddress,
		Handler: router,
	}

	fmt.Printf("Server is listening to %s", srv.Addr)

	cert := config.Confs.Settings.Cert
	key := config.Confs.Settings.Key

	var err error
	if cert == "" && key == "" {
		err = srv.ListenAndServe()
	} else {
		err = srv.ListenAndServeTLS(cert, key)
	}
	// service connections
	if err != nil && err != http.ErrServerClosed {
		fmt.Printf("listen: %s\n", err)
	} else {
		fmt.Printf("Server is listening to %s", srv.Addr)
	}

}

func Default(c *gin.Context) {
	c.Writer.WriteHeader(204)
}
func ok(c *gin.Context) {
	c.Writer.Write([]byte("OK"))
	c.Writer.WriteHeader(200)
}

func create(c *gin.Context) {

	// AUTH
	body, _ := io.ReadAll(c.Request.Body)
	sign := c.GetHeader("s")

	stringToSign := string(body) + ";lkjasdfOoiuwer92384554fsldkjf0)Ojklsdf"

	hash := sha256.Sum256([]byte(stringToSign))
	if sign != fmt.Sprintf("%x", hash[:]) {
		c.JSON(400, "")
	}

	// Get location
	l := models.Location{}
	l.Get(c.Param("location"))

	// Create guest object

	info := zenotiv1.Personal_info{}
	json.Unmarshal(body, &info)

	guest := zenotiv1.Guest{
		Personal_info: info,
	}

	client := zenoti.MustGetClientFromLocation(l)

	guest, err := client.GuestsCreate(guest)
	lvn.GinErr(c, 500, err, "Error creating guest")

	c.Writer.WriteHeader(200)
	c.Writer.Write([]byte(zenoti.GuestsGetLinkByIdLocationId(guest.Id, l.Id)))
}

func newOpportunity(c *gin.Context) {

	info := struct {
		Contact_Id  string `json:"contact_id"`
		Email       string `json:"email"`
		Phone       string `json:"phone"`
		LocationStr string `json:"Location"`
		Location    struct {
			Id string `json:"id"`
		} `json:"location"`
	}{}

	// // Read the raw body first
	// body, err := io.ReadAll(c.Request.Body)
	// if err != nil {
	// 	fmt.Println("Error reading body:", err)
	// 	return
	// }
	// fmt.Println("Raw JSON body:", string(body))

	// // Now bind again (but need to reset the body since it was already read)
	// c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	err := c.BindJSON(&info)
	if err != nil {
		fmt.Println("Error binding json:", err)
	}

	locId := info.Location.Id
	l := models.Location{}
	l.Get(locId)

	svc := runway.GetSvc()

	cli, err := svc.NewClientFromId(locId)
	lvn.GinErr(c, 500, err, "Error creating client")

	op, err := cli.OpportunitiesGetAll(runwayv2.OpportunitiesFilter{
		ContactId:  info.Contact_Id,
		PipelineId: l.PipelineId,
	})

	lvn.GinErr(c, 500, err, "Error finding opportunity")

	for _, o := range op {
		tgbot.Notify("New Opportunity", fmt.Sprintf("New opportunity for %s %s\n%s\n%s\n%s\n%s",
			info.Contact_Id,
			info.Email,
			info.Phone,
			l.Name,
			l.Id,
			o.Id), false)
		runway.UpdateNote(o, l, false)
	}
	c.Writer.WriteHeader(200)
}

func goToContact(c *gin.Context) {

	// finding the runway contact
	locId := c.Param("location")
	contactId := c.Param("id")
	contact, err := runway.GetContactByIdLocationId(locId, contactId)

	lvn.GinErr(c, 500, err, "Error finding location")

	// finding the zenoti guest
	guest, err := zenoti.GuestsGetByPhoneNumberLocationId(
		contact.Phone,
		contact.Email,
		locId,
	)

	if err != nil && err.Error() == "guest not found" {
		c.Redirect(307, fmt.Sprintf("%s/createNew?name=%s&lastName=%s&email=%s&phone=%s&l=%s&id=%s",
			config.Confs.Settings.SrvDomain,
			contact.FirstName,
			contact.LastName,
			contact.Email,
			contact.Phone,
			locId,
			contact.Id))
		c.Writer.WriteHeader(200)
		return
	}

	if err != nil {
		c.Writer.WriteHeader(400)
		c.Writer.Write([]byte(fmt.Sprintf("%s", err)))
		return
	}
	link := zenoti.GuestsGetLinkByIdLocationId(guest.Id, locId)
	fmt.Println("Redirecting to: ", link)
	c.Redirect(302, link)
}

func ErrorHandler(c *gin.Context) {
	//c.Next()

	if len(c.Errors) == 0 {
		return
	}
	msg, _ := json.Marshal(c.Errors)
	c.JSON(http.StatusInternalServerError, msg)
}

func locations(c *gin.Context) {
	type Res struct {
		Name string `json:"name"`
		Id   string `json:"id"`
	}
	res := []Res{}
	for _, l := range config.GetLocations() {
		res = append(res, Res{Name: l.Name, Id: l.Id})
	}
	c.JSON(200, res)
}
