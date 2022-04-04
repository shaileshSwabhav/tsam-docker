import { Component, OnInit } from '@angular/core';
import { FormGroup, FormBuilder, Validators, FormControl } from '@angular/forms';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { AdminService } from 'src/app/service/admin/admin.service';
import { DummyDataService } from 'src/app/service/dummydata/dummy-data.service';
import { NgbModal } from '@ng-bootstrap/ng-bootstrap';

@Component({
      selector: 'app-adminService-role',
      templateUrl: './admin-role.component.html',
      styleUrls: ['./admin-role.component.css']
})
export class AdminRoleComponent implements OnInit {

      //role
      roleForm: FormGroup;
      roles: any[] = [];

      //search
      searchForm: FormGroup;
      searched = false;

      //modal
      modalButton: string;
      modalHeader: string;
      modalAction: () => void;
      formAction: (param?) => void;
      isRoleUpadteClick: boolean;
      selectedRole: any;
      selectedRoleId: string;
      modalRef: any;

      //pagination
      totalRoles: number;
      limit: number;
      offset: number;
      currentPage: number;

      constructor(
            private formBuilder: FormBuilder,
            private dummyDataService: DummyDataService,
            private adminService: AdminService,
            private spinnerService: SpinnerService,
            private modalService: NgbModal

      ) {
            this.limit = 5;
            this.offset = 0;
      }


      get ongoingOperations() {
            return this.spinnerService.ongoingOperations
      }

      ngOnInit() {

            this.onAddRoleClick();
            this.formAction();
            this.loadData();
            this.createSearchForm();
      }

      createSearchForm(): void {
            this.searchForm =
                  this.formBuilder.group({
                        roleName: new FormControl(null),
                        level: new FormControl(null)
                  })
      }
      // On Add New Role Button Click.
      onAddRoleClick(): void {
            this.updateCommonParam('Add Role', 'Add New Role', this.addRole, this.createRoleForm, false);
      }

      // On Add New Role Button Click.
      onUpdateRoleClick(param: string): void {
            this.spinnerService.loadingMessage = "Getting role details";

            this.updateCommonParam('Update Role', 'Update Role',
                  this.updateRole, this.updateRoleForm, true);
            let index = param;
            this.adminService.getRolebyId(index).subscribe(
                  data => {
                        this.updateRoleForm(data)
                  },
                  (error) => {
                        console.log(error);

                        if (error.error) {
                              alert(error.error);
                              return;
                        }
                        alert(error.statusText);
                  }
            )
            // console.log(selectedData)
            // this.updateRoleForm(selectedData);
      }

      // Create Role Form.
      createRoleForm(): void {
            this.roleForm = this.formBuilder.group({
                  roleName: ['', Validators.required],
                  level: ['', Validators.required]
            });
      }

      // Update Role Form.
      updateRoleForm(data: any): void {
            this.roleForm.patchValue(
                  {
                        roleName: data.roleName,
                        level: data.level
                  }
            );

            this.selectedRoleId = data.id;
      }

      // Add New Role
      addRole(): void {
            this.spinnerService.loadingMessage = "Adding role";

            let data = this.roleForm.value;
            this.adminService.postRole(data).subscribe(data => {
                  this.modalRef.close();
                  this.getAllRoles();
                  alert('Data added successfully');
            }, (error) => {
                  console.log(error);

                  if (error.error) {
                        alert(error.error);
                        return;
                  }
                  alert(error.statusText);
            })
      }

      // Update Role.
      updateRole(): void {
            this.spinnerService.loadingMessage = "Updating role";

            let data = this.roleForm.value;
            this.adminService.updateRolebyId(data, this.selectedRoleId).
                  subscribe(data => {
                        this.modalRef.close();
                        this.roleForm.reset();
                        this.getAllRoles();
                        alert("Data Updated Successfully");
                  }, (error) => {
                        console.log(error);

                        if (error.error) {
                              alert(error.error);
                              return;
                        }
                        alert(error.statusText);
                  })
      }

