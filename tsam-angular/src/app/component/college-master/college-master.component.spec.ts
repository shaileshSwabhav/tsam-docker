import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { CollegeMasterComponent } from './college-master.component';

describe('CollegeMasterComponent', () => {
  let component: CollegeMasterComponent;
  let fixture: ComponentFixture<CollegeMasterComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ CollegeMasterComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(CollegeMasterComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
