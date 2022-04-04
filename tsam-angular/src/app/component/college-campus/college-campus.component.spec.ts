import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { CollegeCampusComponent } from './college-campus.component';

describe('CollegeCampusComponent', () => {
  let component: CollegeCampusComponent;
  let fixture: ComponentFixture<CollegeCampusComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ CollegeCampusComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(CollegeCampusComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
