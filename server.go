package main

import (
	"fmt"
	"net/http"
	"strings"

	// "strings"
	// "os"

	"github.com/spf13/viper"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	// Echo instance
	e := echo.New()

	//
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	mongoHost := viper.GetString("mongo.host")
	mongoUser := viper.GetString("mongo.user")
	mongoPass := viper.GetString("mongo.password")
	port := viper.GetString("port")

	connString := fmt.Sprintf("%v:%v@%v", mongoUser, mongoPass, mongoHost)
	session, err := mgo.Dial(connString)
	if err != nil {
		e.Logger.Fatal(err)
		return
	}

	h := &handler{
		m: session,
	}

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Route => handler
	e.GET("/todos", h.list)
	e.GET("/todos/:id", h.view)
	e.PUT("/todos/:id", h.done)
	e.DELETE("/todos/:id", h.delete)
	e.POST("/todos", h.create)

	// Start server
	e.Logger.Fatal(e.Start(port))
}

type todo struct {
	ID    bson.ObjectId `json: "id", bson:"_id"`
	Topic string        `json: "topic", bson:"topic"`
	Done  bool          `json: "done", bson:"done"`
}

type handler struct {
	m *mgo.Session
}

func (h *handler) delete(c echo.Context) error {

	session := h.m.Copy()
	defer session.Close()

	id := bson.ObjectIdHex(c.Param("id"))

	ss := session.DB("workshop").C("tumit-todos")

	if err := ss.RemoveId(id); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, nil)
}

func (h *handler) done(c echo.Context) error {

	session := h.m.Copy()
	defer session.Close()

	id := bson.ObjectIdHex(c.Param("id"))

	ss := session.DB("workshop").C("tumit-todos")
	var t todo
	if err := ss.FindId(id).One(&t); err != nil {
		return err
	}

	t.Done = true
	if err := ss.UpdateId(id, t); err != nil {
		return err
	}

	c.JSON(http.StatusOK, t)

	return nil
}

func (h *handler) view(c echo.Context) error {

	session := h.m.Copy()
	defer session.Close()

	id := bson.ObjectIdHex(c.Param("id"))

	ss := session.DB("workshop").C("tumit-todos")
	var t todo
	if err := ss.FindId(id).One(&t); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, t)
}

func (h *handler) list(c echo.Context) error {

	session := h.m.Copy()
	defer session.Close()

	ss := session.DB("workshop").C("tumit-todos")
	var ts []todo
	if err := ss.Find(nil).All(&ts); err != nil {
		return err
	}

	c.JSON(http.StatusOK, ts)

	return nil
}

func (h *handler) create(c echo.Context) error {

	session := h.m.Copy()
	defer session.Close()

	var t todo
	if err := c.Bind(&t); err != nil {
		return err
	}

	t.ID = bson.NewObjectId()

	ss := session.DB("workshop").C("tumit-todos")
	if err := ss.Insert(t); err != nil {
		return err
	}

	c.JSON(http.StatusOK, t)

	return nil
}

// func create(c echo.Context) error {
// 	var t todo
// 	if err := c.Bind(&t); err != nil {
// 		return err
// 	}

// 	session, err := mgo.Dial("root:example@13.250.119.252")

// 	if err != nil {
// 		return err
// 	}

// 	t.ID = bson.NewObjectId()

// 	ss := session.DB("workshop").C("tumit-todos")
// 	if err := ss.Insert(t); err != nil {
// 		return err
// 	}

// 	c.JSON(http.StatusOK, t)

// 	return nil
// }
