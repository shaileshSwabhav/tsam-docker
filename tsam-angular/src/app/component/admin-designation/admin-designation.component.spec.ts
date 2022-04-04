import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { AdminDesignationComponent } from './admin-designation.component';

describe('AdminDesignationComponent', () => {
  let component: AdminDesignationComponent;
  let fixture: ComponentFixture<AdminDesignationComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ AdminDesignationComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(AdminDesignationComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
