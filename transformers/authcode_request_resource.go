package transformers

import (
	"fmt"
	"strings"
	"time"

	"github.com/companieshouse/chs.go/log"
	"github.com/companieshouse/emergency-auth-code-api/models"
	"github.com/companieshouse/emergency-auth-code-api/utils"
)

// AuthCodeResourceRequestToDB will take the input request from the REST call
// and transform it to a dao ready for insertion into the database
func AuthCodeResourceRequestToDB(req *models.AuthCodeRequest) *models.AuthCodeRequestResourceDao {

	createdAt := time.Now().Truncate(time.Millisecond)

	id := utils.GenerateID()

	etag, err := utils.GenerateEtag()
	if err != nil {
		log.Error(fmt.Errorf("error generating etag: [%s]", err))
	}

	format := "/emergency-auth-code-service/auth-code-requests/%s"
	self := fmt.Sprintf(format, id)

	dao := &models.AuthCodeRequestResourceDao{
		ID: id,
		Data: models.AuthCodeRequestDataDao{
			CompanyNumber:   req.CompanyNumber,
			OfficerID:       req.OfficerID,
			OfficerUraID:    req.OfficerUraID,
			OfficerForename: req.OfficerForename,
			OfficerSurname:  req.OfficerSurname,
			Status:          "pending",
			CreatedAt:       &createdAt,
			SubmittedAt:     nil,
			Kind:            "emergency-auth-code-request",
			Etag:            etag,
			CreatedBy: models.CreatedByDao{
				Email:    req.CreatedBy.Email,
				ID:       req.CreatedBy.ID,
				Forename: req.CreatedBy.Forename,
				Surname:  req.CreatedBy.Surname,
			},
			Links: models.AuthCodeResourceLinksDao{
				Self: self,
			},
		},
	}

	return dao
}

// AuthCodeRequestResourceDaoToResponse will transform an auth code resource dao
// into an http response entity
func AuthCodeRequestResourceDaoToResponse(model *models.AuthCodeRequestResourceDao) *models.AuthCodeRequestResourceResponse {
	return &models.AuthCodeRequestResourceResponse{
		CompanyNumber: model.Data.CompanyNumber,
		CompanyName:   model.Data.CompanyName,
		UserID:        model.Data.CreatedBy.ID,
		UserEmail:     model.Data.CreatedBy.Email,
		OfficerID:     model.Data.OfficerID,
		OfficerUraID:  model.Data.OfficerUraID,
		OfficerName:   strings.Join([]string{model.Data.OfficerForename, model.Data.OfficerSurname}, " "),
		Status:        model.Data.Status,
		CreatedAt:     model.Data.CreatedAt,
		SubmittedAt:   model.Data.SubmittedAt,
		Etag:          model.Data.Etag,
		Kind:          model.Data.Kind,
		Links: models.AuthCodeRequestResourceLinks{
			Self: model.Data.Links.Self,
		},
	}
}
