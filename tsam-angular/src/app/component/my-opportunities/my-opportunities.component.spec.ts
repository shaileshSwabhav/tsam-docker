import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { MyOpportunitiesComponent } from './my-opportunities.component';

describe('MyOpportunitiesComponent', () => {
  let component: MyOpportunitiesComponent;
  let fixture: ComponentFixture<MyOpportunitiesComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ MyOpportunitiesComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(MyOpportunitiesComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
