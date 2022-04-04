package module

import (
	"github.com/techlabs/swabhav/tsam"
	"github.com/techlabs/swabhav/tsam/models/admin"
	"github.com/techlabs/swabhav/tsam/models/batch"
	"github.com/techlabs/swabhav/tsam/models/blog"
	"github.com/techlabs/swabhav/tsam/models/college"
	"github.com/techlabs/swabhav/tsam/models/community"
	"github.com/techlabs/swabhav/tsam/models/company"
	"github.com/techlabs/swabhav/tsam/models/course"
	"github.com/techlabs/swabhav/tsam/models/faculty"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/programming"
	"github.com/techlabs/swabhav/tsam/models/resource"
	"github.com/techlabs/swabhav/tsam/models/talent"
	"github.com/techlabs/swabhav/tsam/models/talentenquiry"
	"github.com/techlabs/swabhav/tsam/models/test"
)

// Configure will configure all modules
func Configure(app *tsam.App) {
	admin := admin.NewAdministrationModuleConfig(app.DB)
	generalModule := general.NewGeneralModuleConfig(app.DB)
	talentModule := talent.NewTalentModuleConfig(app.DB)
	talentEnquiryModule := talentenquiry.NewTalentEnquiryModuleConfig(app.DB)
	courseModule := course.NewCourseModuleConfig(app.DB)
	collegeModule := college.NewCollegeModuleConfig(app.DB)
	facultyModule := faculty.NewFacultyModuleConfig(app.DB)
	batchModule := batch.NewBatchModuleConfig(app.DB)
	companyModule := company.NewCompanyModuleConfig(app.DB)
	communityModule := community.NewCommunityModuleConfig(app.DB)
	testModule := test.NewTestModuleConfig(app.DB)
	resourceModule := resource.NewResourceModuleConfig(app.DB)
	programmingModule := programming.NewProgrammingModuleConfig(app.DB)
	blogModule := blog.NewBlogModuleConfig(app.DB)

	// Need to create/migrate the general tables first!
	generalModule.TableMigration()

	app.MigrateTables([]tsam.ModuleConfig{admin, talentModule, programmingModule, courseModule, collegeModule,
		facultyModule, batchModule, companyModule, communityModule, testModule, talentEnquiryModule,
		resourceModule, blogModule})
}