      // Get Role By ID;
      getRoles(param: string): any {
            for (let index = 0; index < this.roles.length; index++) {
                  if (this.roles[index].rolecode == param) {
                        return index;
                  }
            }
            return undefined;
      }


      assignSelectedRoleId(id: string): void {
            this.selectedRoleId = id;
      }

      // Delete Role
      deleteRole(): void {
            this.spinnerService.loadingMessage = "Deleting role";
            this.modalRef.close();

            this.adminService.deleteRolebyId(this.selectedRoleId).
                  subscribe(
                        data => {
                              this.getAllRoles();
                              alert("Data deleted successfully");
                        },
                        error => {
                              console.log(error)

                              if (error.error) {
                                    alert(error.error);
                                    return;
                              }
                              alert(error.statusText);
                        }
                  )
      }


      // Update Model Param.
      updateCommonParam(modelbutton: string, modelheader: string, modalAction: any, formAction: any, isRoleUpadteClick: boolean): void {
            this.modalAction = modalAction;
            this.modalButton = modelbutton;
            this.modalHeader = modelheader;
            this.isRoleUpadteClick = isRoleUpadteClick;
            this.formAction = formAction;
      }

      // Load data.
      loadData(): void {
            this.getAllRoles();
      }

      // Pagination Controls
      changePage($event: any): void {

            // $event will be the page number & offset will be 1 less than it.
            this.offset = $event - 1
            this.currentPage = $event;
            this.getAllRoles();
      }

      // On Limit Change
      limitChange(): void {
            this.offset = 1;
            this.getAllRoles();
      }

      // On Page Number Change.
      offsetChange(param: number): void {
            this.offset = param;
            this.getAllRoles();
      }

      ///get All Roles
      getAllRoles(): void {
            this.spinnerService.loadingMessage = "Getting roles";
            this.adminService.getAllRole().subscribe(res => {
                  this.roles = [];
                  this.roles = this.roles.concat(res.body);
                  console.log(this.roles)
                  this.totalRoles = parseInt(res.headers.get("X-Total-Count"))

            }, err => {
                  console.log(err.error.error)

                  if (err.error) {
                        alert(err.error);
                        return;
                  }
                  alert(err.statusText);
            })
      }

      //search roles
      search(): void {
            this.spinnerService.loadingMessage = "Seraching roles";

            this.searched = true;
            let data = this.searchForm.value;
            console.log(data)
            this.adminService.getSearchedRole(data, this.limit, this.offset).subscribe(res => {
                  this.roles = [];
                  this.roles = this.roles.concat(res.body);
                  this.totalRoles = parseInt(res.headers.get('X-Total-Count'));

                  console.log(this.roles)
            },
                  err => {
                        console.log(err.error.error)

                        if (err.error) {
                              alert(err.error);
                              return;
                        }
                        alert(err.statusText);
                  }
            )
      }

      //reest the search and get all roles
      reset(): void {
            this.spinnerService.loadingMessage = "Getting roles";

            this.searchForm.reset();
            this.getAllRoles();
            this.searched = false;
      }


      //reset add/update form
      close(): void {
            this.roleForm.reset();
      }

      //open modal for add/update form
      openRoleFormModal(roleFormModal: any): void {
            this.modalRef = this.modalService.open(roleFormModal, { ariaLabelledBy: 'modal-basic-title', keyboard: false, backdrop: 'static', size: 'xl' });
            /*this.modalRef.result.then((result) => {
            }, (reason) => {

            });*/
      }

      //open modal for delete confirmation
      openDeleteConfirmationModal(deleteConfirmationModal: any): void {
            this.modalRef = this.modalService.open(deleteConfirmationModal, { ariaLabelledBy: 'modal-basic-title', keyboard: false, backdrop: 'static' });
            /*this.modalRef.result.then((result) => {
            }, (reason) => {

            });*/
      }

      //check form validation
      validate(): void {
            if (this.roleForm.invalid) {
                  this.roleForm.markAllAsTouched();
            }
            else {
                  this.modalAction();
            }
      }

}
