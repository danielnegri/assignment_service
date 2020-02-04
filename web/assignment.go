package web

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/surajjain36/assignment_service/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//createAssignment inserts the assignment data to the DB
//Body type: raw
//Input params:
//{
// "title" : "test 19",
// "name" :"test desc dd",
// "description": "1",
// "type" : "type",
// "duration" : 20
// }
// @Summary Insert an Assignments to the DB
// @tags v1
// @Accept  json
// @Produce  json
// @Success 200 {array} ReadResponse
// @Router /assignment [post]
func (s *Service) createAssignment(c *gin.Context) {
	var assignment models.Assignment
	msg := "Something went wrong"
	statusCode := http.StatusBadRequest
	var currentTime = time.Now()

	if err := c.ShouldBindJSON(&assignment); err != nil {
		log.Error("Error with i/p params: ", err.Error())
		s.responseWriter(c, gin.H{"message": msg}, statusCode)
		return
	}

	assignment.ID = primitive.NewObjectID()
	assignment.CreatedAt = currentTime.UTC()
	assignment.UpdatedAt = currentTime.UTC()
	assignment.Status = "started"
	if _, err := s.mdb.Insert("assignment", &assignment); err == nil {
		s.responseWriter(c, gin.H{"message": "Assignment created successfully"}, http.StatusOK)
		return
	}

	s.responseWriter(c, gin.H{"message": msg}, statusCode)
	return
}

//getAssignment retrieves the assignment data from DB
//Input params: id = <MongoDB ID>
// @Summary Get the Asiignment
// @tags v1
// @Accept  json
// @Produce  json
// @Success 200 {array} ReadResponse
// @Router /assignment/{id} [get]
func (s *Service) getAssignment(c *gin.Context) {
	var assignment models.Assignment
	statusCode := http.StatusBadRequest
	res := gin.H{"message": "Something went wrong", "data": nil}

	if aID := c.Param("id"); aID != "" {
		aIDHex, err := primitive.ObjectIDFromHex(aID)
		if err != nil {
			log.Println("Error parsing assignment id.", err)
			s.responseWriter(c, res, statusCode)
			return
		}
		if err := s.mdb.FindOne("assignment", bson.M{"_id": aIDHex}, &assignment); err == nil {
			res["message"] = "Assignment retrieved successfully"
			res["data"] = &assignment
			s.responseWriter(c, res, http.StatusOK)
			return
		}
	}

	s.responseWriter(c, res, statusCode)
	return
}

//searchAssignmentByTags retrieves the assignment which are matched by Tags
//Input params:
//		tags = <comma separated strings> (ex, Hello,World)
//		pn = <Page Number>
//		pp = <Per Page>
// @Summary Get the list of Assignments for the matching tags
// @tags v1
// @Accept  json
// @Produce  json
// @Success 200 {array} ReadResponse
// @Router /search/assignment [get]
func (s *Service) searchAssignmentByTags(c *gin.Context) {
	var assignments []models.Assignment
	statusCode := http.StatusBadRequest
	res := gin.H{"message": "Something went wrong", "data": nil}

	pn, err := strconv.Atoi(c.DefaultQuery("pn", "1"))
	if err != nil {
		log.WithError(err).Warn("Invalid page number")
		pn = 1
	}

	pp, err := strconv.Atoi(c.DefaultQuery("pp", "20"))
	if err != nil {
		log.WithError(err).Warn("Invalid per page value")
		pp = 20
	}

	if tags := c.DefaultQuery("tags", ""); tags != "" {
		tagsArr := bson.A{}
		for _, t := range strings.Split(tags, ",") {
			tagsArr = append(tagsArr, t)
		}
		pipeline := bson.A{
			bson.M{"$match": bson.M{"tags": bson.M{"$all": tagsArr}}},
			bson.M{"$skip": pp * (pn - 1)},
			bson.M{"$limit": pp},
		}
		if err := s.mdb.Aggregate("assignment", pipeline, &assignments); err == nil {
			res["message"] = "Assignment retrieved successfully"
			res["data"] = &assignments
			s.responseWriter(c, res, http.StatusOK)
			return
		}
	}

	s.responseWriter(c, res, statusCode)
	return

}
