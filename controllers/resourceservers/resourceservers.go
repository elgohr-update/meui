package resourceservers

import (
  "net/url"
  "net/http"
  "sort"
  "github.com/sirupsen/logrus"
  "github.com/gin-gonic/gin"
  //"github.com/gin-contrib/sessions"
  idp "github.com/charmixer/idp/client"

  "github.com/charmixer/meui/app"
  "github.com/charmixer/meui/config"
  "github.com/charmixer/meui/environment"

  bulky "github.com/charmixer/bulky/client"
)

type ResourceServerTemplate struct {
  Id string
  Name string
  Description string
  Audience string
  GrantsUrl string
}

func ShowResourceServers(env *environment.State) gin.HandlerFunc {
  fn := func(c *gin.Context) {

    log := c.MustGet(environment.LogKey).(*logrus.Entry)
    log = log.WithFields(logrus.Fields{
      "func": "ShowResourceServers",
    })

    identity := app.RequireIdentity(c)
    if identity == nil {
      log.Debug("Missing Identity")
      c.AbortWithStatus(http.StatusForbidden)
      return
    }

    idpClient := app.IdpClientUsingAuthorizationCode(env, c)

    status, responses, err := idp.ReadResourceServers(idpClient, config.GetString("idp.public.url") + config.GetString("idp.public.endpoints.resourceservers.collection"), nil)
    if err != nil {
      log.Debug(err.Error())
      c.AbortWithStatus(http.StatusInternalServerError)
      return
    }

    if status != 200 {
      log.Debug(err.Error())
      c.AbortWithStatus(http.StatusInternalServerError)
      return
    }

    var uiCreatedRs []ResourceServerTemplate

    var resourceservers idp.ReadResourceServersResponse
    status, _ = bulky.Unmarshal(0, responses, &resourceservers)
    if status == 200 {

      for _, rs := range resourceservers {

        uiClient := ResourceServerTemplate{
          Id:        rs.Id,
          Name:      rs.Name,
          Description: rs.Description,
          Audience: rs.Audience,
        }
        uiCreatedRs = append(uiCreatedRs, uiClient)

      }

    }

    sort.Slice(uiCreatedRs, func(i, j int) bool {
		  return uiCreatedRs[i].Name > uiCreatedRs[j].Name
	  })

    c.HTML(http.StatusOK, "resourceservers.html", gin.H{
      "title": "Resource Servers",
      "links": []map[string]string{
        {"href": "/public/css/dashboard.css"},
      },
      "id": identity.Id,
      "user": identity.Username,
      "name": identity.Name,
      "created": uiCreatedRs,
    })
  }
  return gin.HandlerFunc(fn)
}