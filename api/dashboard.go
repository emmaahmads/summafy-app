package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type Link struct {
	Href string
	Rel  string
}

type dashboard_history struct {
	Date     string
	User     string
	Activity string
	Document Link
}

func (server *Server) HandlerLandingPage(c *gin.Context) {
	username_str, _ := sessions.Default(c).Get("username").(string)
	var activity []dashboard_history
	activity_type := map[int]string{
		0: "uploaded",
		1: "generated a summary on",
		2: "deleted",
		3: "changed the summary on",
		4: "downloaded"}
	activities, err := server.store.Queries.GetAllActivities(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	//  TODO - add a new table that consists of existing documents on the S3 bucket
	// this table should not be published
	// activity lists and view lists should only display the documents that are in this table

	// TODO - do not display redundant documents activities, display only the latest
	// example when user keeps uploading the same file

	// TODO - enable session

	for a := range activities {
		user, _ := server.store.Queries.GetUser(context.Background(), activities[a].Username)
		doc, err := server.store.Queries.GetDocument(context.Background(), activities[a].DocumentID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		act := dashboard_history{
			Date:     string(activities[a].CreatedAt.Format("2006-01-02 15:04:05")),
			User:     user.FullName,
			Activity: activity_type[int(activities[a].Activity)],
			Document: Link{
				Href: fmt.Sprintf("/viewdoc/%d", (activities[a].DocumentID)),
				Rel:  doc.FileName,
			},
		}

		activity = append(activity, act)
	}

	c.HTML(200, "dashboard.html", gin.H{"act": activity, "user": username_str})
	c.Header("Content-Type", "text/html")
}
