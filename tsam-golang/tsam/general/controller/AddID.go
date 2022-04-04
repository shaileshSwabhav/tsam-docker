package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/batch"
	cmp "github.com/techlabs/swabhav/tsam/models/company"
	"github.com/techlabs/swabhav/tsam/models/faculty"
	tal "github.com/techlabs/swabhav/tsam/models/talent"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// AddIDToEntity adds Id value to the respective fields with created_at
func (controller *CountryController) AddIDToEntity(w http.ResponseWriter, r *http.Request) {
	fmt.Println("==========================================================")
	fmt.Println("===============ADD ID TO ENTITY CALLED=======================")
	fmt.Println("==========================================================")
	tableName := mux.Vars(r)["tableName"]
	stringToStructMap := map[string]interface{}{
		"talent_academics":    &tal.Academic{},
		"talent_call_records": &tal.CallRecord{},
		"talent_experiences":  &tal.Experience{},
		// "faculties_technologies":       &faculty.Technology{},
		"faculty_experiences":          &faculty.Experience{},
		"faculty_academics":            &faculty.Academic{},
		"batch_talents":                &batch.MappedTalent{},
		"company_enquiry_call_records": &cmp.CallRecord{},
	}
	db := controller.CountryService.DB
	// 	SELECT @i:= 0;
	// UPDATE tsmdb.batch_talents SET `id`=@i:=@i+1;
	var id uuid.UUID
	var total int
	err := db.Model(stringToStructMap[tableName]).Count(&total).Error
	fmt.Println("===============Total=======================", total)
	if err != nil {
		fmt.Println("===============ErrorxxxCount=======================", err)
		return
	}
	for i := 0; i < total; i++ {
		id = util.GenerateUUID()
		err := db.Debug().Exec("UPDATE "+tableName+" SET `id`=?, `created_at`=?,`updated_at`=?,`tenant_id`=? WHERE `id`=?",
			id, time.Now(), time.Now(), "7ca2664b-f379-43db-bdf9-7fdd40707219", fmt.Sprintf("%d", i+1)).
			Error
		if err != nil {
			fmt.Println("===============Error=======================", err)
			return
		}
		fmt.Println("===============YAY=======================")
	}

}

// ChangeTalentCode changes the code of talent to new format
func (controller *CountryController) ChangeTalentCode(w http.ResponseWriter, r *http.Request) {

	talents := &[]tal.Talent{}

	uow := repository.NewUnitOfWork(controller.CountryService.DB, false)
	err := controller.CountryService.Repository.GetAll(uow, talents, repository.Select([]string{"`first_name`", "`code`", "`id`"}))
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	fmt.Println("================================================================================")
	for _, talent := range *talents {
		fmt.Print("Talents ->", talent.Code, " ")
		if _, err := strconv.ParseInt(talent.Code, 10, 64); err == nil {
			fmt.Printf("%q looks like a number.", talent.Code)
			// Assign Talent Code
			var codeError error
			fmt.Println("===================================CODE=============================================")
			talent.Code, codeError = util.GenerateUniqueCode(uow.DB, talent.FirstName, "`code` = ?", &tal.Talent{})
			if codeError != nil {
				log.NewLogger().Error(codeError.Error())
				uow.RollBack()
				web.RespondError(w, errors.NewHTTPError("Internal server error", http.StatusInternalServerError))
				return
			}
			fmt.Println(" code ->", talent.Code)
			fmt.Println("================================================================================")
		}
		fmt.Println()

		err = controller.CountryService.Repository.UpdateWithMap(uow, &tal.Talent{}, map[interface{}]interface{}{
			"code": talent.Code,
		}, repository.Filter("`id` = ?", talent.ID))
		if err != nil {
			log.NewLogger().Error(err.Error())
			web.RespondError(w, err)
			return
		}
	}
	fmt.Println("================================================================================")
	uow.Commit()
}
