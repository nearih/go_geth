package main

import (
	"fmt"
	"net/http"

	"github.com/globalsign/mgo"

	"github.com/globalsign/mgo/bson"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/viper"
)

func main() {

	//yml
	//mongo:
	// host: 13.250.119.252

	//env
	//MONGO_HOST

	mongoHost := "13.250.119.252"
	mongoUser := "root"
	mongoPass := "example"
	port := ":1323"
	viper.SetConfigName("env.json")
	viper.AddConfigPath(".")
	// err := viper.ReadInConfig()
	// if err != nil {
	// 	panic(fmt.Errorf("Fatal error config file: %s \n", err))
	// }
	// // viper.AutomaticEnv()
	// viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	// mongoHost := viper.GetString("mongo.host")
	// mongoUser := viper.GetString("mongo.user")
	// mongoPass := viper.GetString("mongo.pass")
	// port := ":" + viper.GetString("port")
	// fmt.Println("ddd	", mongoPass)
	connstring := fmt.Sprintf("%v:%v@%v", mongoUser, mongoPass, mongoHost)

	session, err := mgo.Dial(connstring)
	e := echo.New()
	if err != nil {
		e.Logger.Fatal(err)
		return
	}

	h := &handler{
		m: session,
	}

	e.Use(middleware.Logger())
	e.POST("/todos", h.create)
	e.GET("/todos", h.list)
	e.GET("/todos/:id", h.view)
	e.POST("/todos/:id", h.done)
	e.DELETE("/todos/:id", h.delete)
	e.Logger.Fatal(e.Start(port))
}

type handler struct {
	m *mgo.Session
}

type todo struct {
	ID    bson.ObjectId `json:"id" bson:"_id"`
	Topic string        `json:"topic" bson:"topic"`
	Done  bool          `json:"done" bson:"done"`
}

func (h *handler) create(c echo.Context) error {

	session := h.m.Copy()
	defer session.Close()
	var t todo
	if err := c.Bind(&t); err != nil {
		return err
	}

	t.ID = bson.NewObjectId()
	col := session.DB("workshop").C("todos")
	if err := col.Insert(t); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, t)
}

func (h *handler) list(c echo.Context) error {

	session := h.m.Copy()
	defer session.Close()

	var ts []todo
	col := session.DB("workshop").C("todos")
	if err := col.Find(nil).All(&ts); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ts)
}

func (h *handler) view(c echo.Context) error {

	session := h.m.Copy()
	defer session.Close()
	id := bson.ObjectIdHex(c.Param("id"))

	var t todo
	col := session.DB("workshop").C("todos")
	if err := col.FindId(id).One(&t); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, t)
}

func (h *handler) done(c echo.Context) error {

	session := h.m.Copy()
	defer session.Close()
	id := bson.ObjectIdHex(c.Param("id"))

	var t todo
	col := session.DB("workshop").C("todos")
	if err := col.FindId(id).One(&t); err != nil {
		return err
	}

	t.Done = false
	if err := col.UpdateId(id, t); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, t)
}

func (h *handler) delete(c echo.Context) error {

	session := h.m.Copy()
	defer session.Close()
	id := bson.ObjectIdHex(c.Param("id"))

	col := session.DB("workshop").C("todos")

	if err := col.RemoveId(id); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{"result": "success"})
}
