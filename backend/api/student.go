package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo"

	jsonpatch "github.com/evanphx/json-patch/v5"

	authn "gitlab.com/timeterm/timeterm/backend/auhtn"
	"gitlab.com/timeterm/timeterm/backend/database"
)

func (s *Server) getStudent(c echo.Context) error {
	id := c.Param("id")

	uid, err := uuid.Parse(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}

	user, ok := authn.UserFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Not authenticated")
	}

	dbStudent, err := s.db.GetStudent(c.Request().Context(), uid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not read student from database")
	}

	if user.OrganizationID != dbStudent.OrganizationID {
		return echo.NewHTTPError(http.StatusUnauthorized, "Student does not belong to user's organization")
	}

	apiStudent := StudentFrom(dbStudent)
	return c.JSON(http.StatusOK, apiStudent)
}

func (s *Server) getStudents(c echo.Context) error {
	var params paginationParams
	err := c.Bind(&params)
	if err != nil {
		return err
	}

	user, ok := authn.UserFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Not authenticated")
	}

	dbStudents, err := s.db.GetStudents(c.Request().Context(), database.GetStudentsOpts{
		OrganizationID: user.OrganizationID,
		Limit:          params.MaxAmount,
		Offset:         params.Offset,
	})
	if err != nil {
		s.log.Error(err, "could not read students from database")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not read students from database")
	}

	apiStudents := PaginatedStudentsFrom(dbStudents)
	return c.JSON(http.StatusOK, apiStudents)
}

func (s *Server) createStudent(c echo.Context) error {
	user, ok := authn.UserFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Not authenticated")
	}

	var student Student
	err := c.Bind(&student)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not bind data")
	}
	student.OrganizationID = user.OrganizationID

	dbStudent, err := s.db.CreateStudent(c.Request().Context(), StudentToDB(student))
	if err != nil {
		s.log.Error(err, "could not create student")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not create student")
	}

	if errors.Is(err, database.ErrConflict) {
		return echo.NewHTTPError(http.StatusConflict, "User with same name already exists")
	}

	apiStudent := StudentFrom(dbStudent)
	return c.JSON(http.StatusOK, apiStudent)
}

func (s *Server) patchStudent(c echo.Context) error {
	ctx := c.Request().Context()
	studentID := c.Param("id")

	user, ok := authn.UserFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Not logged in")
	}

	uid, err := uuid.Parse(studentID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}

	patchData, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		s.log.Error(err, "could not read request body")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not read request body")
	}

	oldDBStudent, err := s.db.GetStudent(ctx, uid)
	if err != nil {
		s.log.Error(err, "could not read student from database")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not read student from database")
	}

	if oldDBStudent.OrganizationID != user.OrganizationID {
		return echo.NewHTTPError(http.StatusUnauthorized, "User not in organization")
	}

	oldAPIStudent := StudentFrom(oldDBStudent)

	jsonStudent, err := json.Marshal(oldAPIStudent)
	if err != nil {
		s.log.Error(err, "could not marshal the old student")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not marshal the old student")
	}

	newJSONStudent, err := jsonpatch.MergePatch(jsonStudent, patchData)
	if err != nil {
		s.log.Error(err, "could not patch the student")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not patch the student")
	}

	var newAPIStudent PatchedStudent
	err = json.Unmarshal(newJSONStudent, &newAPIStudent)
	if err != nil {
		s.log.Error(err, "could not unmarshal patched student")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not unmarshal patched student")
	}

	newAPIStudent.ID = oldDBStudent.ID
	newAPIStudent.OrganizationID = oldDBStudent.OrganizationID

	newDBStudent := StudentToDB(newAPIStudent.Student)

	err = s.db.ReplaceStudent(ctx, newDBStudent)
	if err != nil {
		s.log.Error(err, "could not update the student in the database")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not update the student in the database")
	}

	if newAPIStudent.CardID.Value != nil {
		cardID := []byte(*newAPIStudent.CardID.Value)
		if err = s.db.ReplaceStudentCard(ctx, user.OrganizationID, newDBStudent.ID, cardID); err != nil {
			s.log.Error(err, "could not replace student card")
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not update student card in the database")
		}
	} else if newAPIStudent.CardID.ExplicitlyNull {
		if err = s.db.DeleteStudentCards(ctx, user.ID); err != nil {
			s.log.Error(err, "could not delete student cards")
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not delete student cards from the database")
		}
	}

	return c.JSON(http.StatusOK, newAPIStudent)
}

type deleteStudentsParams struct {
	StudentDs []uuid.UUID `json:"studentIds"`
}

func (s *Server) deleteStudents(c echo.Context) error {
	var params deleteStudentsParams
	err := c.Bind(&params)
	if err != nil {
		return err
	}

	user, ok := authn.UserFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Not authenticated")
	}

	allInOrg, err := s.db.AreStudentsInOrganization(c.Request().Context(), user.OrganizationID, params.StudentDs...)
	if err != nil {
		s.log.Error(err, "could not get students in organization")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve student information")
	}
	if !allInOrg {
		return echo.NewHTTPError(http.StatusUnauthorized, "Not all students are in user's organization")
	}

	err = s.db.DeleteStudents(c.Request().Context(), params.StudentDs)
	if err != nil {
		s.log.Error(err, "could not delete students")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not delete students")
	}

	return c.NoContent(http.StatusNoContent)
}
